package main

import (
	approval "approval-workflow-demo-go"
	"approval-workflow-demo-go/model"
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"
)

const (
	userID        = int64(54321)
	amountInCents = int64(5000000)
	// expireTime    = 60 * time.Second
)

// thiy is actually the client for the workflow
// Create the client object just once per process
func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	options := client.StartWorkflowOptions{
		TaskQueue: approval.EXPENSE_APPROVAL_TASK_QUEUE,
	}
	transferApprovalContent := &model.ExpenseApprovalStartRequest{
		RequestorUserID:     userID,
		ExpenseAmountInCent: amountInCents,
		ExpenseTeam:         model.TeamMarketing,
		ApprovalExpiration:  1 * time.Minute,
	}
	we, err := c.ExecuteWorkflow(
		context.Background(), options, approval.ExpenseApprovalWorkflow, transferApprovalContent)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// c.CompleteActivity

	fmt.Println("Get workflow ID, and run ID:", we.GetID(), we.GetRunID())
}
