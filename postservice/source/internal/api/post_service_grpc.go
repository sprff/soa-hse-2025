package api

import (
	"context"
	"fmt"
	"social/shared/models"
	"sync"
	"time"

	"postservice/internal/storage"
	"postservice/internal/storage/psql"
	pb "postservice/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostServiceServer struct {
	pb.UnimplementedPostServiceServer
	mu    sync.RWMutex
	store storage.CommonRepository
}

func NewPostServiceServer() *PostServiceServer {
	store, err := psql.NewPsqlStorage()
	if err != nil {
		panic(fmt.Sprintf("can't create psql: %s", err.Error()))
	}
	return &PostServiceServer{
		store: store,
	}
}

func (s *PostServiceServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	if req.CreatorId == "" {
		return nil, status.Error(codes.InvalidArgument, "creator_id is required")
	}

	newPost := models.Post{
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   req.CreatorId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsPrivate:   req.IsPrivate,
		Tags:        req.Tags,
	}

	s.mu.Lock()
	id, err := s.store.CreatePost(ctx, newPost)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	newPost.ID = id
	s.mu.Unlock()

	return &pb.CreatePostResponse{
		Post: convertToPbPost(newPost),
	}, nil
}

func (s *PostServiceServer) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	p, err := s.store.GetPostByID(ctx, models.PostID(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if p.IsPrivate && p.CreatorID != req.RequesterId {
		return nil, status.Error(codes.PermissionDenied, "you don't have permission to view this post")
	}

	return &pb.GetPostResponse{
		Post: convertToPbPost(p),
	}, nil
}

func (s *PostServiceServer) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.RequesterId == "" {
		return nil, status.Error(codes.InvalidArgument, "updater_id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	post, err := s.store.GetPostByID(ctx, models.PostID(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if post.CreatorID != req.RequesterId {
		return nil, status.Error(codes.PermissionDenied, "you can only update your own posts")
	}

	// Обновляем поля
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Description != "" {
		post.Description = req.Description
	}
	post.IsPrivate = req.IsPrivate
	post.Tags = req.Tags
	err = s.store.UpdatePost(ctx, post)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdatePostResponse{
		Post: convertToPbPost(post),
	}, nil
}

func (s *PostServiceServer) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.RequesterId == "" {
		return nil, status.Error(codes.InvalidArgument, "deleter_id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	post, err := s.store.GetPostByID(ctx, models.PostID(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if post.CreatorID != req.RequesterId {
		return nil, status.Error(codes.PermissionDenied, "you can only update your own posts")
	}

	err = s.store.DeletePost(ctx, models.PostID(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeletePostResponse{
		Success: true,
	}, nil
}

func (s *PostServiceServer) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var posts []models.Post
	offset := 0
	for {
		ps, err := s.store.GetPosts(ctx, offset, 200)
		if err != nil {
			if err == models.ErrPostNotFound {
				break
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
		offset += 200
		for _, p := range ps {
			if !p.IsPrivate || p.CreatorID == req.RequesterId {
				posts = append(posts, p)
			}
		}
	}
	total := int32(len(posts))
	start := (req.Page - 1) * req.PageSize
	if start >= total {
		return &pb.ListPostsResponse{
			Posts:      []*pb.Post{},
			TotalCount: int32(total),
			Page:       req.Page,
			PageSize:   req.PageSize,
		}, nil
	}

	end := start + req.PageSize
	if end > total {
		end = total
	}

	paginatedPosts := posts[start:end]

	var pbPosts []*pb.Post
	for _, p := range paginatedPosts {
		pbPosts = append(pbPosts, convertToPbPost(p))
	}

	return &pb.ListPostsResponse{
		Posts:      pbPosts,
		TotalCount: int32(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

func convertToPbPost(p models.Post) *pb.Post {
	return &pb.Post{
		Id:          string(p.ID),
		Title:       p.Title,
		Description: p.Description,
		CreatorId:   p.CreatorID,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
		IsPrivate:   p.IsPrivate,
		Tags:        p.Tags,
	}
}
