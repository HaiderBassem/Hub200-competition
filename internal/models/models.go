package models

import (
	"time"
)

// JSONB type for PostgreSQL JSONB
type JSONB map[string]interface{}

type Tenant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`

	Users []User `gorm:"foreignKey:TenantID" json:"users,omitempty"`
	Forms []Form `gorm:"foreignKey:TenantID" json:"forms,omitempty"`
}

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TenantID     uint      `gorm:"not null" json:"tenant_id"`
	Email        string    `gorm:"not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	FullName     string    `json:"full_name"`
	Role         string    `gorm:"default:editor" json:"role"`
	CreatedAt    time.Time `gorm:"default:now()" json:"created_at"`

	// Relations
	Tenant      Tenant       `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Forms       []Form       `gorm:"foreignKey:CreatedBy" json:"forms,omitempty"`
	Submissions []Submission `gorm:"foreignKey:UserID" json:"submissions,omitempty"`
}

type Form struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	TenantID       uint      `gorm:"not null" json:"tenant_id"`
	CreatedBy      *uint     `json:"created_by,omitempty"`
	Title          string    `gorm:"not null" json:"title"`
	Description    string    `json:"description"`
	CurrentVersion int       `gorm:"default:1" json:"current_version"`
	AllowGuest     bool      `gorm:"default:true" json:"allow_guest"`
	RequireLogin   bool      `gorm:"default:false" json:"require_login"`
	Status         string    `gorm:"default:draft" json:"status"`
	PublicURL      string    `gorm:"uniqueIndex" json:"public_url"`
	CreatedAt      time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:now()" json:"updated_at"`

	// Relations
	Tenant      Tenant        `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Versions    []FormVersion `gorm:"foreignKey:FormID" json:"versions,omitempty"`
	Submissions []Submission  `gorm:"foreignKey:FormID" json:"submissions,omitempty"`
}

type FormVersion struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	FormID           uint      `gorm:"not null" json:"form_id"`
	VersionNumber    int       `gorm:"not null" json:"version_number"`
	Title            string    `gorm:"not null" json:"title"`
	Description      string    `json:"description"`
	Fields           JSONB     `gorm:"type:jsonb;not null" json:"fields"`
	SingleSubmission bool      `gorm:"default:false" json:"single_submission"`
	CreatedBy        *uint     `json:"created_by,omitempty"`
	CreatedAt        time.Time `gorm:"default:now()" json:"created_at"`

	// Relations
	Form        Form         `gorm:"foreignKey:FormID" json:"form,omitempty"`
	Submissions []Submission `gorm:"foreignKey:FormVersionID" json:"submissions,omitempty"`
}

type Submission struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TenantID      uint      `gorm:"not null" json:"tenant_id"`
	FormID        uint      `gorm:"not null" json:"form_id"`
	FormVersionID uint      `gorm:"not null" json:"form_version_id"`
	UserID        *uint     `json:"user_id,omitempty"`
	IsGuest       bool      `gorm:"default:false" json:"is_guest"`
	GuestEmail    string    `json:"guest_email,omitempty"`
	Answers       JSONB     `gorm:"type:jsonb;not null" json:"answers"`
	CreatedAt     time.Time `gorm:"default:now()" json:"created_at"`

	// Relations
	Form        Form        `gorm:"foreignKey:FormID" json:"form,omitempty"`
	FormVersion FormVersion `gorm:"foreignKey:FormVersionID" json:"form_version,omitempty"`
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
