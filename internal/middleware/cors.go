package middleware

import (
	"blog_platform/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CORS(cfg *config.Config) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH},
	})
}
