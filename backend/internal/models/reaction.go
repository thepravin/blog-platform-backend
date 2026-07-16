package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reaction struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PostID    uuid.UUID      `gorm:"type:uuid;not null" json:"post_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Type      string         `gorm:"not null;default:'like'" json:"type"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
