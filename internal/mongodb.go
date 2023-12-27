package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type (
	MongoClient struct {
		client   *mongo.Client
		Database *mongo.Database
		ctx      context.Context
		timeOut  time.Duration
	}

	// MongoDatabase Implicitly import mongo.Client functions
	MongoDatabase interface {
		Client() *mongo.Client
		CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error
		Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
		RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *mongo.SingleResult
	}
)

const (
	maxRetryCount = 5
	appName       = "go-jwt-mongo"
)

var (
	currentRetryCount = 1
)

// NewMongoDbConn Singleton instance
// Reuse this client. You can use this same client instance to perform multiple tasks, instead of creating a new one each time.
// The client type is safe for concurrent use by multiple goroutines
func NewMongoDbConn() *MongoClient {
	var (
		err         error
		mongoUri    = os.Getenv("MONGO_DB_URI")
		timeout     = 6 * time.Second
		client      *mongo.Client
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	)
	defer func() {
		cancel()
	}()
	if err = godotenv.Load(); err != nil {
		log.Warn("No .env file found")
	}

	if mongoUri == "" {
		err = errors.New("no `MONGO_DB_URI` found in the environment variable. Please export your connection string to var MONGO_DB_URI to connect to the db")
		//panic(err)
	}

	mongoOption := options.Client()
	mongoOption.SetAppName(appName)
	mongoOption.ApplyURI(mongoUri)
	client, err = mongo.Connect(ctx, mongoOption)
	err = client.Ping(ctx, nil)

	if err != nil && currentRetryCount <= maxRetryCount {
		log.Errorf("[MongoDB] connect attempt failed. ActiveRequest retry %d of %d.\n%v", currentRetryCount, maxRetryCount, err)
		currentRetryCount++
		NewMongoDbConn()
	}
	if currentRetryCount >= maxRetryCount {
		//Throw panic if retries exceeded without any successful connection
		panic(fmt.Sprintf("Connection retries exceeded limit"))
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		panic(fmt.Sprintf("Unable to connect to MongoDB Server with the context timeout %v", timeout))
	}
	log.Debugf("successfully connected to MongoDB %s", mongoUri)

	var database = client.Database(os.Getenv("MONGO_DB_NAME"))
	return &MongoClient{
		client:   client,
		Database: database,
		ctx:      ctx,
		timeOut:  timeout,
	}
}

func (c *MongoClient) CloseClient() {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("No active MongoDB connection to close.")
		}
	}()
	_ = c.client.Disconnect(c.ctx)
}
