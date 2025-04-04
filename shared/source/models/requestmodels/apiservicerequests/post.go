package apiservicerequests

import (
	"social/shared/models"
)

type RequestCreatePost struct {
	Login       string   `json:"login"`
	Password    string   `json:"password"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	CreatorID   string   `json:"creator_id"`
	IsPrivate   bool     `json:"is_private"`
	Tags        []string `json:"tags"`
}
type ResponseCreatePost struct {
	ID string `json:"id"`
}

type RequestUpdatePost struct {
	Login       string   `json:"login"`
	Password    string   `json:"password"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	CreatorID   string   `json:"creator_id"`
	IsPrivate   bool     `json:"is_private"`
	Tags        []string `json:"tags"`
}
type ResponseUpdatePost struct {
	Post models.Post `json:"post"`
}

type RequestGetPostByID struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type ResponseGetPostByID struct {
	Post models.Post `json:"post"`
}

type RequestGetPosts struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type ResponseGetPosts struct {
	Posts []models.Post `json:"posts"`
}

type RequestDeletePost struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type ResponseDeletePost struct {
	Success bool `json:"success"`
}
