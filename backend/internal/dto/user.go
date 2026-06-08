package dto

import "time"

// Update

type UpdateUserSelfDTO struct {
	FullName  *string `json:"full_name,omitempty"`
	IsStudent *bool   `json:"is_student,omitempty"`
	StudentID *string `json:"student_id,omitempty"`
}

type UpdateUserAdminDTO struct {
	FullName              *string    `json:"full_name,omitempty"`
	IsStudent             *bool      `json:"is_student,omitempty"`
	StudentID             *string    `json:"student_id,omitempty"`
	Role                  *RoleType  `json:"role,omitempty"`
	OnboardingCompletedAt *time.Time `json:"onboarding_completed_at,omitempty"`
}
