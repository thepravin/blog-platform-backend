package handlers

import (
	"blog_platform/internal/services"
	"blog_platform/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReactionHandler struct {
	Service *services.ReactionService
}

func NewReactionHandler(s *services.ReactionService) *ReactionHandler {
	return &ReactionHandler{Service: s}
}

type reactReq struct {
	Type string `json:"type"`
}

func (h *ReactionHandler) Toggle(c echo.Context) error {
	postID := c.Param("id")
	var req reactReq
	if err := c.Bind(&req); err != nil {
		return utils.Err(c, http.StatusBadRequest, "invalid payload")
	}
	uid := c.Get("user_id")
	if uid == nil {
		return utils.Err(c, http.StatusUnauthorized, "unauthorized")
	}
	liked, err := h.Service.TogglePostLike(postID, uid.(string))
	if err != nil {
		return utils.Err(c, http.StatusInternalServerError, err.Error())
	}
	if liked {
		return utils.JSON(c, http.StatusOK, true, "liked", nil)
	}
	return utils.JSON(c, http.StatusOK, true, "unliked", nil)
}
