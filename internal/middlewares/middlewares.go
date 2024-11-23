package middlewares

import (
	"net/http"

	"github.com/Bitummit/booking_api/internal/api"
	authclient "github.com/Bitummit/booking_api/internal/service/authClient"
	"github.com/go-chi/render"
	"github.com/Bitummit/booking_api/pkg/config"
)

func SetJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func IsAdmin(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, api.ErrorResponse("unauthorized"))
				return
			}

			authClient, err := authclient.New(cfg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, api.ErrorResponse("internal grpc error"))
				return
			}
			if err = authClient.CheckIsADmin(token); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, api.ErrorResponse("no enough permission"))
				return
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}