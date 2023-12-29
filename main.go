package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"quickstart-go-jwt-mongodb/internal"
	"quickstart-go-jwt-mongodb/middleware"
	"quickstart-go-jwt-mongodb/route"
	"quickstart-go-jwt-mongodb/server"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)
	var (
		mongoClient        *internal.MongoClient
		httpRequestHandler server.RequestHandler
	)
	parentHttpCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(40*time.Second))
	httpRequestHandler = server.NewHttpRequestHandler(parentHttpCtx)

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

	route.Routes(httpRequestHandler, mongoDb, parentHttpCtx)

	//Use middleware to intermediate every requests
	httpRequestHandler.HandleMiddlewares(
		middleware.HeadersMiddleware(),
		middleware.SecureMiddleware(httpRequestHandler.GetControllers()),
	)

	httpRequestHandler.Serve()

}
