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
		httpRequestHandler api.RequestHandler
	)
	timeout, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	httpRequestHandler = api.NewHttpRequestHandler(timeout)

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

	apiResources := api.WithResource{
		HttpRequest:   httpRequestHandler,
		MongoDatabase: mongoDb,
	}

	api.StaticHandlers(&apiResources)
	api.AuthHandlers(&apiResources)
	api.UserHandlers(&apiResources)

	//Use middleware to intermediate every requests
	httpRequestHandler.HandleMiddlewares(middleware.HeadersMiddleware(), httpRequestHandler.SecureMiddleware())

	httpRequestHandler.Serve()

}
