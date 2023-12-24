package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"quickstart-go-jwt-mongodb/api"
	"quickstart-go-jwt-mongodb/internal"
	"time"
)

func main() {

	var (
		mongoClient        *internal.MongoClient
		httpRequestHandler *api.HttpRequestHandler
		ctx, cancel        = context.WithTimeout(context.Background(), 30*time.Second)
	)

	defer func() {
		log.Debug("HTTP PORT TERMINATING...CLOSING RESOURCES!!!")
		mongoClient.CloseClient()
		cancel()
	}()

	mongoClient = internal.NewMongoDbConn()

	if r := recover(); r != nil {
		log.Warnf("[RECOVERY_FROM_FAILURE] %v", r)
	}

	httpRequestHandler = api.NewHttpRequestHandler()

	httpRequestHandler.HandleRequest(api.HandleRequest{
		Uri:    "/",
		Method: api.GET,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("Hello!"))
			return
		},
	})

	//Use middleware to intermediate every requests
	apiResources := api.WithResource{
		HttpRequest: httpRequestHandler,
		Context:     &ctx,
		MongoClient: mongoClient,
	}

	log.Debug(apiResources)

	httpRequestHandler.Serve()

}
