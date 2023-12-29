package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type (
	Verb    string
	Request = http.Request
	// Controller Declare how HTTP needs to be handled
	Controller struct {
		Uri         string
		Method      Verb
		Secure      bool
		PermitRoles []string
		Callback    func(w http.ResponseWriter, req *http.Request)
	}
	ResponseBody struct {
		IsError bool        `json:"is_error"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
	Handler struct {
		router          *mux.Router
		port            string
		httpTimeout     context.Context
		requestRegistry []Controller
	}
	RequestHandler interface {
		HandleMiddlewares(middlewares ...mux.MiddlewareFunc)
		ControllerRegistry(handler ...Controller)
		GetControllers() []Controller
		Serve()
	}
)

func (appRouter *Handler) GetControllers() []Controller {
	return appRouter.requestRegistry
}

const (
	GET    Verb = "GET"
	POST   Verb = "POST"
	PUT    Verb = "PUT"
	DELETE Verb = "DELETE"
)

var BaseUrlPrefix = os.Getenv("BASE_URI_PREFIX")

func NewHttpRequestHandler(httpTimeoutCtx context.Context) RequestHandler {
	return &Handler{
		router:          mux.NewRouter().PathPrefix(BaseUrlPrefix).Subrouter(),
		port:            os.Getenv("HTTP_PORT"),
		httpTimeout:     httpTimeoutCtx,
		requestRegistry: []Controller{},
	}
}

func (appRouter *Handler) HandleMiddlewares(middlewares ...mux.MiddlewareFunc) {
	appRouter.router.Use(middlewares...)
}

func (appRouter *Handler) ControllerRegistry(controllers ...Controller) {
	for i := range controllers {
		appRouter.requestRegistry = append(appRouter.requestRegistry, controllers[i])
	}
}

func (appRouter *Handler) Serve() {

	for _, request := range appRouter.requestRegistry {
		appRouter.router.HandleFunc(request.Uri, request.Callback).Methods(string(request.Method), http.MethodOptions)
	}

	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%s", "0.0.0.0", appRouter.port),
		ReadTimeout: 60 * time.Second,
		Handler:     appRouter.router,
	}
	log.Infof("serving HTTP Request on port %s", appRouter.port)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
		return
	}
}

func ParseReqToJson(req *Request, obj interface{}) error {
	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(&obj)
	if err != nil {
		return err
	}
	return validator.New().Struct(obj)
}

func HttpResponse(w http.ResponseWriter, statusCode int, obj interface{}) {
	err := json.NewEncoder(w).Encode(ResponseBody{
		IsError: false,
		Message: "Request Completed",
		Data:    obj,
	})
	if err != nil {
		log.Error("error sending server response")
	}
	w.WriteHeader(statusCode)
}

func HttpError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	if json.NewEncoder(w).Encode(ResponseBody{
		IsError: true,
		Message: fmt.Sprintf("%s", err),
	}) != nil {
		log.Error("error sending server response")
	}
	w.WriteHeader(http.StatusBadRequest)
}

func AccessDenied(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	if json.NewEncoder(w).Encode(ResponseBody{
		IsError: true,
		Message: fmt.Sprintf("%s", err),
	}) != nil {
		log.Error("error sending server response")
	}
}
