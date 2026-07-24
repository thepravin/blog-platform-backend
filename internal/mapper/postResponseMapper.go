package mapper

import (
	"blog_platform/internal/dto"
	"blog_platform/internal/models"
)

func MapPostToResponse(post *models.Post) dto.PostResponse {
	return dto.PostResponse{
		ID:    post.ID.String(),
		Title: post.Title,
		Author: dto.AuthorType{
			ID:          post.Author.ID.String(),
			UserName:    post.Author.UserName,
			DisplayName: post.Author.DisplayName,
			Email:       post.Author.Email,
		},
		Content:       post.Content,
		ReactionCount: post.ReactionCount,
		CommentCount:  post.CommentCount,
		HasLiked:      post.HasLiked,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
		DeletedAt:     post.DeletedAt,
	}
}
