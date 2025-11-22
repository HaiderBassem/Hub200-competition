package dto

import "time"

// Base Responses
type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Auth Responses
type LoginResponse struct {
	Success bool         `json:"success"`
	Token   string       `json:"token"`
	User    UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	TenantID  uint      `json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Form Responses
type FormResponse struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	CurrentVersion int       `json:"current_version"`
	AllowGuest     bool      `json:"allow_guest"`
	RequireLogin   bool      `json:"require_login"`
	Status         string    `json:"status"`
	PublicURL      string    `json:"public_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type FormVersionResponse struct {
	ID               uint                   `json:"id"`
	VersionNumber    int                    `json:"version_number"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Fields           map[string]interface{} `json:"fields"`
	SingleSubmission bool                   `json:"single_submission"`
	CreatedAt        time.Time              `json:"created_at"`
}

type PublicFormResponse struct {
	ID          uint                   `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Fields      map[string]interface{} `json:"fields"`
	Version     int                    `json:"version"`
}

// Submission Responses
type SubmissionResponse struct {
	ID            uint                   `json:"id"`
	FormID        uint                   `json:"form_id"`
	FormVersionID uint                   `json:"form_version_id"`
	UserID        *uint                  `json:"user_id,omitempty"`
	IsGuest       bool                   `json:"is_guest"`
	GuestEmail    string                 `json:"guest_email,omitempty"`
	Answers       map[string]interface{} `json:"answers"`
	CreatedAt     time.Time              `json:"created_at"`
}

type SubmissionStats struct {
	TotalSubmissions int64 `json:"total_submissions"`
	GuestSubmissions int64 `json:"guest_submissions"`
	UserSubmissions  int64 `json:"user_submissions"`
}
