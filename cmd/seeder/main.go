package main

import (
	"blog_platform/config"
	"blog_platform/internal/auth"
	"blog_platform/internal/database"
	"blog_platform/internal/models"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// Since we are running this locally, force host to localhost
	cfg.DBHost = "localhost"

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB(db)

	log.Println("Seeding database...")

	seedTags(db)
	seedUsers(db)
	seedPosts(db)
	seedComments(db)
	seedReactions(db)

	log.Println("Database seeded successfully!")
}

func seedUsers(db *gorm.DB) {
	hashedPassword, _ := auth.HashedPassword("password123")
	for i := 1; i <= 10; i++ {
		user := models.User{
			UserName:    fmt.Sprintf("user%d", i),
			Email:       fmt.Sprintf("user%d@example.com", i),
			Password:    hashedPassword,
			DisplayName: fmt.Sprintf("User %d", i),
			Bio:         fmt.Sprintf("This is the bio for user %d", i),
			Role:        "reader",
			IsActive:    true,
		}
		// ignore error if exists
		db.FirstOrCreate(&user, models.User{UserName: user.UserName})
	}
}

func seedTags(db *gorm.DB) {
	for i := 1; i <= 10; i++ {
		tag := models.Tag{
			Name: fmt.Sprintf("Tag%d", i),
			Slug: fmt.Sprintf("tag-%d", i),
		}
		db.FirstOrCreate(&tag, models.Tag{Name: tag.Name})
	}
}

func seedPosts(db *gorm.DB) {
	var users []models.User
	db.Find(&users)
	if len(users) == 0 {
		return
	}

	var tags []models.Tag
	db.Find(&tags)

	for i := 1; i <= 10; i++ {
		author := users[i%len(users)]
		post := models.Post{
			AuthorID: author.ID,
			Title:    fmt.Sprintf("Post Title %d", i),
			Slug:     fmt.Sprintf("post-title-%d-%s", i, uuid.NewString()[:6]),
			Excerpt:  fmt.Sprintf("This is an excerpt for post %d", i),
			Content:  fmt.Sprintf("This is the full content for post %d", i),
			Status:   "published",
			Visibility: "public",
			Views:    int64(i * 10),
		}
		if db.Where("title = ?", post.Title).First(&models.Post{}).Error == gorm.ErrRecordNotFound {
			db.Create(&post)
			// attach a couple of tags
			if len(tags) > 0 {
				db.Model(&post).Association("Tags").Append([]models.Tag{tags[i%len(tags)], tags[(i+1)%len(tags)]})
			}
		}
	}
}

func seedComments(db *gorm.DB) {
	var users []models.User
	db.Find(&users)
	var posts []models.Post
	db.Find(&posts)

	if len(users) == 0 || len(posts) == 0 {
		return
	}

	for i := 1; i <= 10; i++ {
		user := users[i%len(users)]
		post := posts[i%len(posts)]

		comment := models.Comment{
			PostID:     post.ID,
			UserID:     &user.ID,
			Content:    fmt.Sprintf("This is a comment %d on post", i),
			IsApproved: true,
		}
		db.Create(&comment)
	}
}

func seedReactions(db *gorm.DB) {
	var users []models.User
	db.Find(&users)
	var posts []models.Post
	db.Find(&posts)

	if len(users) == 0 || len(posts) == 0 {
		return
	}

	for i := 1; i <= 10; i++ {
		user := users[i%len(users)]
		post := posts[i%len(posts)]

		reaction := models.Reaction{
			PostID: post.ID,
			UserID: user.ID,
			Type:   "like",
		}
		db.FirstOrCreate(&reaction, models.Reaction{PostID: post.ID, UserID: user.ID})
	}
}
