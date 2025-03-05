package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"userservice/internal/models"
)

func (a *Api) CreateUser(ctx context.Context, user models.User) (models.UserID, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	slog.Info("API CreateUser")

	if err := user.Validate(); err != nil {
		return models.UserID(""), fmt.Errorf("can't validate user: %w", err)
	}

	id, err := a.storage.CreateUser(ctx, user)
	if err != nil {
		return models.UserID(""), fmt.Errorf("can't write user to storage: %w", err)
	}
	return id, nil
}

func (a *Api) AuthUser(ctx context.Context, login string, password string) (id models.UserID, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	slog.Info("API AuthUser")

	expectedUser := models.User{
		Login:    login,
		Password: password,
	}
	expectedUser.HidePassword()

	user, err := a.storage.GetUserByLogin(ctx, login)
	switch {
	case errors.Is(err, models.ErrUserNotFound):
		return id, nil
	case err != nil:
		return id, err
	}

	if user.Password != expectedUser.Password {
		return id, nil
	}

	id = user.ID
	return
}

func (a *Api) UpdateUser(ctx context.Context, user models.User) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	slog.Info("API UpdateUser")

	if err := user.Validate(); err != nil {
		return fmt.Errorf("can't validate user: %w", err)
	}

	err := a.storage.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("can't write user to storage: %w", err)
	}
	return nil
}

func (a *Api) GetUserByID(ctx context.Context, id models.UserID) (models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	slog.Info("API GetUserByID")

	user, err := a.storage.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("can't get user from storage: %w", err)
	}
	return user, nil
}

func (a *Api) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	slog.Info("API GetUserByLogin")

	user, err := a.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return models.User{}, fmt.Errorf("can't get user from storage: %w", err)
	}
	slog.Debug("API GET User", "user", user)
	return user, nil
}
