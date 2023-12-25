package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"quickstart-go-jwt-mongodb/internal"
)

type tokenRepo struct {
	mongoDb internal.MongoDatabase
}

func (t *tokenRepo) CreateOne(context context.Context, model interface{}) (primitive.ObjectID, error) {
	//TODO implement me
	panic("implement me")
}

func (t *tokenRepo) CreateMany(context context.Context, model []interface{}) ([]primitive.ObjectID, error) {
	//TODO implement me
	panic("implement me")
}

func NewTokenRepository(mongoDb internal.MongoDatabase) CrudOperation {
	return &tokenRepo{
		mongoDb: mongoDb,
	}
}
