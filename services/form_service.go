package services

import (
	"errors"
	"googleforms/internal/dto"
	"googleforms/internal/models"
	"googleforms/repositories"
	"math/rand"
	"time"
)

type formService struct {
	formRepo        repositories.FormRepository
	formVersionRepo repositories.FormVersionRepository
	tenantRepo      repositories.TenantRepository
}

func NewFormService(
	formRepo repositories.FormRepository,
	formVersionRepo repositories.FormVersionRepository,
	tenantRepo repositories.TenantRepository,
) FormService {
	return &formService{
		formRepo:        formRepo,
		formVersionRepo: formVersionRepo,
		tenantRepo:      tenantRepo,
	}
}

func (s *formService) CreateForm(tenantID uint, userID uint, req dto.CreateFormRequest) (*models.Form, error) {

	_, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return nil, errors.New("tenant not found")
	}

	// create the main form
	form := &models.Form{
		TenantID:    tenantID,
		CreatedBy:   &userID,
		Title:       req.Title,
		Description: req.Description,
		AllowGuest:  req.AllowGuest,
		PublicURL:   s.generatePublicURL(),
		Status:      "draft",
	}

	if err := s.formRepo.Create(form); err != nil {
		return nil, errors.New("failed to create form")
	}

	// the first version
	version := &models.FormVersion{
		FormID:           form.ID,
		VersionNumber:    1,
		Title:            req.Title,
		Description:      req.Description,
		Fields:           models.JSONB(req.Fields),
		SingleSubmission: req.SingleSubmission,
		CreatedBy:        &userID,
	}

	if err := s.formVersionRepo.Create(version); err != nil {
		// if the creation version faild, delete the form
		s.formRepo.Delete(tenantID, form.ID)
		return nil, errors.New("failed to create form version")
	}

	return form, nil
}

func (s *formService) GetForm(tenantID uint, formID uint) (*models.Form, error) {
	form, err := s.formRepo.GetByID(tenantID, formID)
	if err != nil {
		return nil, errors.New("form not found")
	}
	return form, nil
}

func (s *formService) GetPublicForm(publicURL string) (*dto.PublicFormResponse, error) {
	form, err := s.formRepo.GetByPublicURL(publicURL)
	if err != nil {
		return nil, errors.New("form not found")
	}

	// check the form are published
	if form.Status != "published" {
		return nil, errors.New("form not found")
	}
	// get the current version
	currentVersion, err := s.formVersionRepo.GetCurrentVersion(form.ID)
	if err != nil {
		return nil, errors.New("failed to get form version")
	}

	return &dto.PublicFormResponse{
		ID:          form.ID,
		Title:       currentVersion.Title,
		Description: currentVersion.Description,
		Fields:      currentVersion.Fields,
		Version:     currentVersion.VersionNumber,
	}, nil
}

func (s *formService) ListForms(tenantID uint, page, limit int) ([]models.Form, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.formRepo.ListByTenant(tenantID, page, limit)
}

func (s *formService) UpdateForm(tenantID uint, formID uint, userID uint, req dto.UpdateFormRequest) (*models.Form, error) {
	form, err := s.formRepo.GetByID(tenantID, formID)
	if err != nil {
		return nil, errors.New("form not found")
	}

	submissionCount, err := s.getSubmissionCount(formID)
	if err != nil {
		return nil, errors.New("failed to check form submissions")
	}

	currentVersion, err := s.formVersionRepo.GetCurrentVersion(formID)
	if err != nil {
		return nil, errors.New("failed to get current version")
	}

	if submissionCount > 0 {

		newVersionNumber := form.CurrentVersion + 1

		newVersion := &models.FormVersion{
			FormID:           form.ID,
			VersionNumber:    newVersionNumber,
			Title:            req.Title,
			Description:      req.Description,
			Fields:           models.JSONB(req.Fields),
			SingleSubmission: req.SingleSubmission,
			CreatedBy:        &userID,
		}

		if err := s.formVersionRepo.Create(newVersion); err != nil {
			return nil, errors.New("failed to create new version")
		}

		/// update the main version
		form.Title = req.Title
		form.Description = req.Description
		form.CurrentVersion = newVersionNumber
	} else {
		// if there are not answers, edition the main form
		currentVersion.Title = req.Title
		currentVersion.Description = req.Description
		currentVersion.Fields = models.JSONB(req.Fields)
		currentVersion.SingleSubmission = req.SingleSubmission

		// if err := s.formVersionRepo.Update(currentVersion); err != nil {
		// 	return nil, errors.New("failed to update current version")
		// }

		// update the same form
		form.Title = req.Title
		form.Description = req.Description
	}

	if err := s.formRepo.Update(form); err != nil {
		return nil, errors.New("failed to update form")
	}

	return form, nil
}

func (s *formService) PublishForm(tenantID uint, formID uint) (*models.Form, error) {
	form, err := s.formRepo.GetByID(tenantID, formID)
	if err != nil {
		return nil, errors.New("form not found")
	}

	if err := s.formRepo.UpdateStatus(tenantID, formID, "published"); err != nil {
		return nil, errors.New("failed to publish form")
	}

	form.Status = "published"
	return form, nil
}

func (s *formService) UnpublishForm(tenantID uint, formID uint) (*models.Form, error) {
	form, err := s.formRepo.GetByID(tenantID, formID)
	if err != nil {
		return nil, errors.New("form not found")
	}

	if err := s.formRepo.UpdateStatus(tenantID, formID, "draft"); err != nil {
		return nil, errors.New("failed to unpublish form")
	}

	form.Status = "draft"
	return form, nil
}

func (s *formService) DeleteForm(tenantID uint, formID uint) error {
	return s.formRepo.Delete(tenantID, formID)
}

func (s *formService) GetFormVersions(tenantID uint, formID uint) ([]models.FormVersion, error) {
	// التحقق من أن النموذج ينتمي للمستأجر
	_, err := s.formRepo.GetByID(tenantID, formID)
	if err != nil {
		return nil, errors.New("form not found")
	}

	return s.formVersionRepo.ListByForm(formID)
}

// ===== الدوال المساعدة =====

func (s *formService) getSubmissionCount(formID uint) (int64, error) {
	// هذه دالة مساعدة - تحتاج لـ submission repository
	// يمكن إضافتها لاحقاً
	return 0, nil
}

func (s *formService) generatePublicURL() string {
	rand.Seed(time.Now().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 10)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
