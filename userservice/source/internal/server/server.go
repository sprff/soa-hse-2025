package server

import (
	"context"
	"log/slog"
	"net/http"
	"userservice/internal/api"

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
	router.Get("/users/{id}", prepareHandler(GetUserByID(a)))
	router.Put("/users/{id}", prepareHandler(UpdateUser(a)))
	router.Get("/users/bylogin/{login}", prepareHandler(GetUserByLogin(a)))

	return router
}

type MyHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) (response any, err error)

type response struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

func prepareHandler(handler MyHandlerFunc) http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := handler(ctx, w, r)
		if err != nil {
			response := writeError(err)
			render.JSON(w, r, response)
			return
		}

		if resp != nil {
			render.JSON(w, r, response{Status: "OK", Data: resp})
		}
	}
}

func writeError(err error) any {
	return response{Status: err.Error(), Data: err}
}

