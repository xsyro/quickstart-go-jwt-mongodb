package api

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"quickstart-go-jwt-mongodb/internal"
	"quickstart-go-jwt-mongodb/repositories"
	"quickstart-go-jwt-mongodb/types"
	"time"
)

func AuthHandlers(resources *WithResource) {
	resources.HttpRequest.HandleRequest(createAccount(resources.MongoDatabase))
}

func createAccount(mongoDb internal.MongoDatabase) HandleRequest {
	return HandleRequest{
		Uri:    "/account/create",
		Method: POST,
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

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateHashPassword(plainText string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	return string(bytes), err
}
