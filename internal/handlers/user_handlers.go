package handlers

import (
	"blog_platform/internal/errs"
	"blog_platform/internal/services"
	"blog_platform/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{
		Service: s,
	}
}

func (u *UserHandler) GetProfile(c echo.Context) error {
	userId, ok := c.Get("user_id").(string)
	if !ok {
		return errs.NewUnauthorizedError("Invalid or missing user id")
	}

	user, err := u.Service.GetProfile(userId)
	if err != nil {
		return err
	}

	return utils.JSON(c, http.StatusOK, true, "Profile found", user)
}
