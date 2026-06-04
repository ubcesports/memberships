package dto

type RoleType string

const (
	RoleMember      RoleType = "member"
	RoleExec        RoleType = "exec"
	RoleCompetitive RoleType = "competitive"
	RoleAdmin       RoleType = "admin"
)

type VerificationTokenType string

const (
	VerificationEmail VerificationTokenType = "email_verification"
	VerificationMagic VerificationTokenType = "magic_link"
)

type MembershipStatusType string

const (
	MembershipActive   MembershipStatusType = "active"
	MembershipExpired  MembershipStatusType = "expired"
	MembershipCancelled MembershipStatusType = "cancelled"
)

type TransactionStatusType string

const (
	TransactionPending   TransactionStatusType = "pending"
	TransactionCompleted TransactionStatusType = "completed"
	TransactionFailed    TransactionStatusType = "failed"
	TransactionRefunded  TransactionStatusType = "refunded"
)