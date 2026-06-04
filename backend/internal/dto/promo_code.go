package dto

// Create

type CreatePromoCodeDTO struct {
	Code               string    `json:"code"`
	DiscountPercentage int       `json:"discount_percentage"`
	IsActive           bool      `json:"is_active"`
	RoleOverride       *RoleType `json:"role_override,omitempty"`
}

// Update

type UpdatePromoCodeDTO struct {
	Code               *string `json:"code,omitempty"`
	IsActive           *bool   `json:"is_active,omitempty"`
	DiscountPercentage *int    `json:"discount_percentage,omitempty"`
	RoleOverride       *string `json:"role_override,omitempty"`
}

// Read

type PromoCodeDTO struct {
	ID                 string `json:"id"`
	Code               string `json:"code"`
	DiscountPercentage int    `json:"discount_percentage"`
	IsActive           bool   `json:"is_active"`
	RoleOverride       string `json:"role_override"`
}
