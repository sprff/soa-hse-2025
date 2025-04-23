package api

import (
	"context"
	"fmt"
	"log"
	post "social/apiservice/internal/proto"
	"social/shared/models"
	userreqs "social/shared/models/requestmodels/userservicerequests"
	"social/shared/network"

	"github.com/IBM/sarama"
)

type Api struct {
	Usclient   UserserviceClient
	PostClient post.PostServiceClient
	Producer   sarama.SyncProducer
}

func NewApi(usClient UserserviceClient, postClient post.PostServiceClient) *Api {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll


	brokers := []string{"kafka:29092"}


	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Ошибка создания Producer: %v", err)
	}
	return &Api{
		Usclient:   usClient,
		PostClient: postClient,
		Producer:   producer,
	}
}

type UserserviceClient struct {
	Url string
}

func (c UserserviceClient) RegisterUser(ctx context.Context, in userreqs.RequestRegister) (out userreqs.ResponseRegister, err error) {
	err = network.MakeRequest(ctx, "POST", fmt.Sprintf("%s/users", c.Url), in, &out)
	return
}

func (c UserserviceClient) AuthUser(ctx context.Context, in userreqs.RequestAuth) (out userreqs.ResponseAuth, err error) {
	err = network.MakeRequest(ctx, "POST", fmt.Sprintf("%s/users/auth", c.Url), in, &out)
	return
}

func (c UserserviceClient) UpdateUser(ctx context.Context, id models.UserID, in userreqs.RequestUpdateUser) (out userreqs.ResponseUpdateUser, err error) {
	err = network.MakeRequest(ctx, "PUT", fmt.Sprintf("%s/users/%v", c.Url, id), in, &out)
	return
}

func (c UserserviceClient) GetUserByID(ctx context.Context, id models.UserID, in userreqs.RequestGetUserByID) (out userreqs.ResponseGetUserByID, err error) {
	err = network.MakeRequest(ctx, "GET", fmt.Sprintf("%s/users/%v", c.Url, id), in, &out)
	return
}
func (c UserserviceClient) GetUserByLogin(ctx context.Context, login string, in userreqs.RequestGetUserByLogin) (out userreqs.ResponseGetUserByLogin, err error) {
	err = network.MakeRequest(ctx, "GET", fmt.Sprintf("%s/users/bylogin/%s", c.Url, login), in, &out)
	return
}
