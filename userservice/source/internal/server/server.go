package server

import (
	"context"
	"log/slog"
	"net/http"
	"social/shared/models"
	"social/userservice/internal/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func GetRouter(a *api.Api) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("/ touched")
		w.Write([]byte("userserver\n"))
	})

	router.Post("/users", prepareHandler(RegisterUser(a)))
	router.Post("/users/auth", prepareHandler(AuthUser(a)))
	router.Get("/users/{id}", prepareHandler(GetUserByID(a)))
	router.Put("/users/{id}", prepareHandler(UpdateUser(a)))
	router.Get("/users/bylogin/{login}", prepareHandler(GetUserByLogin(a)))

	return router
}

type MyHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error)

func prepareHandler(handler MyHandlerFunc) http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		sd := models.NewStatusData(handler(ctx, w, r))
		render.JSON(w, r, sd)
		slog.DebugContext(ctx, "Responsed", "status", sd.Status)
	}
}
