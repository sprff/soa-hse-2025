package models

import (
	"errors"
	"time"
)

type PostID string

type Post struct {
	ID          PostID    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatorID   string    `json:"creator_id" db:"creator_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsPrivate   bool      `json:"is_private" db:"is_private"`
	Tags        []string  `json:"tags" db:"tags"`
}

var ErrPostNotFound = errors.New("user not found")
