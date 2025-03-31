package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"social/apiservice/internal/api"
	pb "social/apiservice/internal/proto"
	"social/apiservice/internal/server"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	slog.SetLogLoggerLevel(slog.LevelDebug)
	conn, err := grpc.Dial("postservice-backend:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPostServiceClient(conn)
	a := api.NewApi(
		api.UserserviceClient{Url: "http://userservice-backend:8080"},
		client,
	)
	router := server.GetRouter(a)
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", 8080),
		Handler:      router,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	slog.Info("Started")
	err = srv.ListenAndServe()
	processError("Failed to start server", err)
	slog.Info("Server stopped")
}

func processError(msg string, err error) {
	if err != nil {
		slog.Error(msg, slog.String("error", err.Error()))
		os.Exit(1)
	}
}
