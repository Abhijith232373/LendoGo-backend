package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JobApplication stores candidates applying for specific Career Openings.
type JobApplication struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	
	// 👇 THE FOREIGN KEY CONNECTION
	CareerOpeningID uuid.UUID      `gorm:"type:uuid;not null;index" json:"career_opening_id"`
	CareerOpening   CareerOpening  `gorm:"foreignKey:CareerOpeningID" json:"CareerOpening"`

	// Applicant Details
	FirstName       string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName        string         `gorm:"type:varchar(100);not null" json:"last_name"`
	Email           string         `gorm:"type:varchar(150);not null;index" json:"email"`
	Phone           string         `gorm:"type:varchar(20);not null" json:"phone"`
	Address         string         `gorm:"type:varchar(255)" json:"address"`
	City            string         `gorm:"type:varchar(100)" json:"city"`
	State           string         `gorm:"type:varchar(100)" json:"state"`
	PostalCode      string         `gorm:"type:varchar(20)" json:"postal_code"`

	// File Storage
	ResumePath      string         `gorm:"type:varchar(255);not null" json:"resume_path"` 
	
	// HR Tracking
	Status          string         `gorm:"type:varchar(50);default:'Under Review';index" json:"status"` 

	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}