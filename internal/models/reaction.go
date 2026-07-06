package models

import (
	"time"

	"github.com/google/uuid"
)

type Reaction struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"post_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Type      string    `gorm:"not null;default:'like'" json:"type"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Post Post `gorm:"foreignKey:PostID" json:"post"`
	User User `gorm:"foreignKey:UserID" json:"user"`
}
