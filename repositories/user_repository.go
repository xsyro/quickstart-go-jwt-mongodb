package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"quickstart-go-jwt-mongodb/internal"
)

type userRepo struct {
	mongoDb    internal.MongoDatabase
	collection string
	timeout    context.Context
	cancelFun  context.CancelFunc
}

func (u *userRepo) FindOne(context context.Context, model interface{}, filters ...Filter) bool {
	singleResult := u.mongoDb.Collection(u.collection).FindOne(context, filterToBsonFilter(filters...))
	_ = singleResult.Decode(&model)
	return !errors.Is(singleResult.Err(), mongo.ErrNoDocuments)
}

func (u *userRepo) CreateOne(context context.Context, model interface{}) (primitive.ObjectID, error) {
	id, err := u.mongoDb.Collection(u.collection).InsertOne(context, model)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return id.InsertedID.(primitive.ObjectID), nil
}

func (u *userRepo) CreateMany(context context.Context, model []interface{}) ([]primitive.ObjectID, error) {

	return nil, nil
}

func NewUserRepository(mongoDb internal.MongoDatabase) CrudOperation {
	return &userRepo{
		mongoDb:    mongoDb,
		collection: "users",
	}
}
