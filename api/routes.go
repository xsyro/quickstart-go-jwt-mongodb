package api

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"quickstart-go-jwt-mongodb/internal"
	"time"
)

type (
	HttpVerb     string
	Request      = http.Request
	WithResource struct {
		HttpRequest   *HttpRequestHandler
		MongoDatabase internal.MongoDatabase
	}
	// HandleRequest Declare how HTTP route needs to be handled
	HandleRequest struct {
		Uri      string
		Method   HttpVerb
		Callback func(w http.ResponseWriter, req *http.Request)
	}
	HttpErrorResp struct {
		IsError bool   `json:"is_error"`
		Message string `json:"message"`
	}
	HttpRequestHandler struct {
		Router *mux.Router
		port   string
	}
)

const (
	GET    HttpVerb = "GET"
	POST   HttpVerb = "POST"
	PUT    HttpVerb = "PUT"
	DELETE HttpVerb = "DELETE"
)

func NewHttpRequestHandler() *HttpRequestHandler {
	return &HttpRequestHandler{
		Router: mux.NewRouter().PathPrefix(os.Getenv("BASE_URI_PREFIX")).Subrouter(),
		port:   os.Getenv("HTTP_PORT"),
	}
}

func (appRouter *HttpRequestHandler) HandleMiddlewares(middlewares ...mux.MiddlewareFunc) {
	appRouter.Router.Use(middlewares...)
}

func (appRouter *HttpRequestHandler) HandleRequest(handler HandleRequest) {
	appRouter.Router.HandleFunc(handler.Uri, handler.Callback).Methods(string(handler.Method), http.MethodOptions)
}

func (appRouter *HttpRequestHandler) Serve() {
	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%s", "0.0.0.0", appRouter.port),
		ReadTimeout: 60 * time.Second,
		Handler:     appRouter.Router,
	}
	log.Infof("serving HTTP Request on port %s", appRouter.port)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
		return
	}
}
