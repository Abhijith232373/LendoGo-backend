package models

import (
	"time"
	"github.com/google/uuid"
)

type SystemWallet struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WalletName string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"wallet_name"`
	Balance    float64   `gorm:"not null;default:0.0" json:"balance"`
	UpdatedAt  time.Time `json:"updated_at"`
}