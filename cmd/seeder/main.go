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

	profiles := []struct {
		UserName    string
		DisplayName string
		Email       string
	}{
		{"johndoe", "John Doe", "john.doe@example.com"},
		{"janesmith", "Jane Smith", "jane.smith@example.com"},
		{"alice_j", "Alice Johnson", "alice.j@example.com"},
		{"robert_b", "Robert Brown", "robert.b@example.com"},
		{"emilyw", "Emily White", "emily.white@example.com"},
		{"michael_t", "Michael Taylor", "michael.t@example.com"},
		{"sarah_d", "Sarah Davis", "sarah.davis@example.com"},
		{"david_m", "David Miller", "david.m@example.com"},
		{"laura_w", "Laura Wilson", "laura.w@example.com"},
		{"james_m", "James Moore", "james.moore@example.com"},
	}

	for _, p := range profiles {
		user := models.User{
			UserName:    p.UserName,
			Email:       p.Email,
			Password:    hashedPassword,
			DisplayName: p.DisplayName,
			Bio:         fmt.Sprintf("Hello! I am %s and I am a software enthusiast. I love exploring new technologies and writing about them.", p.DisplayName),
			Role:        "reader",
			IsActive:    true,
		}
		// ignore error if exists
		db.FirstOrCreate(&user, models.User{UserName: user.UserName})
	}
}

func seedTags(db *gorm.DB) {
	tagNames := []string{"Technology", "Programming", "Lifestyle", "Go", "Web Development", "Database", "Design", "Business", "Health", "React"}
	for _, name := range tagNames {
		tag := models.Tag{
			Name: name,
			Slug: fmt.Sprintf("tag-%s", name), // Note: A proper slug function is better, but this works for basic testing
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

	htmlContent := `
	<h1>Welcome to My Deep Dive!</h1>
	<p>This is a paragraph of text explaining some core concepts. In modern web development, it's very important to format our content using standard HTML tags so it renders beautifully on the frontend!</p>
	<h2>Key Takeaways</h2>
	<ul>
		<li>Understanding the basics is absolutely crucial.</li>
		<li>We can emphasize text using <strong>bold text</strong> or <em>italics</em>.</li>
		<li>Always ensure your architecture scales gracefully.</li>
	</ul>
	<blockquote>"Good code is its own best documentation."</blockquote>
	<p>Thank you for reading my article! Feel free to drop a reaction or leave a comment below sharing your thoughts.</p>
	`

	for i := 1; i <= 10; i++ {
		author := users[i%len(users)]
		post := models.Post{
			AuthorID:   author.ID,
			Title:      fmt.Sprintf("Exploring Modern Architecture %d", i),
			Slug:       fmt.Sprintf("exploring-modern-architecture-%d-%s", i, uuid.NewString()[:6]),
			Excerpt:    "A quick dive into modern concepts, sharing valuable insights on software design and structural patterns.",
			Content:    htmlContent,
			Status:     "published",
			Visibility: "public",
			Views:      int64(i * 15),
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

	commentTexts := []string{
		"This was an absolutely fantastic read, thank you!",
		"I completely agree with the points you made here.",
		"Could you elaborate more on the second point in a future post?",
		"Such great use of HTML styling! The frontend will look amazing.",
		"I learned so much from this.",
	}

	for i := 1; i <= 15; i++ {
		user := users[i%len(users)]
		post := posts[i%len(posts)]
		text := commentTexts[i%len(commentTexts)]

		comment := models.Comment{
			PostID:     post.ID,
			UserID:     &user.ID,
			Content:    text,
			IsApproved: true,
		}
		db.Create(&comment)
	}

	// Sync all comment counts
	db.Exec("UPDATE posts SET comment_count = (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id AND comments.deleted_at IS NULL)")
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

	// Sync all reaction counts
	db.Exec("UPDATE posts SET reaction_count = (SELECT COUNT(*) FROM reactions WHERE reactions.post_id = posts.id AND reactions.deleted_at IS NULL)")
}
