package repositories

import (
	"googleforms/internal/models"

	"gorm.io/gorm"
)

type formVersionRepository struct {
	db *gorm.DB
}

func NewFormVersionRepository(db *gorm.DB) FormVersionRepository {
	return &formVersionRepository{db: db}
}

func (r *formVersionRepository) Create(version *models.FormVersion) error {
	return r.db.Create(version).Error
}

func (r *formVersionRepository) GetByID(id uint) (*models.FormVersion, error) {
	var version models.FormVersion
	err := r.db.First(&version, id).Error
	return &version, err
}

func (r *formVersionRepository) GetByFormAndVersion(formID uint, versionNumber int) (*models.FormVersion, error) {
	var version models.FormVersion
	err := r.db.Where("form_id = ? AND version_number = ?", formID, versionNumber).First(&version).Error
	return &version, err
}

func (r *formVersionRepository) GetCurrentVersion(formID uint) (*models.FormVersion, error) {
	var form models.Form
	if err := r.db.First(&form, formID).Error; err != nil {
		return nil, err
	}

	var version models.FormVersion
	err := r.db.Where("form_id = ? AND version_number = ?", formID, form.CurrentVersion).First(&version).Error
	return &version, err
}

func (r *formVersionRepository) ListByForm(formID uint) ([]models.FormVersion, error) {
	var versions []models.FormVersion
	err := r.db.Where("form_id = ?", formID).Order("version_number DESC").Find(&versions).Error
	return versions, err
}

// func (r *formRepository) Update(oldForm *formRepository){

// }
