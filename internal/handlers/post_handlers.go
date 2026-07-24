package handlers

import (
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalidrequest"})
	}
	userIDstr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid user id"})
	}
	authorID, err := uuid.Parse(userIDstr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid user id"})
	}
	post, err := h.Service.Create(authorID, req.Title, req.Content, req.Tags)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetAll(c echo.Context) error {
	posts, err := h.Service.GetAll()
	if err != nil {
		return utils.Err(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSON(c, http.StatusOK, true, "posts", posts)
}

func (h *PostHandler) GetHistory(c echo.Context) error {
	userId, ok := c.Get("user_id").(string)

	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.APIResponse{
			Success: false,
			Message: "Invalid user id",
		})
	}

	posts, err := h.Service.GetAllDeletedPosts(userId)
	if err != nil {
		return utils.Err(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSON(c, http.StatusOK, true, "posts", posts)
}

func (h *PostHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.Service.Delete(id); err != nil {
		return utils.Err(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSON(c, http.StatusOK, true, "deleted", nil)
}

func (h *PostHandler) UpdatePost(c echo.Context) error {
	postID := c.Param("id")
	var req CreatePostReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	updatePost, err := h.Service.Update(postID, req.Title, req.Content, req.Tags)
	if err != nil {
		return utils.Err(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, updatePost)
}

func (h *PostHandler) GetPost(c echo.Context) error {
	slug := c.Param("slug")
	userId, _ := c.Get("user_id").(string)
	post, err := h.Service.GetBySlug(slug, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.APIResponse{
			Success: false,
			Message: "Post Not found",
		})
	}

	safeResponse := mapper.MapPostToResponse(post)

	return c.JSON(http.StatusOK, utils.APIResponse{
		Success: true,
		Message: "Post fetched ",
		Data:    safeResponse,
	})
}

func (h *PostHandler) GetPostsByUserId(c echo.Context) error {
	userId, ok := c.Get("user_id").(string)

	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.APIResponse{
			Success: false,
			Message: "Invalid user id",
		})
	}

	posts, err := h.Service.GetAllByUserId(userId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.APIResponse{
			Success: false,
			Message: "Internal server Error",
		})
	}

	return utils.JSON(c, http.StatusOK, true, "posts", posts)
}

func (h *PostHandler) GetDeletedPostById(c echo.Context) error {
	id := c.Param("id")
	post, err := h.Service.GetDeletedPostById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.APIResponse{
			Success: false,
			Message: "Post Not found",
		})
	}

	return c.JSON(http.StatusOK, utils.APIResponse{
		Success: true,
		Message: "Post fetched",
		Data:    post,
	})
}

func (h *PostHandler) RestoreDeletedPostById(c echo.Context) error {
	id := c.Param("id")
	if err := h.Service.RestoreDeletedPostById(id); err != nil {
		return utils.Err(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSON(c, http.StatusOK, true, "Post restored", nil)
}
