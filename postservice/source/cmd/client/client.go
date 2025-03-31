package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "postservice/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:5002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPostServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Тестирование создания поста
	createRes, err := client.CreatePost(ctx, &pb.CreatePostRequest{
		Title:       "First Post",
		Description: "This is my first post",
		CreatorId:   "user1",
		IsPrivate:   false,
		Tags:        []string{"golang", "grpc"},
	})
	if err != nil {
		log.Fatalf("CreatePost failed: %v", err)
	}
	log.Printf("Created post: %v", createRes.Post)

	postID := createRes.Post.Id

	// Тестирование получения поста
	getRes, err := client.GetPost(ctx, &pb.GetPostRequest{
		Id:          postID,
		RequesterId: "user1",
	})
	if err != nil {
		log.Fatalf("GetPost failed: %v", err)
	}
	log.Printf("Got post: %v", getRes.Post)

	// Тестирование обновления поста
	updateRes, err := client.UpdatePost(ctx, &pb.UpdatePostRequest{
		Id:          postID,
		Title:       "Updated First Post",
		Description: "This is my updated first post",
		IsPrivate:   true,
		Tags:        []string{"golang", "grpc", "updated"},
		RequesterId: "user1",
	})
	if err != nil {
		log.Fatalf("UpdatePost failed: %v", err)
	}
	log.Printf("Updated post: %v", updateRes.Post)

	// Тестирование списка постов
	listRes, err := client.ListPosts(ctx, &pb.ListPostsRequest{
		Page:        1,
		PageSize:    10,
		RequesterId: "user1",
	})
	if err != nil {
		log.Fatalf("ListPosts failed: %v", err)
	}
	log.Printf("List posts: %v (total: %d)", listRes.Posts, listRes.TotalCount)

	// Тестирование удаления поста
	delRes, err := client.DeletePost(ctx, &pb.DeletePostRequest{
		Id:          postID,
		RequesterId: "user1",
	})
	if err != nil {
		log.Fatalf("DeletePost failed: %v", err)
	}
	log.Printf("Delete post success: %v", delRes.Success)
}
