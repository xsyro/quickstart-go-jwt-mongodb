package main

import (
	log "github.com/sirupsen/logrus"
	"quickstart-go-jwt-mongodb/api"
	"quickstart-go-jwt-mongodb/internal"
)

func main() {

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

	apiResources := api.WithResource{
		HttpRequest:   httpRequestHandler,
		MongoDatabase: mongoDb,
	}

	api.AuthHandlers(&apiResources)
	api.UserHandlers(&apiResources)

	httpRequestHandler.Serve()

}
