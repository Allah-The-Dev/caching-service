package config

import (
	"fmt"
	"log"
	"os"

	"github.com/adammck/venv"
)

const (
	mongoDBURIStr = "mongodb://%s:%s@%s/?authSource=admin&readPreference=primary&ssl=false"
)

var (
	//MongoDBURI ...
	MongoDBURI string
	//RedisURI ...
	RedisURI string
	//KafkaHost ...
	KafkaHost string
	//EmpAPILogger ...
	EmpAPILogger *log.Logger
)

//InitializeAppConfig ...
func InitializeAppConfig(env venv.Env) {
	//logger
	EmpAPILogger = log.New(os.Stdout, "employee-api : ", log.LstdFlags)

	//mongo db config
	dbServer := env.Getenv("MONGODB_SERVER")
	dbUsername, dbPassword := env.Getenv("MONGODB_ADMINUSERNAME"), env.Getenv("MONGODB_ADMINPASSWORD")
	MongoDBURI = fmt.Sprintf(mongoDBURIStr, dbUsername, dbPassword, dbServer)
	EmpAPILogger.Printf("mongodb server URI is : %s", dbServer)

	//redis config
	RedisURI = fmt.Sprintf("%s:%s", env.Getenv("REDIS_SERVER"), env.Getenv("REDIS_PORT"))

	//kafka config
	KafkaHost = env.Getenv("KAFKA_SERVER")

	EmpAPILogger.Printf("redis, kafka : %s %s", RedisURI, KafkaHost)
}
