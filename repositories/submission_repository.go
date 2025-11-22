package repositories

import (
	"googleforms/internal/models"

	"gorm.io/gorm"
)

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}

func (r *submissionRepository) Create(submission *models.Submission) error {
	return r.db.Create(submission).Error
}

func (r *submissionRepository) GetByID(tenantID uint, submissionID uint) (*models.Submission, error) {
	var submission models.Submission
	err := r.db.Where("id = ? AND tenant_id = ?", submissionID, tenantID).First(&submission).Error
	return &submission, err
}

func (r *submissionRepository) ListByForm(tenantID uint, formID uint, page, limit int) ([]models.Submission, int64, error) {
	var submissions []models.Submission
	var total int64

	offset := (page - 1) * limit

	r.db.Model(&models.Submission{}).Where("form_id = ? AND tenant_id = ?", formID, tenantID).Count(&total)

	err := r.db.Preload("FormVersion").
		Where("form_id = ? AND tenant_id = ?", formID, tenantID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&submissions).Error

	return submissions, total, err
}

func (r *submissionRepository) CountByForm(formID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Submission{}).Where("form_id = ?", formID).Count(&count).Error
	return count, err
}

func (r *submissionRepository) CheckDuplicateSubmission(formID uint, email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Submission{}).
		Where("form_id = ? AND guest_email = ?", formID, email).
		Count(&count).Error

	return count > 0, err
}
