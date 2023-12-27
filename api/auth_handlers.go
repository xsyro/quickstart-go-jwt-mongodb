package api

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"quickstart-go-jwt-mongodb/internal"
	"quickstart-go-jwt-mongodb/repositories"
	"quickstart-go-jwt-mongodb/services"
	"quickstart-go-jwt-mongodb/types"
	"time"
)

func AuthHandlers(resources *WithResource) {
	resources.HttpRequest.RequestRegistry(createAccount(resources.MongoDatabase))
	resources.HttpRequest.RequestRegistry(authenticate(resources.MongoDatabase))
	resources.HttpRequest.RequestRegistry(refreshToken(resources.MongoDatabase))
	resources.HttpRequest.RequestRegistry(listCustomers(resources.MongoDatabase))

}

func createAccount(database internal.MongoDatabase) HttpRequest {
	return HttpRequest{
		Uri:    "/account/create",
		Method: POST,
		Secure: false,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			var user = types.User{
				BaseModel: types.NewBaseModel(),
			}
			err := parseReqToJson(req, &user)
			if err != nil {
				httpError(w, err)
				return
			}

			userRepository := repositories.NewUserRepository(database)
			if userRepository.FindOne(ctx, &user, repositories.Filter{Key: "email", Value: user.Email}) {
				httpError(w, errors.New(fmt.Sprintf("%s already exists", user.Email)))
				return
			}
			password, err := generateHashPassword(user.PasswordRequestBody)
			if err != nil {
				httpError(w, errors.New("unable to hash password. Please try again later"))
				return
			}
			user.Password = password
			objectId, err := userRepository.CreateOne(ctx, &user)
			if err != nil {
				httpError(w, err)
				return
			}
			user.ID = objectId
			user.PasswordRequestBody = ""
			httpResponse(w, http.StatusCreated, user)
			return
		},
	}
}

type auth struct {
	Username string `json:"username" validate:"required,email"`
	Password string `json:"password"`
}

func authenticate(database internal.MongoDatabase) HttpRequest {
	return HttpRequest{
		Uri:    "/account/auth",
		Method: POST,
		Secure: false,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			var user types.User
			var auth auth
			err := parseReqToJson(req, &auth)
			if err != nil {
				httpError(w, err)
				return
			}

			userRepository := repositories.NewUserRepository(database)
			tokenRepository := repositories.NewTokenRepository(database)
			findOne := userRepository.FindOne(ctx, &user, repositories.Filter{Key: "email", Value: auth.Username})
			if !findOne {
				httpError(w, errors.New(fmt.Sprintf("Invalid Username. %s does not exists", auth.Username)))
				return
			}

			if !checkPasswordHash(auth.Password, user.Password) {
				httpError(w, errors.New("invalid Credential supplied. Please check username/password"))
				return
			}

			jwtService := services.NewJwtService(ctx)
			extraClaims := map[string]any{
				"iss": req.Host,
			}
			accessTokenStr, err := jwtService.GenerateJWT(user, 30*time.Minute, extraClaims)
			refreshTokenStr, err := jwtService.GenerateJWT(user.Email, 24*time.Hour, extraClaims)
			token := types.Token{
				BaseModel:    types.NewBaseModel(),
				AccessToken:  accessTokenStr,
				RefreshToken: refreshTokenStr,
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "jwt",
				Value:    refreshTokenStr,
				Expires:  time.Now().Add(25 * time.Hour),
				MaxAge:   60 * 60 * 24,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			})
			tokenId, err := tokenRepository.CreateOne(ctx, token)
			if err != nil {
				log.Error("Unable to persist token generated", err)
			}
			token.ID = tokenId

			if err != nil {
				httpError(w, errors.New("unable to generated token. Please try again later"))
				return
			}

			httpResponse(w, http.StatusCreated, token)
			return
		},
	}
}

func refreshToken(database internal.MongoDatabase) HttpRequest {
	return HttpRequest{
		Uri:    "/account/refresh-token",
		Method: GET,
		Secure: false,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			cookie, err := req.Cookie("jwt")
			var email string
			if err != nil {
				accessDenied(w, errors.New("'jwt' cookie do not exist in cookie-header"))
				return
			}
			jwtService := services.NewJwtService(ctx)
			jwt, err := jwtService.ClaimToken(cookie.Value, &email)
			if err != nil {
				accessDenied(w, errors.New("invalid jwt token supplied"))
				return
			}
			expirationTime, _ := jwt.GetExpirationTime()
			if err != nil {
				return
			}
			if time.Now().After(expirationTime.Time) {
				accessDenied(w, errors.New("unauthorized. Token expired"))
				return
			}

			userRepository := repositories.NewUserRepository(database)
			var user types.User
			userRepository.FindOne(ctx, &user, repositories.Filter{Key: "email", Value: email})

			//Create a one-time only access token again
			accessTokenStr, err := jwtService.GenerateJWT(user, 10*time.Hour)
			token := types.Token{
				BaseModel:   types.NewBaseModel(),
				AccessToken: accessTokenStr,
			}

			tokenRepository := repositories.NewTokenRepository(database)
			id, err := tokenRepository.CreateOne(ctx, token)
			if err != nil {
				return
			}
			token.ID = id
			httpResponse(w, http.StatusCreated, token)
			return
		},
	}
}

// listCustomers - Empty '[]PermitRoles' is wildcard access to all users.
// By simply excluding the PermitRole filed from the HttpRequest struct, it permits all secured users to
// access the page
func listCustomers(database internal.MongoDatabase) HttpRequest {
	return HttpRequest{
		Uri:    "/user/customer-records",
		Method: GET,
		Secure: true,
		Callback: func(responseWriter http.ResponseWriter, req *http.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			userRepository := repositories.NewUserRepository(database)
			var users []types.User
			err := userRepository.FindAll(ctx, &users)
			if err != nil {
				httpError(responseWriter, err)
				return
			}
			httpResponse(responseWriter, http.StatusOK, users)
		},
	}
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateHashPassword(plainText string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	return string(bytes), err
}
