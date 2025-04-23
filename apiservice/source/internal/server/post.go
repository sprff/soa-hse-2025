package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"social/apiservice/internal/api"
	post "social/apiservice/internal/proto"
	"social/shared/models"
	kevents "social/shared/models/kafka_events"
	apireqs "social/shared/models/requestmodels/apiservicerequests"
	"social/shared/models/requestmodels/userservicerequests"
	"social/shared/network"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
)

func CreatePost(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		input := apireqs.RequestCreatePost{}
		err = network.ReadBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}
		slog.InfoContext(ctx, "Create Post", "title", input.Title)

		userResp, err := a.Usclient.AuthUser(ctx, userservicerequests.RequestAuth{
			Login:    input.Login,
			Password: input.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("can't auth user: %w", err)
		}

		resp, err := a.PostClient.CreatePost(ctx, &post.CreatePostRequest{
			Title:       input.Title,
			Description: input.Description,
			CreatorId:   userResp.ID,
			IsPrivate:   input.IsPrivate,
			Tags:        input.Tags,
		})
		if err != nil {
			return nil, fmt.Errorf("can't create post: %w", err)
		}
		return apireqs.ResponseCreatePost{ID: resp.Post.Id}, nil
	}
}

func GetPostByID(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		id := chi.URLParam(r, "id")
		input := apireqs.RequestGetPostByID{}
		err = network.ReadBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}

		slog.InfoContext(ctx, "Get Post", "id", id)

		requster := ""
		if input.Login != "" && input.Password != "" {
			userResp, err := a.Usclient.AuthUser(ctx, userservicerequests.RequestAuth{
				Login:    input.Login,
				Password: input.Password,
			})
			if err != nil {
				return nil, fmt.Errorf("can't auth user: %w", err)
			}
			requster = userResp.ID
		}

		resp, err := a.PostClient.GetPost(ctx, &post.GetPostRequest{
			Id:          id,
			RequesterId: requster,
		})
		if err != nil {
			return nil, fmt.Errorf("can't get post: %w", err)
		}

		_, _, err = a.Producer.SendMessage(&sarama.ProducerMessage{
			Topic: "posts",
			Value: sarama.StringEncoder(kevents.Event{Event: "View", UserID: requster, PostID: id}.String()),
		})
		if err != nil {
			return nil, fmt.Errorf("can't send message: %w", err)
		}

		return apireqs.ResponseGetPostByID{Post: convert(resp.Post)}, nil
	}
}

func GetPosts(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		page, err := strconv.Atoi(chi.URLParam(r, "page"))
		if err != nil {
			return nil, fmt.Errorf("can't get page: %w", err)
		}
		input := apireqs.RequestGetPosts{}
		err = network.ReadBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}
		slog.InfoContext(ctx, "Get Posts", "page", page)

		requster := ""
		if input.Login != "" && input.Password != "" {
			userResp, err := a.Usclient.AuthUser(ctx, userservicerequests.RequestAuth{
				Login:    input.Login,
				Password: input.Password,
			})
			if err != nil {
				return nil, fmt.Errorf("can't auth user: %w", err)
			}
			requster = userResp.ID
		}
		resp, err := a.PostClient.ListPosts(ctx, &post.ListPostsRequest{
			Page:        int32(page),
			PageSize:    10,
			RequesterId: requster,
		})
		if err != nil {
			return nil, fmt.Errorf("can't list posts: %w", err)
		}
		res := make([]models.Post, len(resp.Posts))
		for i := range res {
			res[i] = convert(resp.Posts[i])
		}

		return apireqs.ResponseGetPosts{Posts: res}, nil
	}
}

func UpdatePost(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		id := chi.URLParam(r, "id")
		slog.InfoContext(ctx, "Update Post", "id", id)
		input := apireqs.RequestUpdatePost{}
		err = network.ReadBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}

		userResp, err := a.Usclient.AuthUser(ctx, userservicerequests.RequestAuth{
			Login:    input.Login,
			Password: input.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("can't auth user: %w", err)
		}

		resp, err := a.PostClient.UpdatePost(ctx, &post.UpdatePostRequest{
			Id:          id,
			Title:       input.Title,
			Description: input.Description,
			IsPrivate:   input.IsPrivate,
			Tags:        input.Tags,
			RequesterId: userResp.ID,
		})

		if err != nil {
			return nil, fmt.Errorf("can't create post: %w", err)
		}
		return apireqs.ResponseUpdatePost{Post: convert(resp.Post)}, nil
	}
}

func DeletePost(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		id := chi.URLParam(r, "id")
		slog.InfoContext(ctx, "Update Post", "id", id)
		input := apireqs.RequestDeletePost{}
		err = network.ReadBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}

		userResp, err := a.Usclient.AuthUser(ctx, userservicerequests.RequestAuth{
			Login:    input.Login,
			Password: input.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("can't auth user: %w", err)
		}
		resp, err := a.PostClient.DeletePost(ctx, &post.DeletePostRequest{
			Id:          id,
			RequesterId: userResp.ID,
		})

		if err != nil {
			return nil, fmt.Errorf("can't create post: %w", err)
		}
		return apireqs.ResponseDeletePost{Success: resp.Success}, nil
	}
}

func LikePost(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		postId := chi.URLParam(r, "post_id")
		input := apireqs.RequestLikePost{}
		err = network.ReadBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}

		userResp, err := a.Usclient.AuthUser(ctx, userservicerequests.RequestAuth{
			Login:    input.Login,
			Password: input.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("can't auth user: %w", err)
		}
		requster := userResp.ID

		_, _, err = a.Producer.SendMessage(&sarama.ProducerMessage{
			Topic: "posts",
			Value: sarama.StringEncoder(kevents.Event{Event: "Like", UserID: requster, PostID: postId}.String()),
		})
		if err != nil {
			return nil, fmt.Errorf("can't send message: %w", err)
		}

		return apireqs.ResponseLikePost{}, nil
	}
}

func NewComment(a *api.Api) MyHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error) {
		postId := chi.URLParam(r, "post_id")
		input := apireqs.RequestNewComment{}
		err = network.ReadBody(r, &input)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}

		userResp, err := a.Usclient.AuthUser(ctx, userservicerequests.RequestAuth{
			Login:    input.Login,
			Password: input.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("can't auth user: %w", err)
		}
		requster := userResp.ID
		_, _, err = a.Producer.SendMessage(&sarama.ProducerMessage{
			Topic: "posts",
			Value: sarama.StringEncoder(kevents.Event{Event: "Like", UserID: requster, PostID: postId}.String()),
		})
		if err != nil {
			return nil, fmt.Errorf("can't send message: %w", err)
		}

		return apireqs.ResponseNewComment{}, nil
	}
}

func convert(p *post.Post) models.Post {
	return models.Post{
		ID:          models.PostID(p.GetId()),
		Title:       p.GetTitle(),
		Description: p.GetDescription(),
		CreatorID:   p.GetCreatorId(),
		CreatedAt:   p.GetCreatedAt().AsTime(),
		UpdatedAt:   p.GetUpdatedAt().AsTime(),
		IsPrivate:   p.GetIsPrivate(),
		Tags:        p.GetTags(),
	}
}
