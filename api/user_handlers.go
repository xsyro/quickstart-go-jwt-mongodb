package api

import (
	"context"
	"net/http"
	"quickstart-go-jwt-mongodb/internal"
	"quickstart-go-jwt-mongodb/repositories"
	"quickstart-go-jwt-mongodb/types"
	"time"
)

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

func UserHandlers(resources *WithResource) {

	resources.HttpRequest.RequestRegistry(listCustomers(resources.MongoDatabase))

}
