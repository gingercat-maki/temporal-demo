package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		log.Info("ping")
		c.JSON(200, gin.H{
			"message": "pong, hello world",
		})
	})
}

type ApprovalHandlers struct {
}

// TODO input-data, and input-response should specified
// input data should include
func (ah *ApprovalHandlers) StartApprovalProcess(ctx context.Context, data interface{}) error {
	log.Info("StartApprovalProcess called")
	return nil
}

// TODO for all queries, none of pagination, namespace isolation are set for demo simplicity
// TODO this should be changed into one general query, by status, by workflowID
func (ah *ApprovalHandlers) QueryApprovalProcessByStarter(ctx context.Context, userID int64) error {
	log.Info("StartApprovalProcess called")
	return nil
}

func (ah *ApprovalHandlers) QueryApprovalProcessByApprover(ctx context.Context, userID int64) error {
	log.Info("StartApprovalProcess called")
	return nil
}
