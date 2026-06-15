package dto

type OnboardUserDTO struct {
	IsStudent bool    `json:"isStudent"`
	StudentID *string `json:"studentID"`
}
