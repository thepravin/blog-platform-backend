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
		ID:         uuid.New(),
		AuthorID:   authorID,
		Title:      title,
		Slug:       utils.MakeSlugSimple(title + "-" + uuid.NewString()[:6]),
		Status:     "published",
		Visibility: "public",
		CreatedAt:  now,
		UpdatedAt:  now,
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

func (s *PostService) GetAll() ([]models.Post, error) { return s.repo.GetAll() }

func (s *PostService) GetAllDeletedPosts() ([]models.Post, error) { return s.repo.GetAllDeletedPosts() }

func (s *PostService) GetByID(id string) (*models.Post, error) {
	return s.repo.GetByID(id)
}

func (s *PostService) GetDeletedPostById(id string) (*models.Post, error) {
	return s.repo.GetDeletedPostById(id)
}

func (s *PostService) RestoreDeletedPostById(id string) error {
	return s.repo.RestoreDeletedPost(id)
}
