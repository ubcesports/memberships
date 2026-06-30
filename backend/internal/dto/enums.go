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

type TransactionStatusType string

const (
	TransactionPending   TransactionStatusType = "pending"
	TransactionCompleted TransactionStatusType = "completed"
	TransactionFailed    TransactionStatusType = "failed"
	TransactionRefunded  TransactionStatusType = "refunded"
)

type PurchaseType string

const (
	PurchaseNew     PurchaseType = "new"
	PurchaseUpgrade PurchaseType = "upgrade"
)
