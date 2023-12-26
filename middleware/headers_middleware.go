package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func HeadersMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization, Accept-Language")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTION")
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Access-Control-Max-Age", "86000") //browser cache cors preflight request for 14secs

			next.ServeHTTP(w, req)
		})
	}
}
