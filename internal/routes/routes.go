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
	e.Use(middleware.CORS(cfg))
	e.Use(middleware.InjectDB(db))

	// Repos or Services Call
	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo)
	userH := handlers.NewUserHandler(userSvc)
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
	public := e.Group("/api/v1")

	public.POST("/signup", authH.Signup)
	public.POST("/login", authH.Login)

	public.GET("/posts", postH.GetAll, middleware.OptionalMiddleware(cfg))
	public.GET("/posts/:slug", postH.GetPost, middleware.OptionalMiddleware(cfg))
	public.GET("/posts/:id/comments", commentH.List)
	public.POST("/posts/:id/view", postH.RecordView, middleware.OptionalMiddleware(cfg))

	public.GET("/tags", func(c echo.Context) error {
		db := c.Get("db").(*gorm.DB)
		var tags []models.Tag
		_ = db.Find(&tags).Error
		return utils.JSON(c, 200, true, "tags", tags)
	})

	// Protected Routes
	protected := e.Group("/api/v1")
	protected.Use(middleware.JWTMiddleware(cfg))

	protected.GET("/profile/me", userH.GetProfile)

	protected.POST("/posts", postH.CreatePost)
	protected.GET("/posts/history", postH.GetHistory)
	protected.GET("/posts/history/:id", postH.GetDeletedPostById)
	protected.GET("/posts/me", postH.GetPostsByUserId)
	protected.POST("/posts/:id/restore", postH.RestoreDeletedPostById)
	protected.DELETE("/posts/:id", postH.Delete)
	protected.PUT("/posts/:id/edit", postH.UpdatePost)

	protected.POST("/posts/:id/comments", commentH.Add)

	protected.POST("/posts/:id/reactions", reactionH.Toggle)
}
