package repositories

import "googleforms/internal/models"

type TenantRepository interface {
	Create(tenant *models.Tenant) error
	GetByID(id uint) (*models.Tenant, error)
	GetBySlug(slug string) (*models.Tenant, error)
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByTenantAndEmail(tenantID uint, email string) (*models.User, error)
}

type FormRepository interface {
	Create(form *models.Form) error
	GetByID(tenantID uint, formID uint) (*models.Form, error)
	GetByPublicURL(publicURL string) (*models.Form, error)
	ListByTenant(tenantID uint, page, limit int) ([]models.Form, int64, error)
	Update(form *models.Form) error
	Delete(tenantID uint, formID uint) error
	UpdateStatus(tenantID uint, formID uint, status string) error
}

type FormVersionRepository interface {
	Create(version *models.FormVersion) error
	GetByID(id uint) (*models.FormVersion, error)
	GetByFormAndVersion(formID uint, versionNumber int) (*models.FormVersion, error)
	GetCurrentVersion(formID uint) (*models.FormVersion, error)
	ListByForm(formID uint) ([]models.FormVersion, error)
}

type SubmissionRepository interface {
	Create(submission *models.Submission) error
	GetByID(tenantID uint, submissionID uint) (*models.Submission, error)
	ListByForm(tenantID uint, formID uint, page, limit int) ([]models.Submission, int64, error)
	CountByForm(formID uint) (int64, error)
	CheckDuplicateSubmission(formID uint, email string) (bool, error)
}
