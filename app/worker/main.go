package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	approval "approval-workflow-demo-go"
	"approval-workflow-demo-go/external"
)

func main() {

	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, approval.EXPENSE_APPROVAL_TASK_QUEUE, worker.Options{})

	// warning, should register child workflows
	w.RegisterWorkflow(approval.ExpenseApprovalWorkflow)
	w.RegisterWorkflow(external.ExpenseNotificationWorkflow)

	w.RegisterActivity(&approval.ExpenseApprovalActivity{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
