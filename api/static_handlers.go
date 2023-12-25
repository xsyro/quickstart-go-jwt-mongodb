package api

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"quickstart-go-jwt-mongodb/internal"
	"quickstart-go-jwt-mongodb/types"
)

func index(mongoDb internal.MongoDatabase) HandleRequest {
	return HandleRequest{
		Uri:    "/",
		Method: GET,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			one, err := mongoDb.Collection("users").InsertOne(context.Background(), types.User{})
			if err != nil {
				log.Error(err)
				return
			}
			log.Info(one)
		},
	}
}

// healthCheck return HTTP_STATUS_CODE of 200 to inform container service of it activeness
func healthCheck() HandleRequest {
	return HandleRequest{
		Uri:    "/status",
		Method: GET,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Good!"))
			return
		},
	}
}

func StaticHandlers(resources *WithResource) {
	healthCheck()
	index(resources.MongoDatabase)
}
