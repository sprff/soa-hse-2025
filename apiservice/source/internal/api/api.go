package api

import (
	"context"
	"fmt"
	"social/shared/models"
	userreqs "social/shared/models/requestmodels/userservicerequests"
	"social/shared/network"
)

type Api struct {
	Usclient UserserviceClient
}

func NewApi() *Api {
	return &Api{
		Usclient: UserserviceClient{"http://userservice-backend:8080"},
	}
}

type UserserviceClient struct {
	url string
}

func (c UserserviceClient) RegisterUser(ctx context.Context, in userreqs.RequestRegister) (out userreqs.ResponseRegister, err error) {
	err = network.MakeRequest(ctx, "POST", fmt.Sprintf("%s/users", c.url), in, &out)
	return
}

func (c UserserviceClient) AuthUser(ctx context.Context, in userreqs.RequestAuth) (out userreqs.ResponseAuth, err error) {
	err = network.MakeRequest(ctx, "POST", fmt.Sprintf("%s/users/auth", c.url), in, &out)
	return
}

func (c UserserviceClient) UpdateUser(ctx context.Context, id models.UserID, in userreqs.RequestUpdateUser) (out userreqs.ResponseUpdateUser, err error) {
	err = network.MakeRequest(ctx, "PUT", fmt.Sprintf("%s/users/%v", c.url, id), in, &out)
	return
}

func (c UserserviceClient) GetUserByID(ctx context.Context, id models.UserID, in userreqs.RequestGetUserByID) (out userreqs.ResponseGetUserByID, err error) {
	err = network.MakeRequest(ctx, "GET", fmt.Sprintf("%s/users/%v", c.url, id), in, &out)
	return
}
func (c UserserviceClient) GetUserByLogin(ctx context.Context, login string, in userreqs.RequestGetUserByLogin) (out userreqs.ResponseGetUserByLogin, err error) {
	err = network.MakeRequest(ctx, "GET", fmt.Sprintf("%s/users/bylogin/%s", c.url, login), in, &out)
	return
}
