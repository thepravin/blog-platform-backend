package routes

import (
	"blog_platform/config"
	"blog_platform/internal/handlers"
	"blog_platform/internal/middleware"
	"blog_platform/internal/models"
	"blog_platform/internal/repositories"
	"blog_platform/internal/services"
	"blog_platform/internal/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	e.Use(middleware.InjectDB(db))

	// Repos or Services Call
	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo)
	authH := handlers.NewAuthHandler(userSvc)

	postRepo := repositories.NewPostRepository(db)
	postSvc := services.NewPostService(postRepo, db)
	postH := handlers.NewPostHandler(postSvc)

	commentRepo := repositories.NewCommentRepository(db)
	commentSvc := services.NewCommentService(commentRepo, db)
	commentH := handlers.NewCommentHandler(commentSvc)

	reactionRepo := repositories.NewReactionRepository(db)
	reactionSvc := services.NewReactionService(reactionRepo, db)
	reactionH := handlers.NewReactionHandler(reactionSvc)

	// Public Endpoints
	e.POST("/api/v1/signup", authH.Signup)
	e.POST("/api/v1/login", authH.Login)

	e.GET("/api/v1/posts", postH.GetAll)
	e.GET("/api/v1/posts/:id", postH.GetPost)
	e.GET("/api/v1/posts/:id/comments", commentH.List)
	e.PUT("/api/v1/posts/:id", postH.UpdatePost)

	e.GET("/api/v1/tags", func(c echo.Context) error {
		db := c.Get("db").(*gorm.DB)
		var tags []models.Tag
		_ = db.Find(&tags).Error
		return utils.JSON(c, 200, true, "tags", tags)
	})

	// Protected Routes

	g := e.Group("/api/v1")
	g.Use(middleware.JWTMiddleware(cfg))
	g.POST("/posts", postH.CreatePost)
	g.GET("/posts/history", postH.GetHistory)
	g.DELETE("/posts/:id", postH.Delete)
	g.POST("/posts/:id/comments", commentH.Add)
	g.POST("/posts/:id/reactions", reactionH.Toggle)
}
