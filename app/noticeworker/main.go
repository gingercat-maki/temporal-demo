package main

import (
	"approval-workflow-demo-go"
	"approval-workflow-demo-go/external"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, approval.NOTIFICATION_TASK_QUEUE, worker.Options{})

	w.RegisterWorkflow(external.ExpenseNotificationWorkflow)
	w.RegisterActivity(&external.ExpenseNotificationActivity{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
