package approval

const (
	// These may be stored in database as metadata/config to run workflows
	// TODO this is another subdomain of problem, solve user/group/role problems later
	MockApproverID      = int64(123456)
	MockSuperApproverID = int64(999999)

	// should be in config
	NoApprovalNeededShreshold    = int64(1000)
	SuperApprovalNeededShreshold = int64(10000)
)

// permission, user logic in the approver
func GetDirectApprover() int64 {
	return MockApproverID
}

func GetSecondLevelApprover() int64 {
	return MockSuperApproverID
}

// the amount for approval process
func IsNeedApproval(amountInCent int64) bool {
	return amountInCent > NoApprovalNeededShreshold
}

func IsNeedSuperApproval(amountInCent int64) bool {
	return amountInCent > SuperApprovalNeededShreshold
}
