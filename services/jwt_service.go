package services

import (
	"context"
	"encoding/json"
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

// GenerateJWT to generate a JWT Payload
// sub must not be a JSON string. This would be done by the GenerateJWT method
func (j *jwtService) GenerateJWT(sub interface{}, expiresAt time.Duration, extraClaims ...map[string]any) (string, error) {
	token := jwt.New(JwtSigningAlg)
	claims := token.Claims.(jwt.MapClaims)

	marshal, err := json.Marshal(sub)
	claims["sub"] = string(marshal)
	claims["iat"] = jwt.NewNumericDate(time.Now())
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(expiresAt))
	for i := range extraClaims {
		for k, v := range extraClaims[i] {
			claims[k] = v
		}
	}

	tokenStr, err := token.SignedString(jwtSigningKey)

	if err != nil {
		log.Error("error creating jwt token", err)
		return "", errors.New("error creating jwt token")
	}
	return tokenStr, nil
}

func (j *jwtService) ClaimToken(tokenizedString string, sub interface{}) (jwt.Claims, error) {
	jwtToken, err := jwt.Parse(tokenizedString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return jwtSigningKey, nil
	})
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	subStr, err := claims.GetSubject()
	err = json.Unmarshal([]byte(subStr), sub)
	if ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, err
}
