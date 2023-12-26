package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Filter struct {
	Key   string
	Value interface{}
}
type CrudOperation interface {
	CreateOne(context context.Context, model interface{}) (primitive.ObjectID, error)
	CreateMany(context context.Context, model []interface{}) ([]primitive.ObjectID, error)
	FindOne(context context.Context, model interface{}, filters ...Filter) bool
}

func filterToBsonFilter(filters ...Filter) bson.D {
	f := bson.D{}
	for i := range filters {
		f = append(f, bson.E{Key: filters[i].Key, Value: filters[i].Value})
	}
	return f
}
