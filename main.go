package main

import (
	log "github.com/sirupsen/logrus"
	"quickstart-go-jwt-mongodb/internal"
)

func main() {

	var (
		mongoClient *internal.MongoClient
	)

	defer func() {
		mongoClient.CloseClient()
	}()

	mongoClient = internal.NewMongoDbConn()

	if r := recover(); r != nil {
		log.Warnf("[RECOVERY_FROM_FAILURE] %v", r)
	}
}
