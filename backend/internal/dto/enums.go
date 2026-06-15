package dto

type RoleType string

const (
	RoleMember RoleType = "member"
	RoleAdmin  RoleType = "admin"
)

type GroupType string

const (
	GroupMember          GroupType = "member"
	GroupCompetitiveTeam GroupType = "competitive_team"
	GroupExecutive       GroupType = "executive"
	GroupDirector        GroupType = "director"
	GroupBoard           GroupType = "board"
)

type StudentStatusType string

const (
	StudentStatusStudent    StudentStatusType = "student"
	StudentStatusNonStudent StudentStatusType = "non_student"
)

type TierCodeType string

const (
	TierCodeRegular TierCodeType = "regular"
	TierCodePremium TierCodeType = "premium"
	TierCodeCab     TierCodeType = "cab"
	TierCodeDay     TierCodeType = "day"
)

type TransactionStatusType string

const (
	TransactionPending   TransactionStatusType = "pending"
	TransactionCompleted TransactionStatusType = "completed"
	TransactionFailed    TransactionStatusType = "failed"
	TransactionRefunded  TransactionStatusType = "refunded"
)
