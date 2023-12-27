package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

type jwtService struct {
	//Tokenizer   interface{}
	//TokenString string
	ctx context.Context
}

var (
	JwtSigningAlg = jwt.SigningMethodHS256
	jwtSigningKey = []byte(os.Getenv("JWT_SECRET"))
)

func NewJwtService(ctx context.Context) *jwtService {
	return &jwtService{
		ctx: ctx,
	}
}

func (j *jwtService) GenerateJWT(payload interface{}, expiresAt time.Duration) (string, error) {
	token := jwt.New(JwtSigningAlg)
	claims := token.Claims.(jwt.MapClaims)

	claims["obj"] = payload
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(expiresAt).Unix()

	log.Infoln(claims)
	tokenStr, err := token.SignedString(jwtSigningKey)

	if err != nil {
		log.Error("error creating jwt token", err)
		return "", errors.New("error creating jwt token")
	}
	return tokenStr, nil
}

func (j *jwtService) ParseJWT(tokenizedString string) (map[string]any, error) {
	jwtToken, err := jwt.Parse(tokenizedString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return jwtSigningKey, nil
	})
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		obj := make(map[string]any)
		for k, v := range claims {
			obj[k] = v
		}
		return obj, nil
	}

	return nil, err
}
