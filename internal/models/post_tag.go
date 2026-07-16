package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostTag struct {
	PostID    uuid.UUID      `gorm:"primaryKey"`
	TagID     uuid.UUID      `gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
