package repositories

import (
	"blog_platform/internal/models"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}
func (r *PostRepository) Create(p *models.Post) error { return r.db.Create(p).Error }

func (r *PostRepository) GetAll() ([]models.Post, error) {
	var posts []models.Post

	err := r.db.Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *PostRepository) GetByID(id string) (*models.Post, error) {
	var post models.Post

	err := r.db.Preload("Comments").Preload("Reactions").Preload("Tags").Where("id=?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var post models.Post

		// Fetch the associated tags with the post
		if err := tx.Preload("Tags").First(&post, "id=?", id).Error; err != nil {
			return nil
		}
		// Delete all teh rows in post_tabs
		if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
			return err
		}

		// Delete associated comments
		if err := tx.Where("post_id = ?", id).Delete(&models.Comment{}).Error; err != nil {
			return err
		}

		// Delete associated reactions
		if err := tx.Where("post_id = ?", id).Delete(&models.Reaction{}).Error; err != nil {
			return err
		}

		// Delete Post
		if err := tx.Delete(&post).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *PostRepository) GetAllDeletedPosts() ([]models.Post, error) {
	var posts []models.Post

	err := r.db.Unscoped().Where("deleted_at IS NOT NULL").Find(&posts).Error

	return posts, err
}

func (r *PostRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}
