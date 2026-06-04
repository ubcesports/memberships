package dto

// Create

type CreateVerificationTokenDTO struct {
	UserID    string                `json:"user_id"`
	Type      VerificationTokenType `json:"type"`
	ExpiresAt string                `json:"expires_at"`
}

// Read

type VerificationTokenDTO struct {
	ID        string                `json:"id"`
	UserID    string                `json:"user_id"`
	Type      VerificationTokenType `json:"type"`
	ExpiresAt string                `json:"expires_at"`
	CreatedAt string                `json:"created_at"`
}
