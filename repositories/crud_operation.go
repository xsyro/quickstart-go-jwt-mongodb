package repositories

type CrudOperation interface {
	CreateOne(model *interface{}) error
	CreateMany(model *[]interface{}) error
}
