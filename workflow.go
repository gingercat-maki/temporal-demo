package approval

import (
	"approval-workflow-demo-go/external"
	"approval-workflow-demo-go/model"
	"errors"
	"fmt"
	"time"

	"github.com/luci/go-render/render"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var (
	activities = &ExpenseApprovalActivity{}
)

// This is how expense directly use the temporal with no middle-layer in tech architecture
// temporal is more for SAGA than business workflows
func ExpenseApprovalWorkflow(ctx workflow.Context, request *model.ExpenseApprovalStartRequest) (*model.ExpenseApprovalContent, error) {

	logger := workflow.GetLogger(ctx)
	current_version := 2
	v := workflow.GetVersion(ctx, "ExpenseApprovalWorkflowVersion", workflow.DefaultVersion, workflow.Version(current_version))
	logger.Debug("current version of definition is %d, the run version is %d", current_version, v)

	// this get actually marks the version data, the version can be granular which we don't use right now
	if err := checkValidityOfRequest(request); err != nil {
		return nil, err
	}

	// try some bad thing here
	logger.Debug("before sleep")
	time.Sleep(1 * time.Microsecond)
	// time.Sleep(1 * time.Second)
	now := time.Now()
	fmt.Printf("Now: %v\n", now)
	logger.Debug("after sleep")

	// the domain-content defined in this workflow, used in the lifecycle of execuation
	// we inject it in the registered query, so that it will stored in the temporal's db
	content := model.NewExpenseApprovalContent(request)

	// CONFIG OF THIS WORKFLOW DEFINITION
	// This function is in fact the TEMPLATE OF WORKFLOW DEFINITION
	// with the right config, the workflow can definition can complete
	config := getWorkflowConfig(request.ExpenseTeam)
	if config == nil {
		return nil, fmt.Errorf(
			"No ExpenseApprovalWorkflow is configed for %s", request.ExpenseTeam)
	}
	logger.Debug("Expenseflow config %v", render.Render(config))

	// KEY REACTIVE CAPABLITIES ARE DEFINED HERE:
	// - selectors mean how to interact with this workflow
	// - runtime queries mean what can be obtained from outside
	submitSelector := createSubmitSelector(ctx, content)
	err := registerRuntimeQueries(ctx, content)
	if err != nil {
		logger.Error(err.Error())
		return content, err
	}
	logger.Info("Regiestered queries and selectors on workflowExecuation info: %v", workflow.GetInfo(ctx).WorkflowExecution)

	// RUN STEPS IN THE CONFIG
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    1 * time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Second * 3, // 100 * InitialInterval
		MaximumAttempts:    1,               // Unlimited
	}
	options := workflow.ActivityOptions{
		WaitForCancellation:    false,
		RetryPolicy:            retrypolicy,
		StartToCloseTimeout:    time.Second * 5,
		ScheduleToCloseTimeout: time.Second * 10,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	cwo := workflow.ChildWorkflowOptions{
		// Namespace: "notification", // TODO check with namespace
		TaskQueue:                NOTIFICATION_TASK_QUEUE,
		WorkflowExecutionTimeout: 10 * time.Minute,
		WorkflowTaskTimeout:      time.Minute,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	for _, step := range config.ApprovalSteps {
		if step.ConditionSpendRange.InRange(content.ExpenseAmountInCent) {

			// TODO what to close?
			// task-1: sync back expense approval start data
			err = workflow.ExecuteActivity(
				ctx, activities.CallbackStoreExpenseData, content).Get(ctx, nil)
			if err != nil {
				logger.Error(err.Error())
				return content, err
			}

			// task-2...n: loop for approver task-suite
			for _, ApproverID := range step.ApproverIDs {
				content.CurrentApproverID = ApproverID

				// task2-1: sync back approver task
				err = workflow.ExecuteActivity(
					ctx, activities.CallbackStoreApproverTask, content).Get(ctx, nil)
				if err != nil {
					logger.Error(err.Error())
					return content, err
				}

				// task2-2: wait for approver's submit
				logger.Debug("wait-1: approver with timeout")
				// submitSelector.Select(ctx)
				// TODO: seem can just call select and skip the await, but this may not be one state?
				isNotTimeout, err := workflow.AwaitWithTimeout(
					ctx, content.ApprovalExpiration, submitSelector.HasPending)
				if err != nil {
					logger.Error(err.Error())
					return content, err
				}
				if !isNotTimeout {
					logger.Debug("wait-2: expired, will wait on child workflow")
					content.ApprovalExpired = true
					// TODO new change of work, to see what happened, last time panic before 2-3
					err = workflow.ExecuteChildWorkflow(
						ctx, external.ExpenseNotificationWorkflow, "Maki").Get(ctx, nil)
					return content, err
				} else {
					// warning: if timeout, this channel stil blocks
					submitSelector.Select(ctx)
					logger.Debug("wait-2: get the submitter")
				}

				// task2-3: sync back approver result
				logger.Debug("lenth of content.ApprovalSubmits", len(content.ApprovalSubmits))
				latestSubmit := content.ApprovalSubmits[len(content.ApprovalSubmits)-1]
				err = workflow.ExecuteActivity(
					ctx, activities.CallbackUpdateApproverTask, latestSubmit).Get(ctx, nil)
				if err != nil {
					logger.Error(err.Error())
					return content, err
				}

				// step4: check if the flow should end,
				// other ending cases are omitted, take the reject only
				if content.ApprovalRejected {
					err = workflow.ExecuteChildWorkflow(
						ctx, external.ExpenseNotificationWorkflow, "Maki").Get(ctx, nil)
					return content, err
				}
			}
			break // only one branch applies, no need to rotate again
		}
	}

	err = workflow.ExecuteChildWorkflow(
		ctx, external.ExpenseNotificationWorkflow, "Fred").Get(ctx, nil)
	return content, err
}

// seems Team is the only index-key here to config
// their should be more indexings
// the configs can be the definitions shown to clients I suppose
func getWorkflowConfig(teamName model.TeamName) *model.ExpenseApprovalWorkflowConfig {
	return model.GetTheDemoConfig()
}

func checkValidityOfRequest(request *model.ExpenseApprovalStartRequest) error {
	if request == nil {
		return errors.New("request err: nil request")
	}
	if request.RequestorUserID <= 0 {
		return errors.New("request err: no userID")
	}
	if request.ExpenseAmountInCent <= 0 {
		return errors.New("request err: no valid amount")
	}
	if request.ApprovalExpiration <= 0 {
		return errors.New("request err: no valid expiration")
	}
	return nil
}

func createSubmitSelector(ctx workflow.Context, content *model.ExpenseApprovalContent) workflow.Selector {
	logger := workflow.GetLogger(ctx)
	submitSelector := workflow.NewSelector(ctx)
	submitChannel := workflow.GetSignalChannel(ctx, SINGALCHANNEL_SUBMIT)

	submitSelector.AddReceive(submitChannel, func(ch workflow.ReceiveChannel, more bool) {
		submitReult := &model.ApproverSubmitResult{}
		isChClosed := ch.Receive(ctx, submitReult)
		logger.Debug("[submit] called:", render.Render(submitReult), more)
		logger.Debug("[submit] channel closed?:", isChClosed)
		// other submits are omitted
		if submitReult.SubmitResult == model.SubmitReject {
			content.ApprovalCancelled = true
		}
		content.ApprovalSubmits = append(content.ApprovalSubmits, submitReult)
	})
	return submitSelector
}

func registerRuntimeQueries(ctx workflow.Context, content *model.ExpenseApprovalContent) error {
	return workflow.SetQueryHandler(ctx, QUERYNAME_CURRENAPPROVALCONTENT,
		func() (model.ExpenseApprovalContent, error) {
			return *content, nil
		})
}

// sub-workflows
func ExpenseApprovaSuccessWorkflow(ctx workflow.Context, content *model.ExpenseApprovalContent) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("ExpenseApprovaSuccessWorkflow enters")
	return nil
}

func ExpenseApprovaRejectedWorkflow(ctx workflow.Context, content *model.ExpenseApprovalContent) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("ExpenseApprovaRejectedWorkflow enters")
	return nil
}
