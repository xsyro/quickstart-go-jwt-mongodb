package controllers

import (
	"context"
	"net/http"
	"quickstart-go-jwt-mongodb/server"
)

func Homepage(ctx context.Context) server.Controller {
	return server.Controller{
		Uri:    "/",
		Secure: false,
		Method: server.GET,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Hello!"))
			return
		},
	}
}

func HealthCheck(ctx context.Context) server.Controller {
	return server.Controller{
		Uri:    "/status",
		Secure: false,
		Method: server.GET,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Good!"))
			return
		},
	}
}
