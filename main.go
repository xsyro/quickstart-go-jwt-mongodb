package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"quickstart-go-jwt-mongodb/api"
	"quickstart-go-jwt-mongodb/internal"
	"quickstart-go-jwt-mongodb/middleware"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)
	var (
		mongoClient        *internal.MongoClient
		httpRequestHandler *api.HttpRequestHandler
	)
	timeout, cancel := context.WithTimeout(context.Background(), 40*time.Second)

	defer func() {
		log.Debug("HTTP PORT TERMINATING...CLOSING RESOURCES!!!")
		mongoClient.CloseClient()
		cancel()
	}()

	mongoClient = internal.NewMongoDbConn()
	var mongoDb internal.MongoDatabase = mongoClient.Database

	if r := recover(); r != nil {
		log.Warnf("[RECOVERY_FROM_FAILURE] %v", r)
	}

	httpRequestHandler = api.NewHttpRequestHandler(timeout)
	apiResources := api.WithResource{
		HttpRequest:   httpRequestHandler,
		MongoDatabase: mongoDb,
	}

	//Use middleware to intermediate every requests
	httpRequestHandler.HandleMiddlewares(middleware.HeadersMiddleware())

	api.StaticHandlers(&apiResources)
	api.AuthHandlers(&apiResources)
	api.UserHandlers(&apiResources)

	httpRequestHandler.Serve()

}
