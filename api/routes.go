package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
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
	HttpResponseBody struct {
		IsError bool        `json:"is_error"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
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

func parseReqToJson(req *Request, obj interface{}) error {
	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(&obj)
	if err != nil {
		return err
	}
	return validator.New().Struct(obj)
}

func httpResponse(w http.ResponseWriter, statusCode int, obj interface{}) {
	err := json.NewEncoder(w).Encode(HttpResponseBody{
		IsError: false,
		Message: "Request Completed",
		Data:    obj,
	})
	if err != nil {
		log.Error("error sending http response")
	}
	w.WriteHeader(statusCode)
}

func httpError(w http.ResponseWriter, err error) {
	if json.NewEncoder(w).Encode(HttpResponseBody{
		IsError: true,
		Message: fmt.Sprintf("%s", err),
	}) != nil {
		log.Error("error sending http response")
	}
	w.WriteHeader(http.StatusBadRequest)
}
