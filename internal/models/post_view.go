package models

import (
	"time"

	"github.com/google/uuid"
)

type PostView struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PostID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	UserID    *uuid.UUID `gorm:"type:uuid;index"`
	IPAddress string     `gorm:"index"`
	Processed bool       `gorm:"default:false"`
	CreatedAt time.Time
}
