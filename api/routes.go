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
	"strings"
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
		Router      *mux.Router
		port        string
		httpTimeout context.Context
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

func NewHttpRequestHandler(httpTimeoutCtx context.Context) *HttpRequestHandler {
	return &HttpRequestHandler{
		Router:      mux.NewRouter().PathPrefix(os.Getenv("BASE_URI_PREFIX")).Subrouter(),
		port:        os.Getenv("HTTP_PORT"),
		httpTimeout: httpTimeoutCtx,
	}
}

func (appRouter *HttpRequestHandler) HandleMiddlewares(middlewares ...mux.MiddlewareFunc) {
	appRouter.Router.Use(middlewares...)
}

func (appRouter *HttpRequestHandler) HandleRequest(handler HandleRequest) {
	appRouter.authorizationMiddleware(handler)
	appRouter.Router.HandleFunc(handler.Uri, handler.Callback).Methods(string(handler.Method), http.MethodOptions)
}

func (appRouter *HttpRequestHandler) authorizationMiddleware(handler HandleRequest) {
	appRouter.HandleMiddlewares(appRouter.securedMiddleware(handler))
}

func (appRouter HttpRequestHandler) securedMiddleware(request HandleRequest) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			if req.Method == http.MethodOptions {
				next.ServeHTTP(w, req)
			}

			//Check if the HandleRequest is not Secured. Otherwise, continue with security validation checks
			if !request.Secure {
				next.ServeHTTP(w, req)
				return
			}

			if req.Header[HeaderName] == nil {
				accessDenied(w, errors.New(fmt.Sprintf("'%s' not found int the HTTP Request Header", HeaderName)))
				return
			}
			jwtTokenizedStr := strings.TrimLeft(strings.Trim(req.Header[HeaderName][0], " "), HeaderScheme)
			jwt := services.NewJwtService(appRouter.httpTimeout)
			claims, err := jwt.ParseJWT(jwtTokenizedStr)
			if err != nil {
				w.Header().Add("Expires", "true")
				accessDenied(w, errors.New("access denied"))
				return
			}

			var user types.User
			jsonData, _ := json.Marshal(claims["obj"])
			err = json.NewDecoder(strings.NewReader(string(jsonData))).Decode(&user)

			if err != nil {
				log.Error("Error", err)
			}
			//var allowUrls []string

			//allowUrls := AclUrl[user.Roles]

			//switch user.Roles {
			//case "CASHIER", "ACCOUNTANT", "OPERATION":
			//	for i := range allowUrls {
			//		if strings.HasPrefix(req.URL.Path, allowUrls[i]) {
			//			next.ServeHTTP(w, req)
			//			return
			//		}
			//	}
			//	//var err = routes.HttpErrorResp{}
			//	//err = routes.SetHttpErrorResp(err, "Not Authorized")
			//	//routes.WriteHttpResponse(w, http.StatusUnauthorized, &err)
			//
			//	next.ServeHTTP(w, req)
			//	return

			//default:
			//}
			//

			next.ServeHTTP(w, req)

		})
	}
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

func accessDenied(w http.ResponseWriter, err error) {
	if json.NewEncoder(w).Encode(HttpResponseBody{
		IsError: true,
		Message: fmt.Sprintf("%s", err),
	}) != nil {
		log.Error("error sending http response")
	}
	w.WriteHeader(http.StatusUnauthorized)
}
