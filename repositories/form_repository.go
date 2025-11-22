package repositories

import (
	"googleforms/internal/models"

	"gorm.io/gorm"
)

type formRepository struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) FormRepository {
	return &formRepository{db: db}
}

func (r *formRepository) Create(form *models.Form) error {
	return r.db.Create(form).Error
}

func (r *formRepository) GetByID(tenantID uint, formID uint) (*models.Form, error) {
	var form models.Form
	err := r.db.Where("id = ? AND tenant_id = ?", formID, tenantID).First(&form).Error
	return &form, err
}

func (r *formRepository) GetByPublicURL(publicURL string) (*models.Form, error) {
	var form models.Form
	err := r.db.Where("public_url = ?", publicURL).First(&form).Error
	return &form, err
}

func (r *formRepository) ListByTenant(tenantID uint, page, limit int) ([]models.Form, int64, error) {
	var forms []models.Form
	var total int64

	offset := (page - 1) * limit

	r.db.Model(&models.Form{}).Where("tenant_id = ?", tenantID).Count(&total)

	err := r.db.Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&forms).Error

	return forms, total, err
}

func (r *formRepository) Update(form *models.Form) error {
	return r.db.Save(form).Error
}

func (r *formRepository) Delete(tenantID uint, formID uint) error {
	return r.db.Where("id = ? AND tenant_id = ?", formID, tenantID).Delete(&models.Form{}).Error
}

func (r *formRepository) UpdateStatus(tenantID uint, formID uint, status string) error {
	return r.db.Model(&models.Form{}).
		Where("id = ? AND tenant_id = ?", formID, tenantID).
		Update("status", status).Error
}
