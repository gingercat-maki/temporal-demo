package model

import "math"

const (
	// These may be stored in database as metadata/config to run workflows
	// TODO this is another subdomain of problem, solve user/group/role problems later
	MockApproverID      = int64(111)
	MockSuperApproverID = int64(222)
	// should be in config
	NoApprovalNeededShreshold    = int64(1000)
	SuperApprovalNeededShreshold = int64(10000)
)

var (
	defaultMockApproversLevel1 = []int64{MockApproverID}
	defaultMockApproversLevel2 = []int64{MockApproverID, MockSuperApproverID}
)

// TODO if this is not for domain, but for a general purpose
// we are defining a workflow DSL ourselves
// for now its in accordance with the Frontend page
type ExpenseApprovalWorkflowConfig struct {
	ConditionTeamRequired TeamName
	ApprovalSteps         []*ExpenseApproveStepConfig
}

type ExpenseApproveStepConfig struct {
	// Branching Conditions
	ConditionSpendRange ExpenseAmountRange

	// Approvers: can be different rule to find one approver, here simplified to just userID
	ApproverIDs []int64 // simplified to userID, can be complicated
}

// just suppose [left, right)
type ExpenseAmountRange struct {
	Left  int64
	Right int64
}

func NewExpenseAmountRange(left int64, right int64) ExpenseAmountRange {
	return ExpenseAmountRange{left, right}
}

func (e ExpenseAmountRange) InRange(amount int64) bool {
	return amount >= e.Left && amount < e.Right
}

func GetTheDemoConfig() *ExpenseApprovalWorkflowConfig {

	firstLevel := &ExpenseApproveStepConfig{
		ConditionSpendRange: NewExpenseAmountRange(
			NoApprovalNeededShreshold, SuperApprovalNeededShreshold),
		ApproverIDs: defaultMockApproversLevel1,
	}

	secondLevel := &ExpenseApproveStepConfig{
		ConditionSpendRange: NewExpenseAmountRange(
			SuperApprovalNeededShreshold, math.MaxInt64),
		ApproverIDs: defaultMockApproversLevel2,
	}

	steps := []*ExpenseApproveStepConfig{firstLevel, secondLevel}
	config := &ExpenseApprovalWorkflowConfig{
		ConditionTeamRequired: TeamMarketing,
		ApprovalSteps:         steps,
	}

	return config
}
