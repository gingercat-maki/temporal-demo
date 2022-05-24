// external represents some other workflow providers that can be directly used
package external

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

type ExpenseNotificationActivity struct {
}

func (a *ExpenseNotificationActivity) SendEmail(ctx context.Context, name string) error {
	message := fmt.Sprintf("ExpenseNotificationActivity-SendEmail to %s", name)
	activity.GetLogger(ctx).Info(message)
	return nil
}

var (
	noticeActivities = &ExpenseNotificationActivity{}
)

// some other standard workflows can be directly used
func ExpenseNotificationWorkflow(ctx workflow.Context, name string) error {
	// TODO add names, to try parallel, see how workflow manage parallel states
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	return workflow.ExecuteActivity(
		ctx, noticeActivities.SendEmail, name).Get(ctx, nil)
}
