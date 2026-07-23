package middleware

import (
	"blog_platform/config"
	"blog_platform/internal/auth"
	"strings"

	"github.com/labstack/echo/v4"
)

func OptionalMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return next(c)
			}

			if !strings.HasPrefix(authHeader, "Bearer") {
				return next(c)
			}

			token := strings.TrimPrefix(authHeader, "Bearer")
			claims, err := auth.ValidateJWT(token)
			if err != nil {
				return next(c)
			}
			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)
			return next(c)
		}
	}
}
