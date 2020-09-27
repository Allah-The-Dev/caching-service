package data

import (
	"caching-service/config"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	employeeCollection = "employee"
	db                 = "cacheService"
	collection         = "employee"
	//MongoClient ...
	MongoClient *mongo.Client
)

//InitializeMongoClient ...
func InitializeMongoClient() error {

	var err error
	MongoClient, err = mongo.NewClient(options.Client().ApplyURI(config.MongoDBURI))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	if err = MongoClient.Connect(ctx); err != nil {
		return err
	}
	return nil
}

func getCollection() *mongo.Collection {
	return MongoClient.Database(db).Collection(collection)
}
