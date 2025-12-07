package models

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"match_id"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index;column:user_id" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"-"`

	ProductID   uuid.UUID   `gorm:"type:uuid;not null;index;column:product_id" json:"product_id"`
	LoanProduct LoanProduct `gorm:"constraint:OnDelete:CASCADE;foreignKey:ProductID" json:"-"`

	MatchConfidence bool      `gorm:"default:false" json:"match_confidence"`
	IsNotified      bool      `gorm:"default:false" json:"is_notified"`
	MatchedAt       time.Time `gorm:"autoCreateTime" json:"matched_at"`
	Reason          string    `gorm:"type:text" json:"reason"`
}
