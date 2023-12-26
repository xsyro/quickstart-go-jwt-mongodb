package main

import (
	log "github.com/sirupsen/logrus"
	"quickstart-go-jwt-mongodb/api"
	"quickstart-go-jwt-mongodb/internal"
	"quickstart-go-jwt-mongodb/middleware"
)

func main() {
	log.SetLevel(log.DebugLevel)
	var (
		mongoClient        *internal.MongoClient
		httpRequestHandler *api.HttpRequestHandler
	)

	defer func() {
		log.Debug("HTTP PORT TERMINATING...CLOSING RESOURCES!!!")
		mongoClient.CloseClient()
	}()

	mongoClient = internal.NewMongoDbConn()
	var mongoDb internal.MongoDatabase = mongoClient.Database

	if r := recover(); r != nil {
		log.Warnf("[RECOVERY_FROM_FAILURE] %v", r)
	}

	httpRequestHandler = api.NewHttpRequestHandler()

	//Use middleware to intermediate every requests
	httpRequestHandler.HandleMiddlewares(middleware.HeadersMiddleware())

	apiResources := api.WithResource{
		HttpRequest:   httpRequestHandler,
		MongoDatabase: mongoDb,
	}

	api.StaticHandlers(&apiResources)
	api.AuthHandlers(&apiResources)
	api.UserHandlers(&apiResources)

	httpRequestHandler.Serve()

}
