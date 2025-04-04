package main

import (
	"context"
	"fmt"
	"log"
	"postservice/internal/storage/psql"
	"social/shared/models"
)

func main() {
	s, err := psql.NewPsqlStorage()
	if err != nil {
		log.Fatalf("err: %s\n", err)
	}
	ctx := context.Background()
	id, err := s.CreatePost(ctx, models.Post{
		Title:       "First Post",
		Description: "This is my first post",
		CreatorID:   "user1",
		IsPrivate:   false,
		Tags:        []string{"golang", "grpc"},
	})
	fmt.Println(id, err)
	post, err := s.GetPostByID(ctx, id)
	fmt.Println(post, err)
	posts, err := s.GetPosts(ctx, 1, 3)
	fmt.Println(posts, err)

}
