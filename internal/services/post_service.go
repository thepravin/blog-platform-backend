package services

import (
	"blog_platform/internal/repositories"

	"gorm.io/gorm"
)

type PostService struct {
	repo *repositories.PostRepository
	db   *gorm.DB
}

func NewPostService(r *repositories.PostRepository, db *gorm.DB) *PostService {
	return &PostService{repo: r, db: db}
}
