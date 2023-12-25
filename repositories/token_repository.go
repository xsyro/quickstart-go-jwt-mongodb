package repositories

import "quickstart-go-jwt-mongodb/internal"

type tokenRepo struct {
	mongoDb *internal.MongoDatabase
}

func (t *tokenRepo) CreateOne(model *interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (t *tokenRepo) CreateMany(model *[]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func NewTokenRepository(mongoDb *internal.MongoDatabase) CrudOperation {
	return &tokenRepo{
		mongoDb: mongoDb,
	}
}
