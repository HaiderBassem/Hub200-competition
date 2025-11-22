package services

import (
	"errors"
	"googleforms/internal/dto"
	"googleforms/internal/models"
	"googleforms/repositories"
)

type submissionService struct {
	submissionRepo  repositories.SubmissionRepository
	formRepo        repositories.FormRepository
	formVersionRepo repositories.FormVersionRepository
}

func NewSubmissionService(
	submissionRepo repositories.SubmissionRepository,
	formRepo repositories.FormRepository,
	formVersionRepo repositories.FormVersionRepository,
) SubmissionService {
	return &submissionService{
		submissionRepo:  submissionRepo,
		formRepo:        formRepo,
		formVersionRepo: formVersionRepo,
	}
}

func (s *submissionService) SubmitForm(publicURL string, req dto.SubmitFormRequest) (*models.Submission, error) {

	form, err := s.formRepo.GetByPublicURL(publicURL)
	if err != nil {
		return nil, errors.New("form not found")
	}

	// check the form are published
	if form.Status != "published" {
		return nil, errors.New("form is not published")
	}

	currentVersion, err := s.formVersionRepo.GetCurrentVersion(form.ID)
	if err != nil {
		return nil, errors.New("failed to get form version")
	}

	if currentVersion.SingleSubmission && req.Email != "" {
		isDuplicate, err := s.submissionRepo.CheckDuplicateSubmission(form.ID, req.Email)
		if err != nil {
			return nil, errors.New("failed to check duplicate submission")
		}
		if isDuplicate {
			return nil, errors.New("you have already submitted this form")
		}
	}

	submission := &models.Submission{
		TenantID:      form.TenantID,
		FormID:        form.ID,
		FormVersionID: currentVersion.ID,
		Answers:       models.JSONB(req.Answers),
		IsGuest:       true,
		GuestEmail:    req.Email,
	}

	if err := s.submissionRepo.Create(submission); err != nil {
		return nil, errors.New("failed to submit form")
	}

	return submission, nil
}

func (s *submissionService) GetSubmissions(tenantID uint, formID uint, page, limit int) ([]models.Submission, int64, error) {

	_, err := s.formRepo.GetByID(tenantID, formID)
	if err != nil {
		return nil, 0, errors.New("form not found")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.submissionRepo.ListByForm(tenantID, formID, page, limit)
}

func (s *submissionService) GetSubmissionStats(tenantID uint, formID uint) (*dto.SubmissionStats, error) {
	// check the form releated t o tentant
	_, err := s.formRepo.GetByID(tenantID, formID)
	if err != nil {
		return nil, errors.New("form not found")
	}

	total, err := s.submissionRepo.CountByForm(formID)
	if err != nil {
		return nil, errors.New("failed to get submission stats")
	}

	return &dto.SubmissionStats{
		TotalSubmissions: total,
	}, nil
}
