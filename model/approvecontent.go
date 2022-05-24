package model

import "time"

type ExpenseApprovalStartRequest struct {
	RequestorUserID     int64
	ExpenseAmountInCent int64
	ExpenseTeam         TeamName
	ApprovalExpiration  time.Duration
}

type ExpenseApprovalContent struct {

	// fill by starter content
	RequestorUserID     int64
	ExpenseAmountInCent int64
	ApprovalExpiration  time.Duration
	TeamMarketing       TeamName

	// fill by approvals
	CurrentApproverID int64
	ApprovalSubmits   []*ApproverSubmitResult

	// status (to be changed into one status)
	ApprovalCancelled bool
	ApprovalRejected  bool
	ApprovalExpired   bool

	Description string
}

func NewExpenseApprovalContent(request *ExpenseApprovalStartRequest) *ExpenseApprovalContent {

	content := &ExpenseApprovalContent{}
	content.RequestorUserID = request.RequestorUserID
	content.ExpenseAmountInCent = request.ExpenseAmountInCent
	content.ApprovalSubmits = make([]*ApproverSubmitResult, 0)
	content.ApprovalExpiration = request.ApprovalExpiration
	return content

}
