package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type (
	BaseModel struct {
		ID        primitive.ObjectID `bson:"_id" json:"id"`
		CreatedAt time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
		DeletedAt time.Time          `bson:"deleted_at" json:"deleted_at"`
	}
	Address struct {
	}
	User struct {
		BaseModel
		FirstName   string    `bson:"first_name,omitempty" json:"first_name,omitempty"`
		LastName    string    `bson:"last_name,omitempty" json:"last_name,omitempty"`
		Email       string    `bson:"email,omitempty" json:"email,omitempty"`
		Phone       string    `bson:"phone,omitempty" json:"phone,omitempty"`
		DateOfBirth time.Time `bson:"date_of_birth,omitempty" json:"date_of_birth,omitempty"`
		Roles       []string  `bson:"roles,omitempty" json:"roles"`
		Address     Address   `bson:"address,inline,omitempty" json:"address,omitempty"`
	}
	Token struct {
	}
)
