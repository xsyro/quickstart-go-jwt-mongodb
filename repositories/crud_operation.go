package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CrudOperation interface {
	CreateOne(context context.Context, model interface{}) (primitive.ObjectID, error)
	CreateMany(context context.Context, model []interface{}) ([]primitive.ObjectID, error)
}
