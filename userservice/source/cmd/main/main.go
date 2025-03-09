package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"social/userservice/internal/api"
	"social/userservice/internal/server"
	"social/userservice/internal/storage/psql"
	"time"
)

func main() {

	store, err := psql.NewPsqlStorage()
	processError("Failed to create psql store", err)
	slog.SetLogLoggerLevel(slog.LevelDebug)
	a := api.NewApi(store)
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
