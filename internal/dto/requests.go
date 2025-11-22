package dto

// Auth Requests
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	TenantID uint   `json:"tenant_id" validate:"required"`
}

// Form Requests
type CreateFormRequest struct {
	Title            string                 `json:"title" validate:"required"`
	Description      string                 `json:"description"`
	Fields           map[string]interface{} `json:"fields" validate:"required"`
	AllowGuest       bool                   `json:"allow_guest"`
	SingleSubmission bool                   `json:"single_submission"`
}

type UpdateFormRequest struct {
	Title            string                 `json:"title" validate:"required"`
	Description      string                 `json:"description"`
	Fields           map[string]interface{} `json:"fields" validate:"required"`
	SingleSubmission bool                   `json:"single_submission"`
}

// Submission Requests
type SubmitFormRequest struct {
	Answers map[string]interface{} `json:"answers" validate:"required"`
	Email   string                 `json:"email,omitempty"`
}
