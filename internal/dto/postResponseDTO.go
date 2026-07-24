package dto

import (
	"time"

	"gorm.io/gorm"
)

type AuthorType struct {
	ID          string `json:"id"`
	UserName    string `json:"user_name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

type PostResponse struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Author        AuthorType     `json:"author"`
	Content       string         `json:"content"`
	CommentCount  int64          `json:"comment_count"`
	ReactionCount int64          `json:"reaction_count"`
	HasLiked      bool           `json:"has_liked"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}
