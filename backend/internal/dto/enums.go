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
	GroupCentralDirector GroupType = "central_director"
	GroupGameDirector    GroupType = "game_director"
	GroupBoard           GroupType = "board"
	GroupPresident       GroupType = "president"
)

type TransactionStatusType string

const (
	TransactionPending   TransactionStatusType = "pending"
	TransactionCompleted TransactionStatusType = "completed"
	TransactionFailed    TransactionStatusType = "failed"
	TransactionRefunded  TransactionStatusType = "refunded"
	TransactionExpired   TransactionStatusType = "expired"
)

type PurchaseType string

const (
	PurchaseNew     PurchaseType = "new"
	PurchaseUpgrade PurchaseType = "upgrade"
)

type AdminAuditLogOutcomeType string

const (
	AuditLogSuccess AdminAuditLogOutcomeType = "success"
	AuditLogFailed  AdminAuditLogOutcomeType = "failed"
	AuditLogDenied  AdminAuditLogOutcomeType = "denied"
)

type ExecSocialPlatformType string

const (
	ExecSocialPlatformInstagram ExecSocialPlatformType = "instagram"
	ExecSocialPlatformX         ExecSocialPlatformType = "x"
	ExecSocialPlatformTwitch    ExecSocialPlatformType = "twitch"
	ExecSocialPlatformYoutube   ExecSocialPlatformType = "youtube"
	ExecSocialPlatformTiktok    ExecSocialPlatformType = "tiktok"
	ExecSocialPlatformLinkedIn  ExecSocialPlatformType = "linkedin"
)
