package handlers

import (
	"blog_platform/internal/errs"
	"blog_platform/internal/mapper"
	"blog_platform/internal/services"
	"blog_platform/internal/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	Service *services.PostService
}

func NewPostHandler(s *services.PostService) *PostHandler { return &PostHandler{Service: s} }

type CreatePostReq struct {
	Title   string   `json:"title,omitempty"`
	Content string   `json:"content,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	var req CreatePostReq
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}
	userIDstr, ok := c.Get("user_id").(string)
	if !ok {
		return errs.NewUnauthorizedError("Invalid or missing user id")
	}
	authorID, err := uuid.Parse(userIDstr)
	if err != nil {
		return errs.NewUnauthorizedError("Invalid user id format")
	}
	post, err := h.Service.Create(authorID, req.Title, req.Content, req.Tags)
	if err != nil {
		return err
	}
	return utils.JSON(c, http.StatusCreated, true, "Post created", post)
}

func (h *PostHandler) GetAll(c echo.Context) error {
	sortParam := c.QueryParam("sort")
	if sortParam == "" {
		sortParam = "latest"
	}

	posts, err := h.Service.GetAll(sortParam)
	if err != nil {
		return err
	}
	return utils.JSON(c, http.StatusOK, true, "posts", posts)
}

func (h *PostHandler) GetHistory(c echo.Context) error {
	userId, ok := c.Get("user_id").(string)

	if !ok {
		return errs.NewUnauthorizedError("Invalid user id")
	}

	posts, err := h.Service.GetAllDeletedPosts(userId)
	if err != nil {
		return err
	}

	return utils.JSON(c, http.StatusOK, true, "posts", posts)
}

func (h *PostHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.Service.Delete(id); err != nil {
		return err
	}
	return utils.JSON(c, http.StatusOK, true, "deleted", nil)
}

func (h *PostHandler) UpdatePost(c echo.Context) error {
	postID := c.Param("id")
	var req CreatePostReq
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}
	updatePost, err := h.Service.Update(postID, req.Title, req.Content, req.Tags)
	if err != nil {
		return err
	}
	return utils.JSON(c, http.StatusOK, true, "Post updated", updatePost)
}

func (h *PostHandler) GetPost(c echo.Context) error {
	slug := c.Param("slug")
	userId, _ := c.Get("user_id").(string)
	post, err := h.Service.GetBySlug(slug, userId)
	if err != nil {
		return err
	}

	safeResponse := mapper.MapPostToResponse(post)

	return utils.JSON(c, http.StatusOK, true, "Post fetched", safeResponse)
}

func (h *PostHandler) GetPostsByUserId(c echo.Context) error {
	userId, ok := c.Get("user_id").(string)

	if !ok {
		return errs.NewUnauthorizedError("Invalid user id")
	}

	posts, err := h.Service.GetAllByUserId(userId)

	if err != nil {
		return err
	}

	return utils.JSON(c, http.StatusOK, true, "posts", posts)
}

func (h *PostHandler) GetDeletedPostById(c echo.Context) error {
	id := c.Param("id")
	post, err := h.Service.GetDeletedPostById(id)
	if err != nil {
		return err
	}

	return utils.JSON(c, http.StatusOK, true, "Post fetched", post)
}

func (h *PostHandler) RestoreDeletedPostById(c echo.Context) error {
	id := c.Param("id")
	if err := h.Service.RestoreDeletedPostById(id); err != nil {
		return err
	}
	return utils.JSON(c, http.StatusOK, true, "Post restored", nil)
}

func (h *PostHandler) RecordView(c echo.Context) error {
	postID := c.Param("id")

	// 1. Get the IP Address for anonymous fallback
	ipAddress := c.RealIP()

	// 2. Try to get the UserID (If the OptionalAuth middleware found a token)
	var userID *string
	if u := c.Get("user_id"); u != nil {
		uid := u.(string)
		userID = &uid
	}

	// 3. Call the service to record the view
	err := h.Service.RecordView(postID, userID, ipAddress)
	if err != nil {
		return err
	}

	return utils.JSON(c, http.StatusOK, true, "View recorded", nil)
}
