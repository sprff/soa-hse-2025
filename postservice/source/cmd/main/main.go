package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"postservice/internal/api"
	pb "postservice/proto"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterPostServiceServer(server, api.NewPostServiceServer())

	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
