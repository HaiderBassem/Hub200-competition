package repositories

import (
	"googleforms/internal/models"

	"gorm.io/gorm"
)

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(tenant *models.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *tenantRepository) GetByID(id uint) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.First(&tenant, id).Error
	return &tenant, err
}

func (r *tenantRepository) GetBySlug(slug string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.Where("slug = ?", slug).First(&tenant).Error
	return &tenant, err
}
