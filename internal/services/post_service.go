package services

import (
	"blog_platform/internal/models"
	"blog_platform/internal/repositories"
	"blog_platform/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostService struct {
	repo *repositories.PostRepository
	db   *gorm.DB
}

func NewPostService(r *repositories.PostRepository, db *gorm.DB) *PostService {
	return &PostService{
		repo: r,
		db:   db,
	}
}
func (s *PostService) Create(authorID uuid.UUID, title, content string, tags []string) (*models.Post, error) {
	now := time.Now()

	post := &models.Post{
		ID:          uuid.New(),
		AuthorID:    authorID,
		Title:       title,
		Content:     content,
		Slug:        utils.MakeSlugSimple(title + "-" + uuid.NewString()[:6]),
		Status:      "published",
		Visibility:  "public",
		PublishedAt: &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		for _, tagName := range tags {
			var tag models.Tag
			if err := tx.Where("name=?", tagName).First(&tag).Error; err != nil {
				tag = models.Tag{
					ID:        uuid.New(),
					Name:      tagName,
					Slug:      utils.MakeSlugSimple(tagName),
					CreatedAt: now,
					UpdatedAt: now,
				}
				if err := tx.Create(&tag).Error; err != nil {
					return err
				}
			}
			if err := tx.Model(post).Association("Tags").Append(&tag); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) Delete(id string) error { return s.repo.Delete(id) }

func (s *PostService) Update(postID string, title, content string, tags []string) (*models.Post, error) {
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return nil, err
	}
	post.Title = title
	post.Content = content
	post.Slug = utils.MakeSlugSimple(title + "-" + uuid.NewString()[:6])
	post.UpdatedAt = time.Now()

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&post).Error; err != nil {
			return err
		}
		if len(tags) > 0 {
			var tagModels []models.Tag

			for _, tagName := range tags {
				var tag models.Tag
				if err := tx.Where("name=?", tagName).First(&tag).Error; err != nil {
					tag = models.Tag{
						ID:        uuid.New(),
						Name:      tagName,
						Slug:      utils.MakeSlugSimple(tagName),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
					if err := tx.Create(&tag).Error; err != nil {
						return err
					}
				}
				tagModels = append(tagModels, tag)
			}
			if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
				return err
			}
			if err := tx.Model(&post).Association("Tags").Append(tagModels); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetAll(sortParam string) ([]models.Post, error) {
	return s.repo.GetAll(sortParam)
}

func (s *PostService) GetAllDeletedPosts(id string) ([]models.Post, error) {
	return s.repo.GetAllDeletedPosts(id)
}

func (s *PostService) GetBySlug(slug, userID string) (*models.Post, error) {
	post, err := s.repo.GetByID(slug)
	if err != nil {
		return nil, err
	}

	if userID != "" {
		for _, reaction := range post.Reactions {
			if reaction.UserID.String() == userID {
				post.HasLiked = true
				break
			}
		}
	}

	return post, nil
}

func (s *PostService) GetAllByUserId(id string) ([]models.Post, error) {
	return s.repo.GetAllByUserId(id)
}

func (s *PostService) GetDeletedPostBySlug(slug string) (*models.Post, error) {
	return s.repo.GetDeletedPostBySlug(slug)
}

func (s *PostService) RestoreDeletedPostById(id string) error {
	return s.repo.RestoreDeletedPost(id)
}
func (s *PostService) RecordView(postID string, userID *string, ipAddress string) error {
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return err
	}

	var userUUIDPtr *uuid.UUID
	if userID != nil {
		u, err := uuid.Parse(*userID)
		if err == nil {
			userUUIDPtr = &u
		}
	}

	// Calculate the threshold for 24 hours ago
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)

	// Check if the user/IP has already viewed this post within the last 24 hours
	var existingView models.PostView
	query := s.db.Where("post_id = ? AND created_at > ?", postUUID, twentyFourHoursAgo)

	if userUUIDPtr != nil {
		query = query.Where("user_id = ?", userUUIDPtr)
	} else {
		query = query.Where("ip_address = ?", ipAddress)
	}

	err = query.First(&existingView).Error
	if err == nil {
		// A view already exists in the last 24 hours Do nothing.
		return nil
	} else if err != gorm.ErrRecordNotFound {
		return err // database error occurred
	}

	// No recent view found. Record a new view and update the post views count in a transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		newView := models.PostView{
			ID:        uuid.New(),
			PostID:    postUUID,
			UserID:    userUUIDPtr,
			IPAddress: ipAddress,
			CreatedAt: time.Now(),
		}

		if err := tx.Create(&newView).Error; err != nil {
			return err
		}

		return nil
	})
}
