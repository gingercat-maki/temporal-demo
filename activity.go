// Since we don't use the common code as main entry directory
// we package it as the domain approval
package approval

import (
	"approval-workflow-demo-go/model"
	"context"
	"fmt"

	"github.com/luci/go-render/render"
	"go.temporal.io/sdk/activity"
)

// The grandularity of the approval business can be changed,
// but for demo we make it compact and simple
// TODO @maki find the right domain-names, check if the namings (approver/ticker...) should be changed
type ExpenseApprovalActivity struct {
}

// Callbacks to store domain data in the expense service
func (a *ExpenseApprovalActivity) CallbackStoreExpenseData(ctx context.Context, transferContent *model.ExpenseApprovalContent) error {
	message := fmt.Sprintf("[CallbackStoreExpenseData] to store expense data: %v", render.Render(transferContent))
	activity.GetLogger(ctx).Info(message)
	// return errors.New("try with some error")
	return nil
}

func (a *ExpenseApprovalActivity) CallbackStoreApproverTask(ctx context.Context, transferContent *model.ExpenseApprovalContent) error {
	message := fmt.Sprintf("[CallbackStoreApproverTask] current approver %v", transferContent.CurrentApproverID)
	activity.GetLogger(ctx).Info(message)
	return nil
}

func (a *ExpenseApprovalActivity) CallbackUpdateApproverTask(ctx context.Context, submitResult *model.ApproverSubmitResult) error {
	message := fmt.Sprintf("[CallbackUpdateApproverTask] approver result %v", submitResult)
	activity.GetLogger(ctx).Info(message)
	return nil
}

func (a *ExpenseApprovalActivity) SomeAsyncActivity(ctx context.Context, submitResult *model.ApproverSubmitResult) error {
	// Retrieve the Activity information needed to asynchronously complete the Activity.
	// activityInfo := activity.GetInfo(ctx)
	// taskToken := activityInfo.TaskToken

	// Send the taskToken to the external service that will complete the Activity.

	// Return from the Activity a function indicating that Temporal should wait for an async completion
	// message.
	return nil
}
