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
func InitializeMongoClient() (*mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoDBURI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	return client, nil
}

func getCollection() *mongo.Collection {
	return MongoClient.Database(db).Collection(collection)
}
