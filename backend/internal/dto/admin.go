package dto

import "time"

type AdminAuditLogResponse struct {
	Actor       AdminAuditLogActor       `json:"actor"`
	OccuredAt   time.Time                `json:"occured_at"`
	Action      string                   `json:"action"`
	Description *string                  `json:"description"`
	Outcome     AdminAuditLogOutcomeType `json:"outcome"`
	RequestId   string                   `json:"request_id"`
	TargetUser  *AdminAuditLogActor      `json:"target_user"`
}

// AdminUpdateUserRequest is the body of PATCH /admin/users/{id}. Every field is
// optional and only the ones present are applied.
type AdminUpdateUserRequest struct {
	StudentID        *string     `json:"student_id"`
	IsStudent        *bool       `json:"is_student"`
	GroupsAdd        []GroupType `json:"groups_add"`
	GroupsRemove     []GroupType `json:"groups_remove"`
	Role             *RoleType   `json:"role"`
	CancelMembership bool        `json:"cancel_membership"`
}

type AdminAuditLogActor struct {
	ActorUserId    string `json:"actor_user_id"`
	ActorFullName  string `json:"actor_full_name"`
	ActorAvatarURL string `json:"actor_avatar_url"`
}
