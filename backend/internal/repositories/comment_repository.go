package repositories

import (
	"blog_platform/internal/models"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(c *models.Comment) error { return r.db.Create(c).Error }

func (r *CommentRepository) ListByPost(postID string) ([]models.Comment, error) {
	var cs []models.Comment
	if err := r.db.Where("post_id=?", postID).Find(&cs).Error; err != nil {
		return nil, err
	}
	return cs, nil
}
