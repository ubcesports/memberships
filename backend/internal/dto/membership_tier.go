package dto

// Create

type CreateMembershipTierDTO struct {
	Title           string  `json:"title"`
	Description     *string `json:"description,omitempty"`
	PriceStudent    float64 `json:"price_student"`
	PriceNonStudent float64 `json:"price_non_student"`
}

// Update

type UpdateMembershipTierDTO struct {
	Title           *string  `json:"title,omitempty"`
	Description     *string  `json:"description,omitempty"`
	PriceStudent    *float64 `json:"price_student,omitempty"`
	PriceNonStudent *float64 `json:"price_non_student,omitempty"`
}

// Read

type MembershipTierDTO struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Description     *string `json:"description,omitempty"`
	PriceStudent    float64 `json:"price_student"`
	PriceNonStudent float64 `json:"price_non_student"`
}
