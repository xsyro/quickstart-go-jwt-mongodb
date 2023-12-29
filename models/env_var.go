package models

import "os"

type EnvVar struct {
	HttpPort,
	MongoDbUri,
	MongoDbName,
	BaseUrlPrefix,
	JwtSecret,
	Value string
}

func LoadEnvironmentVariables() EnvVar {
	return EnvVar{
		HttpPort:      os.Getenv("HTTP_PORT"),
		MongoDbUri:    os.Getenv("MONGO_DB_URI"),
		MongoDbName:   os.Getenv("MONGO_DB_NAME"),
		BaseUrlPrefix: os.Getenv("BASE_URI_PREFIX"),
	}
}
