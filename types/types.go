package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type (
	BaseModel struct {
		ID        primitive.ObjectID `bson:"_id" json:"_id"`
		CreatedAt time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
		DeletedAt time.Time          `bson:"deleted_at" json:"deleted_at"`
	}
	Address struct {
		Address string
	}
	PasswordEnc struct {
		Password string `bson:"password" json:"-"`
		Salt     string `bson:"salt" json:"-"`
	}
	User struct {
		BaseModel   `bson:"-,inline"`
		FirstName   string      `bson:"first_name,omitempty" json:"first_name" validate:"required"`
		LastName    string      `bson:"last_name,omitempty" json:"last_name" validate:"required"`
		Email       string      `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"`
		Phone       string      `bson:"phone,omitempty" json:"phone,omitempty" validate:"required"`
		Password    PasswordEnc `bson:"password"`
		DateOfBirth time.Time   `bson:"date_of_birth,omitempty" json:"date_of_birth,omitempty"`
		Roles       []string    `bson:"roles,omitempty" json:"roles"`
		Address     Address     `bson:"address,inline,omitempty" json:"address,omitempty"`
	}
	Token struct {
	}
)
