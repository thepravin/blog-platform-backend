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

	err := r.db.
		Preload("Author").
		Preload("Comments").
		Preload("Reactions").
		Preload("Tags").
		Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *PostRepository) GetAllByUserId(id string) ([]models.Post, error) {
	var posts []models.Post

	err := r.db.
		Preload("Tags").
		Where("author_id=?", id).
		Find(&posts).
		Error

	return posts, err
}

func (r *PostRepository) GetByID(slug string) (*models.Post, error) {
	var post models.Post

	err := r.db.
		Preload("Author").
		Preload("Comments").
		Preload("Reactions").
		Preload("Tags").
		Where("slug=?", slug).
		First(&post).
		Error

	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var post models.Post

		// Fetch the post to ensure it exists
		if err := tx.First(&post, "id=?", id).Error; err != nil {
			return err
		}

		// Delete associated tags in the join table
		if err := tx.Where("post_id = ?", id).Delete(&models.PostTag{}).Error; err != nil {
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

func (r *PostRepository) GetAllDeletedPosts(id string) ([]models.Post, error) {
	var posts []models.Post

	err := r.db.Unscoped().Where("author_id=? and deleted_at IS NOT NULL", id).Find(&posts).Error

	return posts, err
}

func (r *PostRepository) GetDeletedPostById(id string) (*models.Post, error) {
	var post models.Post

	err := r.db.Unscoped().
		Where("deleted_at IS NOT NULL").
		Preload("Author").
		Preload("Comments").
		Preload("Tags").
		Preload("Reactions").
		Where("id = ?", id).
		First(&post).
		Error

	if err != nil {
		return nil, err
	}

	return &post, err
}

func (r *PostRepository) RestoreDeletedPost(id string) error {

	return r.db.Transaction(func(tx *gorm.DB) error {

		// restore posts
		if err := tx.Unscoped().Model(&models.Post{}).Where("id =?", id).Update("deleted_at", nil).Error; err != nil {
			return err
		}

		// restore comments
		if err := tx.Unscoped().Model(&models.Comment{}).Where("post_id = ?", id).Update("deleted_at", nil).Error; err != nil {
			return err
		}

		// restore reactions
		if err := tx.Unscoped().Model(&models.Reaction{}).Where("post_id = ?", id).Update("deleted_at", nil).Error; err != nil {
			return err
		}

		// restore tags (in the post_tags join table)
		if err := tx.Unscoped().Model(&models.PostTag{}).Where("post_id = ?", id).Update("deleted_at", nil).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PostRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}
