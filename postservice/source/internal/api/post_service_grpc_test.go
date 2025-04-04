package api

import (
	"context"
	"postservice/internal/storage"
	"social/shared/models"
	"testing"

	pb "postservice/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var _ storage.CommonRepository = &storeMock{}

type storeMock struct {
	mock.Mock
}

func (s *storeMock) CreatePost(ctx context.Context, post models.Post) (models.PostID, error) {
	args := s.Called(post)
	return args.Get(0).(models.PostID), args.Error(1)
}

// DeletePost implements storage.CommonRepository.
func (s *storeMock) DeletePost(ctx context.Context, id models.PostID) error {
	args := s.Called(id)
	return args.Error(0)
}

// GetPostByID implements storage.CommonRepository.
func (s *storeMock) GetPostByID(ctx context.Context, id models.PostID) (models.Post, error) {
	args := s.Called(id)
	return args.Get(0).(models.Post), args.Error(1)
}

// GetPosts implements storage.CommonRepository.
func (s *storeMock) GetPosts(ctx context.Context, offset int, limit int) ([]models.Post, error) {
	args := s.Called(offset, limit)
	return args.Get(0).([]models.Post), args.Error(1)
}

// UpdatePost implements storage.CommonRepository.
func (s *storeMock) UpdatePost(ctx context.Context, post models.Post) error {
	args := s.Called(post)
	return args.Error(0)
}

func TestCreatePost(t *testing.T) {
	s := &storeMock{}
	ps := PostServiceServer{
		store: s,
	}
	req := pb.CreatePostRequest{
		Title:       "Title",
		Description: "Desc",
		CreatorId:   "sprff",
		IsPrivate:   false,
		Tags:        []string{"hype", "soa"},
	}
	s.On("CreatePost", mock.AnythingOfType("models.Post")).Return(models.PostID("id"), nil)
	res, err := ps.CreatePost(context.TODO(), &req)
	assert.NoError(t, err)
	t.Logf("res=%v", res)

	t.Run("invalid", func(t *testing.T) {
		req.Title = ""
		_, err := ps.CreatePost(context.TODO(), &req)
		assert.Error(t, err)
		req.Title = "Title"

		req.Description = ""
		_, err = ps.CreatePost(context.TODO(), &req)
		assert.NoError(t, err)
		req.Description = "Desc"

		req.CreatorId = ""
		_, err = ps.CreatePost(context.TODO(), &req)
		assert.Error(t, err)
		req.CreatorId = "sprff"
	})
}

func TestGet(t *testing.T) {
	s := &storeMock{}
	ps := PostServiceServer{
		store: s,
	}

	post := models.Post{
		ID:        "qwe",
		CreatorID: "sprff",
		Title:     "Some",
	}
	postPrivate := models.Post{
		ID:        "qwe_",
		CreatorID: "sprff",
		Title:     "Some",
		IsPrivate: true,
	}
	s.On("GetPostByID", models.PostID("qwe")).Return(post, nil)
	s.On("GetPostByID", models.PostID("qwe_")).Return(postPrivate, nil)

	res, err := ps.GetPost(context.TODO(), &pb.GetPostRequest{Id: "qwe", RequesterId: "sprff"})
	assert.NoError(t, err)
	assert.Equal(t, convertToPbPost(post), res.Post)
	res, err = ps.GetPost(context.TODO(), &pb.GetPostRequest{Id: "qwe_", RequesterId: "sprff"})
	assert.NoError(t, err)
	assert.Equal(t, convertToPbPost(postPrivate), res.Post)
	res, err = ps.GetPost(context.TODO(), &pb.GetPostRequest{Id: "qwe", RequesterId: "sprff2"})
	assert.NoError(t, err)
	assert.Equal(t, convertToPbPost(post), res.Post)
	_, err = ps.GetPost(context.TODO(), &pb.GetPostRequest{Id: "qwe_", RequesterId: "sprff2"})
	assert.Error(t, err)
}

func TestUpdate(t *testing.T) {
	s := &storeMock{}
	ps := PostServiceServer{
		store: s,
	}

	post := models.Post{
		ID:        "qwe",
		CreatorID: "sprff",
		Title:     "Some",
	}
	postPrivate := models.Post{
		ID:        "qwe_",
		CreatorID: "sprff",
		Title:     "Some",
		IsPrivate: true,
	}
	newPost := models.Post{
		ID:        "qwe",
		CreatorID: "sprff",
		Title:     "Some2",
	}
	newPostPrivate := models.Post{
		ID:        "qwe_",
		CreatorID: "sprff",
		Title:     "Some2",
	}
	s.On("GetPostByID", models.PostID("qwe")).Return(post, nil)
	s.On("GetPostByID", models.PostID("qwe_")).Return(postPrivate, nil)
	s.On("UpdatePost", newPost).Return(nil)
	s.On("UpdatePost", newPostPrivate).Return(nil)

	res, err := ps.UpdatePost(context.TODO(), &pb.UpdatePostRequest{Id: "qwe", Title: "Some2", RequesterId: "sprff"})
	assert.NoError(t, err)
	assert.Equal(t, convertToPbPost(newPost), res.Post)

	res, err = ps.UpdatePost(context.TODO(), &pb.UpdatePostRequest{Id: "qwe_", Title: "Some2", RequesterId: "sprff"})
	assert.NoError(t, err)
	assert.Equal(t, convertToPbPost(newPostPrivate), res.Post)

	_, err = ps.UpdatePost(context.TODO(), &pb.UpdatePostRequest{Id: "qwe", Title: "Some2", RequesterId: "sprff2"})
	assert.Error(t, err)

	_, err = ps.UpdatePost(context.TODO(), &pb.UpdatePostRequest{Id: "qwe_", Title: "Some2", RequesterId: "sprff2"})
	assert.Error(t, err)
}
