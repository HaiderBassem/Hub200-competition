package services

import (
	"googleforms/internal/dto"
	"googleforms/internal/models"
)

type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	ValidateToken(token string) (*models.User, error)
}

type FormService interface {
	CreateForm(tenantID uint, userID uint, req dto.CreateFormRequest) (*models.Form, error)
	GetForm(tenantID uint, formID uint) (*models.Form, error)
	GetPublicForm(publicURL string) (*dto.PublicFormResponse, error)
	ListForms(tenantID uint, page, limit int) ([]models.Form, int64, error)
	UpdateForm(tenantID uint, formID uint, userID uint, req dto.UpdateFormRequest) (*models.Form, error)
	PublishForm(tenantID uint, formID uint) (*models.Form, error)
	UnpublishForm(tenantID uint, formID uint) (*models.Form, error)
	DeleteForm(tenantID uint, formID uint) error
	GetFormVersions(tenantID uint, formID uint) ([]models.FormVersion, error)
}

type SubmissionService interface {
	SubmitForm(publicURL string, req dto.SubmitFormRequest) (*models.Submission, error)
	GetSubmissions(tenantID uint, formID uint, page, limit int) ([]models.Submission, int64, error)
	GetSubmissionStats(tenantID uint, formID uint) (*dto.SubmissionStats, error)
}
