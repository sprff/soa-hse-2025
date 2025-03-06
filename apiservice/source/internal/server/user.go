package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"social/apiservice/internal/api"
	"social/shared/models"
	apireqs "social/shared/models/requestmodels/apiservicerequests"
	userreqs "social/shared/models/requestmodels/userservicerequests"
	"social/shared/network"

	"github.com/go-chi/chi/v5"
)

func RegisterUser(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		input := apireqs.RequestRegister{}
		err = network.ReadBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}
		slog.InfoContext(ctx, "Create User", "login", input.Login)

		out, err := a.Usclient.RegisterUser(ctx, userreqs.RequestRegister{
			Login:    input.Login,
			Password: input.Password,
			Email:    input.Email,
		})
		if err != nil {
			return nil, fmt.Errorf("can't register user: %w", err)
		}
		return apireqs.ResponseRegister{ID: out.ID}, nil
	}
}

func AuthUser(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		input := apireqs.RequestAuth{}
		err = network.ReadBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}
		slog.InfoContext(ctx, "Auth User", "login", input.Login)

		out, err := a.Usclient.AuthUser(ctx, userreqs.RequestAuth{Login: input.Login, Password: input.Password})
		if err != nil {
			return nil, fmt.Errorf("can't auth user: %w", err)
		}

		return apireqs.ResponseAuth{ID: out.ID}, nil
	}
}

func UpdateUser(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		input := apireqs.RequestUpdateUser{}
		err = network.ReadBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}
		id := models.UserID(chi.URLParam(r, "id"))
		slog.InfoContext(ctx, "UpdateUser", "id", id)

		out, err := a.Usclient.UpdateUser(ctx, id, userreqs.RequestUpdateUser{
			Name:        input.Name,
			Surname:     input.Surname,
			DateOfBirth: input.DateOfBirth,
			Email:       input.Email,
			Phone:       input.Phone,
		})
		if err != nil {
			return nil, fmt.Errorf("can't update user: %w", err)
		}

		return apireqs.ResponseUpdateUser{
			ID:          out.ID,
			Login:       out.Login,
			Name:        out.Name,
			Surname:     out.Surname,
			DateOfBirth: out.DateOfBirth,
			Email:       out.Email,
			Phone:       out.Phone,
			CreateDt:    out.CreateDt,
			UpdateDt:    out.UpdateDt,
		}, nil
	}
}

func GetUserByID(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		id := models.UserID(chi.URLParam(r, "id"))
		slog.InfoContext(ctx, "GetUserByID", "id", id)

		out, err := a.Usclient.GetUserByID(ctx, id, userreqs.RequestGetUserByID{})
		if err != nil {
			return nil, fmt.Errorf("can't get user by id: %w", err)
		}

		return apireqs.ResponseGetUserByID{User: out.User}, nil
	}
}

func GetUserByLogin(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		login := chi.URLParam(r, "login")
		slog.InfoContext(ctx, "GetUserByLogin", "login", login)

		out, err := a.Usclient.GetUserByLogin(ctx, login, userreqs.RequestGetUserByLogin{})
		if err != nil {
			return nil, fmt.Errorf("can't get user by id: %w", err)
		}

		return apireqs.ResponseGetUserByLogin{User: out.User}, nil
	}
}
