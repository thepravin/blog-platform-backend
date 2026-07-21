package services

import (
	"blog_platform/internal/models"
	"blog_platform/internal/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentService struct {
	repo *repositories.CommentRepository
	db   *gorm.DB
}

func NewCommentService(r *repositories.CommentRepository, db *gorm.DB) *CommentService {
	return &CommentService{repo: r, db: db}
}

func (s *CommentService) Add(postID string, userID *string, parentId *string, content string) (*models.Comment, error) {
	var post models.Post
	if err := s.db.First(&post, "id=?", postID).Error; err != nil {
		return nil, err
	}
	// Proceed with the comment creation
	comment := &models.Comment{
		ID:        uuid.New(),
		PostID:    uuid.MustParse(postID),
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if userID != nil {
		uid := uuid.MustParse(*userID)
		comment.UserID = &uid
	}
	if parentId != nil && *parentId != "" {
		var parent models.Comment
		if err := s.db.First(&parent, "id=?", *parentId).Error; err != nil {
			return nil, err
		}
		pid := uuid.MustParse(*parentId)
		comment.ParentID = &pid
	}
	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	// Increment comment count in the post
	s.db.Model(&models.Post{}).Where("id = ?", postID).Update("comment_count", gorm.Expr("comment_count + ?", 1))

	return comment, nil
}

func (s *CommentService) ListByPost(postID string) ([]models.Comment, error) {
	return s.repo.ListByPost(postID)
}
