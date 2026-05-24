package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey"` // Use UUID instead of uint
    CreatedAt time.Time
    FullName  string         `json:"full_name"`
    Email     string         `json:"email" gorm:"unique"`
    Password  string         `json:"password"`
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}