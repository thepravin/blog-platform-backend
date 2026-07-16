package handlers

import (
	"blog_platform/internal/services"
	"blog_platform/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	Service *services.CommentService
}

func NewCommentHandler(s *services.CommentService) *CommentHandler {
	return &CommentHandler{Service: s}
}

type addCommentReq struct {
	Content  string  `json:"content,omitempty"`
	ParentID *string `json:"parent_id,omitempty"`
}

func (h *CommentHandler) Add(c echo.Context) error {
	postID := c.Param("id")
	var req addCommentReq
	if err := c.Bind(&req); err != nil {
		return utils.Err(c, http.StatusBadRequest, "invalid payload")
	}
	var userID *string
	if u := c.Get("user_id"); u != nil {
		uid := u.(string)
		userID = &uid
	}
	comment, err := h.Service.Add(postID, userID, req.ParentID, req.Content)
	if err != nil {
		return utils.Err(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSON(c, http.StatusCreated, true, "comment added", comment)

}

func (h *CommentHandler) List(c echo.Context) error {
	postID := c.Param("id")
	comment, err := h.Service.ListByPost(postID)
	if err != nil {
		return utils.Err(c, http.StatusInternalServerError, err.Error())

	}
	return utils.JSON(c, http.StatusOK, true, "comments", comment)
}
