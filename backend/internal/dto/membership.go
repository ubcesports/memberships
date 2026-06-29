package dto

type MembershipTierDTO struct {
	ID          string                   `json:"id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Prices      []MembershipTierPriceDTO `json:"prices"`
}

type MembershipTierPriceDTO struct {
	Price             string `json:"price"`
	IsStudentRequired *bool  `json:"is_student_required"`
}
