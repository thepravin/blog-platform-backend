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

	re, err := s.repo.Find(postID, userID)

	// record not present
	if err != nil {

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
		// Increment the reaction count
		s.db.Model(&models.Post{}).Where("id = ?", postID).Update("reaction_count", gorm.Expr("reaction_count + ?", 1))
		return true, nil
	}

	if re.DeletedAt.Valid { // current- unliked
		if err := s.repo.Restore(re.ID.String()); err != nil { // like
			return false, err
		}

		// Increment
		s.db.Model(&models.Post{}).Where("id = ?", postID).Update("reaction_count", gorm.Expr("reaction_count + ?", 1))
	} else { //current- liked
		if err := s.repo.Delete(re.ID.String()); err != nil {
			return false, err
		}
		// Decrement
		s.db.Model(&models.Post{}).Where("id = ?", postID).Update("reaction_count", gorm.Expr("reaction_count - ?", 1))
		return false, nil

	}

	return true, nil
}
