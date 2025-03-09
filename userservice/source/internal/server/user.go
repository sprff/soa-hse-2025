package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"social/shared/models"
	"social/shared/models/requestmodels/userservicerequests"
	"social/userservice/internal/api"

	"github.com/go-chi/chi/v5"
)

func RegisterUser(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		input := userservicerequests.RequestRegister{}
		err = readBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}
		slog.InfoContext(ctx, "Create User", "login", input.Login)

		id, err := a.CreateUser(ctx, input.Login, input.Password, input.Email)
		if err != nil {
			return nil, fmt.Errorf("can't create user: %w", err)
		}

		return userservicerequests.ResponseRegister{ID: string(id)}, nil
	}
}

func AuthUser(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		input := userservicerequests.RequestAuth{}
		err = readBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}
		slog.InfoContext(ctx, "Auth User", "login", input.Login)

		id, err := a.AuthUser(ctx, input.Login, input.Password)
		if err != nil {
			return nil, fmt.Errorf("can't create user: %w", err)
		}

		return userservicerequests.ResponseAuth{ID: string(id)}, nil
	}
}

func UpdateUser(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		input := userservicerequests.RequestUpdateUser{}
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
		})
		if err != nil {
			return nil, fmt.Errorf("can't update user: %w", err)
		}

		newUser, err := a.GetUserByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("can't get user user: %w", err)
		}
		return userservicerequests.ResponseUpdateUser{
			ID:          string(id),
			Login:       newUser.Login,
			Name:        newUser.Name,
			Surname:     newUser.Surname,
			DateOfBirth: newUser.DateOfBirth,
			Email:       newUser.Email,
			Phone:       newUser.Phone,
			CreateDt:    newUser.CreateDt,
			UpdateDt:    newUser.UpdateDt}, nil
	}
}

func GetUserByID(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		id := models.UserID(chi.URLParam(r, "id"))
		slog.InfoContext(ctx, "GetUserByID", "id", id)

		user, err := a.GetUserByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("can't get user by id: %w", err)
		}

		return userservicerequests.ResponseGetUserByID{User: user}, nil
	}
}

func GetUserByLogin(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		login := chi.URLParam(r, "login")
		slog.InfoContext(ctx, "GetUserByLogin", "login", login)

		user, err := a.GetUserByLogin(ctx, login)
		if err != nil {
			return nil, fmt.Errorf("can't get user by login: %w", err)
		}

		return userservicerequests.ResponseGetUserByLogin{User: user}, nil
	}
}
