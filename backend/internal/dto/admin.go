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

type AdminAuditLogActor struct {
	ActorUserId    string `json:"actor_user_id"`
	ActorFullName  string `json:"actor_full_name"`
	ActorAvatarURL string `json:"actor_avatar_url"`
}
