package dto

// Create

type CreateUserDTO struct {
	Email        string  `json:"email"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	UbcStudentID *string `json:"ubc_student_id,omitempty"`
}

// Update

type UpdateUserSelfDTO struct {
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	UbcStudentID *string `json:"ubc_student_id,omitempty"`
}

type UpdateUserAdminDTO struct {
	FirstName    *string   `json:"first_name,omitempty"`
	LastName     *string   `json:"last_name,omitempty"`
	UbcStudentID *string   `json:"ubc_student_id,omitempty"`
	Role         *RoleType `json:"role,omitempty"`
	IsVerified   *bool     `json:"is_verified,omitempty"`
}

// Read

type UserDTO struct {
	ID           string   `json:"id"`
	Email        string   `json:"email"`
	IsVerified   bool     `json:"is_verified"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	UbcStudentID *string  `json:"ubc_student_id,omitempty"`
	Role         RoleType `json:"role"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}
