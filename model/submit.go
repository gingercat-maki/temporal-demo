package model

type SubmitAction int

const (
	SubmitReserved SubmitAction = iota
	SubmitApproval
	SubmitReject
	SubmitChange
	SubmitCancel
)

type ApproverSubmitResult struct {
	ApproverID   int64
	SubmitResult SubmitAction
	ChangeAmount int64
}
