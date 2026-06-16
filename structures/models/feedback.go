package models

import (
	"time"
	"github.com/google/uuid"
)

type Feedback struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"` 
	Rating    int       `gorm:"not null" json:"rating"`        
	Comment   string    `gorm:"type:text" json:"comment"`
	Status    string    `gorm:"type:varchar(20);default:'Pending'" json:"status"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}