package repositories

import "quickstart-go-jwt-mongodb/internal"

type userRepo struct {
	mongoDb *internal.MongoDatabase
}

func (u *userRepo) CreateOne(model *interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) CreateMany(model *[]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func NewUserRepository(mongoDb *internal.MongoDatabase) CrudOperation {
	return &userRepo{
		mongoDb: mongoDb,
	}
}
