package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"quickstart-go-jwt-mongodb/internal"
	"quickstart-go-jwt-mongodb/services"
	"quickstart-go-jwt-mongodb/types"
	"slices"
	"strings"
	"time"
)

type (
	HttpVerb     string
	Request      = http.Request
	WithResource struct {
		HttpRequest   RequestHandler
		MongoDatabase internal.MongoDatabase
	}
	// HttpRequest Declare how HTTP route needs to be handled
	HttpRequest struct {
		Uri         string
		Method      HttpVerb
		Secure      bool
		PermitRoles []string
		Callback    func(w http.ResponseWriter, req *http.Request)
	}
	HttpResponseBody struct {
		IsError bool        `json:"is_error"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
	HttpRequestHandler struct {
		router          *mux.Router
		port            string
		httpTimeout     context.Context
		requestRegistry []HttpRequest
	}
	RequestHandler interface {
		HandleMiddlewares(middlewares ...mux.MiddlewareFunc)
		RequestRegistry(handler HttpRequest)
		SecureMiddleware() mux.MiddlewareFunc
		Serve()
	}
)

const (
	GET    HttpVerb = "GET"
	POST   HttpVerb = "POST"
	PUT    HttpVerb = "PUT"
	DELETE HttpVerb = "DELETE"
)

const (
	HeaderName   = "Authorization"
	HeaderScheme = "Bearer"
)

var pathPrefix = os.Getenv("BASE_URI_PREFIX")

func NewHttpRequestHandler(httpTimeoutCtx context.Context) RequestHandler {
	return &HttpRequestHandler{
		router:          mux.NewRouter().PathPrefix(pathPrefix).Subrouter(),
		port:            os.Getenv("HTTP_PORT"),
		httpTimeout:     httpTimeoutCtx,
		requestRegistry: []HttpRequest{},
	}
}

func (appRouter *HttpRequestHandler) HandleMiddlewares(middlewares ...mux.MiddlewareFunc) {
	appRouter.router.Use(middlewares...)
}

func (appRouter *HttpRequestHandler) RequestRegistry(handler HttpRequest) {
	appRouter.requestRegistry = append(appRouter.requestRegistry, handler)
}

func (appRouter HttpRequestHandler) SecureMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			if req.Method == http.MethodOptions {
				next.ServeHTTP(w, req)
			}
			//Check if the RequestRegistry is not Secured. Otherwise, continue with security validation checks
			var currentHttpRequest HttpRequest
			for i := range appRouter.requestRegistry {
				if fmt.Sprintf("%s%s", pathPrefix, appRouter.requestRegistry[i].Uri) == req.URL.Path {
					currentHttpRequest = appRouter.requestRegistry[i]
					break
				}
			}

			//non-secured page should be served their corresponding http handler
			if !currentHttpRequest.Secure {
				next.ServeHTTP(w, req)
				return
			}

			//Access-Control-Level and JWT Expiring Validations
			if req.Header[HeaderName] == nil {
				accessDenied(w, errors.New(fmt.Sprintf("'%s' not found in the HTTP Request Header", HeaderName)))
				return
			}

			var user types.User
			jwtTokenizedStr := strings.Trim(strings.TrimLeft(req.Header[HeaderName][0], HeaderScheme), " ")
			jwt := services.NewJwtService(appRouter.httpTimeout)
			claims, err := jwt.ClaimToken(jwtTokenizedStr, &user)
			if err != nil {
				accessDenied(w, errors.New(fmt.Sprintf("access denied. %v", err)))
				return
			}

			expirationTime, _ := claims.GetExpirationTime()

			if time.Now().After(expirationTime.Time) {
				w.Header().Add("Expires", "true")
				accessDenied(w, errors.New("unauthorized. Token expired"))
				return
			}

			//Empty 'PermitRoles' signifies wild card ACL. Authorization check isn't required. Just authentication
			if len(currentHttpRequest.PermitRoles) == 0 {
				next.ServeHTTP(w, req)
				return
			}
			for i := range user.Roles {
				if slices.Contains(currentHttpRequest.PermitRoles, user.Roles[i]) {
					next.ServeHTTP(w, req)
					return
				}
			}
			accessDenied(w, errors.New(fmt.Sprintf("unauthorised accces to this URL")))
		})
	}
}

func (appRouter *HttpRequestHandler) Serve() {

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
	w.WriteHeader(http.StatusUnauthorized)
	if json.NewEncoder(w).Encode(HttpResponseBody{
		IsError: true,
		Message: fmt.Sprintf("%s", err),
	}) != nil {
		log.Error("error sending http response")
	}
	w.WriteHeader(http.StatusBadRequest)
}

func accessDenied(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	if json.NewEncoder(w).Encode(HttpResponseBody{
		IsError: true,
		Message: fmt.Sprintf("%s", err),
	}) != nil {
		log.Error("error sending http response")
	}
}
