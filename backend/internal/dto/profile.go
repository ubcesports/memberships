package dto

import "time"

type ProfileDTO struct {
	ID                    string      `json:"id"`
	Email                 string      `json:"email"`
	StudentID             *string     `json:"student_id"`
	Role                  RoleType    `json:"role"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
	FullName              string      `json:"full_name"`
	EmailVerifiedAt       *time.Time  `json:"email_verified_at"`
	IsStudent             bool        `json:"is_student"`
	OnboardingCompletedAt *time.Time  `json:"onboarding_completed_at"`
	AvatarURL             *string     `json:"avatar_url"`
	Groups                []GroupType `json:"groups"`
}

type OnboardUserRequest struct {
	IsStudent bool    `json:"is_student"`
	StudentID *string `json:"student_id"`
}

type AdminAuditLogResponse struct {
	Actor       AdminAuditLogActor       `json:"actor"`
	OccuredAt   time.Time                `json:"occured_at"`
	Action      string                   `json:"action"`
	Description *string                  `json:"description"`
	Outcome     AdminAuditLogOutcomeType `json:"outcome"`
	RequestId   string                   `json:"request_id"`
	TargetUser  *AdminAuditLogActor      `json:"target_user"`
}

type AdminAuditLogActor struct {
	ActorUserId    string `json:"actor_user_id"`
	ActorFullName  string `json:"actor_full_name"`
	ActorAvatarURL string `json:"actor_avatar_url"`
}
