package main

import (
	"context"
	"flag"
	"log"

	approval "approval-workflow-demo-go"

	"approval-workflow-demo-go/model"

	"go.temporal.io/sdk/client"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowID := flag.String("workflowID", "", "workflowID, required")
	runID := flag.String("runID", "", "runID, optional")
	isReject := flag.Bool("isReject", false, "optional")
	flag.Parse()

	request := model.ApproverSubmitResult{
		SubmitResult: model.SubmitApproval,
	}
	if *isReject {
		request.SubmitResult = model.SubmitReject
	}
	err = c.SignalWorkflow(
		context.Background(), *workflowID, *runID, approval.SINGALCHANNEL_SUBMIT, request)
	if err != nil {
		log.Fatalln("Unable to signal workflow", err)
	}
}
