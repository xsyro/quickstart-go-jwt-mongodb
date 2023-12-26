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
	resources.HttpRequest.HandleRequest(createAccount(resources.MongoDatabase))
	resources.HttpRequest.HandleRequest(authenticate(resources.MongoDatabase))
}

func createAccount(mongoDb internal.MongoDatabase) HandleRequest {
	return HandleRequest{
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

			userRepository := repositories.NewUserRepository(mongoDb)
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

func authenticate(mongoDb internal.MongoDatabase) HandleRequest {
	return HandleRequest{
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

			userRepository := repositories.NewUserRepository(mongoDb)
			tokenRepository := repositories.NewTokenRepository(mongoDb)
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
			tokenizedStr, err := jwtService.GenerateJWT(user)
			token := types.Token{
				BaseModel:   types.NewBaseModel(),
				AccessToken: tokenizedStr,
			}
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

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateHashPassword(plainText string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	return string(bytes), err
}
