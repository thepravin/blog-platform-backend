package handlers

import (
	"blog_platform/internal/services"
	"blog_platform/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Service *services.UserService
}

func NewAuthHandler(s *services.UserService) *AuthHandler {
	return &AuthHandler{
		Service: s,
	}
}

type SignupReq struct {
	UserName    string `json:"user_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

func (h *AuthHandler) Signup(c echo.Context) error {
	var req SignupReq
	if err := c.Bind(&req); err != nil {
		return utils.Err(c, http.StatusBadRequest, "Invalid Payload")
	}

	u, err := h.Service.Register(req.UserName, req.Email, req.Password)
	if err != nil {
		return utils.Err(c, http.StatusBadRequest, err.Error())
	}

	return utils.JSON(c, http.StatusCreated, true, "user_created", u)
}

type LoginReq struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return utils.Err(c, http.StatusBadRequest, "Invalid Payload")
	}

	user, token, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		return utils.Err(c, http.StatusBadRequest, "Invalid Email or Password")
	}

	return utils.JSON(
		c,
		http.StatusOK,
		true,
		"Login Successfull",
		map[string]interface{}{
			"token": token,
			"user":  user,
		},
	)
}
