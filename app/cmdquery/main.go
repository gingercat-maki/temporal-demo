package main

import (
	"context"
	"flag"
	"log"

	approval "approval-workflow-demo-go"
	"approval-workflow-demo-go/model"

	"github.com/luci/go-render/render"
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
	flag.Parse()

	resp, err := c.QueryWorkflow(
		context.Background(), *workflowID, *runID, approval.QUERYNAME_CURRENAPPROVALCONTENT)
	if err != nil {
		log.Fatalln("Unable to query workflow", err)
	}
	var result model.ExpenseApprovalContent
	if err := resp.Get(&result); err != nil {
		log.Fatalln("Unable to decode query result", err)
	}
	log.Println("QueryExpenseApprovalWorkflow:", render.Render(result))
}
