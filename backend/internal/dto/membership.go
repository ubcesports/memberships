package dto

// Create

type CreateMembershipDTO struct {
	UserID    string               `json:"user_id"`
	TierID    string               `json:"tier_id"`
	Status    MembershipStatusType `json:"status"`
	StartAt   string               `json:"started_at"`
	ExpiresAt string               `json:"expires_at"`
}

// Update

type UpdateMembershipAdminDTO struct {
	TierID    *string `json:"tier_id,omitempty"`
	Status    *string `json:"status,omitempty"`
	ExpiresAt *string `json:"expires_at,omitempty"`
}

// Read

type MembershipDTO struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	TierID    string `json:"tier_id"`
	Status    string `json:"status"`
	StartedAt string `json:"started_at"`
	ExpiresAt string `json:"expires_at"`
}
