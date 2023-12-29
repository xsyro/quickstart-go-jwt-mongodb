package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"quickstart-go-jwt-mongodb/models"
	"quickstart-go-jwt-mongodb/server"
	"quickstart-go-jwt-mongodb/services"
	"slices"
	"strings"
	"time"
)

const (
	HeaderName   = "Authorization"
	HeaderScheme = "Bearer"
)

func SecureMiddleware(controllers []server.Controller, envVar models.EnvVar) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			if req.Method == http.MethodOptions {
				next.ServeHTTP(w, req)
			}
			//Check if the ControllerRegistry is not Secured. Otherwise, continue with security validation checks
			var currentHttpRequest server.Controller
			for i := range controllers {
				if fmt.Sprintf("%s%s", envVar.BaseUrlPrefix, controllers[i].Uri) == req.URL.Path {
					currentHttpRequest = controllers[i]
					break
				}
			}

			//non-secured page should be served their corresponding server handler
			if !currentHttpRequest.Secure {
				next.ServeHTTP(w, req)
				return
			}

			//Access-Control-Level and JWT Expiring Validations
			if req.Header[HeaderName] == nil {
				server.AccessDenied(w, errors.New(fmt.Sprintf("'%s' not found in the HTTP Request Header", HeaderName)))
				return
			}

			var user models.User
			jwtTokenizedStr := strings.Trim(strings.TrimLeft(req.Header[HeaderName][0], HeaderScheme), " ")
			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
			defer cancel()
			jwt := services.NewJwtService(ctx, envVar)
			claims, err := jwt.ClaimToken(jwtTokenizedStr, &user)
			if err != nil {
				server.AccessDenied(w, errors.New(fmt.Sprintf("access denied. %v", err)))
				return
			}

			expirationTime, _ := claims.GetExpirationTime()

			if time.Now().After(expirationTime.Time) {
				w.Header().Add("Expires", "true")
				server.AccessDenied(w, errors.New("unauthorized. Token expired"))
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
			server.AccessDenied(w, errors.New(fmt.Sprintf("unauthorised accces to this URL")))
		})
	}
}
