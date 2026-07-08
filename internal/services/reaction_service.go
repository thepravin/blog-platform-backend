package services

import (
	"blog_platform/internal/models"
	"blog_platform/internal/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReactionService struct {
	repo *repositories.ReactionRepository
	db   *gorm.DB
}

func NewReactionService(r *repositories.ReactionRepository, db *gorm.DB) *ReactionService {
	return &ReactionService{repo: r, db: db}
}

func (s *ReactionService) TogglePostLike(postID, userID string) (bool, error) {
	if re, err := s.repo.Find(postID, userID); err == nil && re != nil {
		if err := s.repo.Delete(re.ID.String()); err != nil {
			return false, err
		}
		return false, nil
	}

	nr := &models.Reaction{
		ID:        uuid.New(),
		PostID:    uuid.MustParse(postID),
		UserID:    uuid.MustParse(userID),
		Type:      "like",
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(nr); err != nil {
		return false, err
	}
	return true, nil
}
