package models

import (
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CareerOpening struct {
	ID               uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title            string          `gorm:"type:varchar(150);not null" json:"title"`
	Department       string          `gorm:"type:varchar(100);not null" json:"department"`
	Location         string          `gorm:"type:varchar(150);not null" json:"location"`
	ExperienceRange  string          `gorm:"type:varchar(50);not null" json:"experience_range"`
	WorkMode         string          `gorm:"type:varchar(50);not null" json:"work_mode"`
	EmploymentType   string          `gorm:"type:varchar(50);not null" json:"employment_type"`
	Skills           json.RawMessage `gorm:"type:jsonb" json:"skills"`
	ShortDescription string          `gorm:"type:text;not null" json:"short_description"`
	AboutRole        string          `gorm:"type:text;not null" json:"about_role"`
	Responsibilities json.RawMessage `gorm:"type:jsonb" json:"responsibilities"`
	Requirements     json.RawMessage `gorm:"type:jsonb" json:"requirements"`
	Benefits         json.RawMessage `gorm:"type:jsonb" json:"benefits"`
	Status           string          `gorm:"type:varchar(20);default:'Open';index" json:"status"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}
