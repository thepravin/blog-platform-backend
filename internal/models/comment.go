package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PostID     uuid.UUID  `gorm:"type:uuid;not null" json:"post_id"`
	UserID     *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	ParentID   *uuid.UUID `gorm:"type:uuid" json:"parent_id"`
	Content    string     `gorm:"type:text;not null" json:"content"`
	IsApproved bool       `gorm:"default:true" json:"is_approved"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
