package storage

import (
	"context"
	"social/shared/models"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post models.Post) (models.PostID, error)
	GetPostByID(ctx context.Context, id models.PostID) (models.Post, error)
	UpdatePost(ctx context.Context, post models.Post) error
	GetPosts(ctx context.Context, offset, limit int) ([]models.Post, error)
	DeletePost(ctx context.Context, id models.PostID) error
}

type CommonRepository interface {
	PostRepository
}
