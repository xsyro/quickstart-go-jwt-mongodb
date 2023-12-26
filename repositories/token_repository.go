package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"quickstart-go-jwt-mongodb/internal"
)

type tokenRepo struct {
	mongoDb    internal.MongoDatabase
	collection string
}

func (t *tokenRepo) FindOne(context context.Context, model interface{}, filters ...Filter) bool {
	//TODO implement me
	panic("implement me")
}

func (t *tokenRepo) CreateOne(context context.Context, model interface{}) (primitive.ObjectID, error) {
	id, err := t.mongoDb.Collection(t.collection).InsertOne(context, model)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return id.InsertedID.(primitive.ObjectID), nil
}

func (t *tokenRepo) CreateMany(context context.Context, model []interface{}) ([]primitive.ObjectID, error) {
	//TODO implement me
	panic("implement me")
}

func NewTokenRepository(mongoDb internal.MongoDatabase) CrudOperation {
	return &tokenRepo{
		mongoDb:    mongoDb,
		collection: "tokens",
	}
}
