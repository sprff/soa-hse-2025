package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"userservice/internal/api"
	"userservice/internal/models"

	"github.com/go-chi/chi/v5"
)

type RegisterUserRequest struct {
	Login    string
	Password string
	Email    string
}
type RegisterUserResponse struct {
	ID string `json:"id"`
}

func RegisterUser(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		input := RegisterUserRequest{}
		err = readBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}
		slog.InfoContext(ctx, "Create User", "login", input.Login)

		id, err := a.CreateUser(ctx, models.User{
			Email:    input.Email,
			Login:    input.Login,
			Password: input.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("can't create user: %w", err)
		}

		return RegisterUserResponse{ID: string(id)}, nil
	}
}

type AuthUserRequest struct {
	Login    string
	Password string
}
type AuthUserResponse struct {
	ID string `json:"id"`
}

func AuthUser(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		input := AuthUserRequest{}
		err = readBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}
		slog.InfoContext(ctx, "Auth User", "login", input.Login)

		id, err := a.AuthUser(ctx, input.Login, input.Password)
		if err != nil {
			return nil, fmt.Errorf("can't create user: %w", err)
		}

		return AuthUserResponse{ID: string(id)}, nil
	}
}

type UpdateUserRequest struct {
	Name        *string    `json:"name"`
	Surname     *string    `json:"surname"`
	DateOfBirth *time.Time `json:"dob"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone"`
	CreateDt    time.Time  `json:"created_at"`
	UpdateDt    time.Time  `json:"updated_at"`
}
type UpdateUserResponse struct{}

func UpdateUser(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		input := UpdateUserRequest{}
		err = readBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}
		id := models.UserID(chi.URLParam(r, "id"))
		slog.InfoContext(ctx, "UpdateUser", "id", id)

		err = a.UpdateUser(ctx, models.User{
			ID:          id,
			Name:        input.Name,
			Surname:     input.Surname,
			DateOfBirth: input.DateOfBirth,
			Email:       input.Email,
			Phone:       input.Phone,
			CreateDt:    input.CreateDt,
			UpdateDt:    input.UpdateDt,
		})

		if err != nil {
			return nil, fmt.Errorf("can't create user: %w", err)
		}

		return UpdateUserResponse{}, nil
	}
}

type GetUserByIDRequest struct{}
type GetUserByIDResponse struct {
	User models.User `json:"user"`
}

func GetUserByID(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		id := models.UserID(chi.URLParam(r, "id"))
		slog.InfoContext(ctx, "GetUserByID", "id", id)

		user, err := a.GetUserByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("can't get user by id: %w", err)
		}

		return GetUserByIDResponse{User: user}, nil
	}
}

type GetUserByLoginRequest struct{}
type GetUserByLoginResponse struct {
	User models.User `json:"user"`
}

func GetUserByLogin(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		login := chi.URLParam(r, "login")
		slog.InfoContext(ctx, "GetUserByLogin", "login", login)

		user, err := a.GetUserByLogin(ctx, login)
		if err != nil {
			return nil, fmt.Errorf("can't get user by login: %w", err)
		}

		return GetUserByIDResponse{User: user}, nil
	}
}
