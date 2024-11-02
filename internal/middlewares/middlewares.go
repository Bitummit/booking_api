package middlewares

import "net/http"

func SetJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func IsAdmin(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// grpc auth microservice call to check admin
		next.ServeHTTP(w, r)
	})
}