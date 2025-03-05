package storage

import (
	"context"
	"userservice/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.UserID, error)
	UpdateUser(ctx context.Context, user models.User) error
	GetUserByID(ctx context.Context, id models.UserID) (models.User, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}

type CommonRepository interface {
	UserRepository
}
