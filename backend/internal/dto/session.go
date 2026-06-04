package dto

// Create

type CreateSessionDTO struct {
	UserID    string  `json:"user_id"`
	ExpiresAt string  `json:"expires_at"`
	IPAddress *string `json:"ip_address,omitempty"`
	UserAgent *string `json:"user_agent,omitempty"`
}

// Update

type UpdateSessionDTO struct {
	LastSeenAt *string `json:"last_seen_at,omitempty"`
	IPAddress  *string `json:"ip_address,omitempty"`
}

// Read

type SessionDTO struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	CreatedAt  string  `json:"created_at"`
	ExpiresAt  string  `json:"expires_at"`
	LastSeenAt *string `json:"last_seen_at,omitempty"`
	IPAddress  *string `json:"ip_address,omitempty"`
	UserAgent  *string `json:"user_agent,omitempty"`
}
