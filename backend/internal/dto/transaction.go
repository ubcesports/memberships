package dto

// Create

type CreateTransactionDTO struct {
	UserID       string                `json:"user_id"`
	MembershipID string                `json:"membership_id"`
	PromoCodeID  *string               `json:"promo_code_id,omitempty"`
	PriceAmount  float64               `json:"price_amount"`
	Status       TransactionStatusType `json:"status"`
}

// Update

type UpdateTransactionDTO struct {
	Status *TransactionStatusType `json:"status,omitempty"`
}

// Read

type TransactionDTO struct {
	ID           string                `json:"id"`
	UserID       string                `json:"user_id"`
	MembershipID string                `json:"membership_id"`
	PromoCodeID  *string               `json:"promo_code_id,omitempty"`
	PriceAmount  float64               `json:"price_amount"`
	Status       TransactionStatusType `json:"status"`
	CreatedAt    string                `json:"created_at"`
	UpdatedAt    string                `json:"updated_at"`
}
