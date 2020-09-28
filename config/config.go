package config

import (
	"fmt"
	"log"
	"os"
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
func InitializeAppConfig() {
	//logger
	EmpAPILogger = log.New(os.Stdout, "employee-api : ", log.LstdFlags)

	//mongo db config
	dbServer := os.Getenv("MONGODB_SERVER")
	dbUsername, dbPassword := os.Getenv("MONGODB_ADMINUSERNAME"), os.Getenv("MONGODB_ADMINPASSWORD")
	MongoDBURI = fmt.Sprintf(mongoDBURIStr, dbUsername, dbPassword, dbServer)
	EmpAPILogger.Printf("mongodb server URI is : %s", dbServer)

	//redis config
	RedisURI = fmt.Sprintf("%s:%s", os.Getenv("REDIS_SERVER"), os.Getenv("REDIS_PORT"))

	//kafka config
	KafkaHost = os.Getenv("KAFKA_SERVER")

	EmpAPILogger.Printf("redis, kafka : %s %s", RedisURI, KafkaHost)
}
