package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	employeeCollection = "employee"
	db                 = "cacheService"
	collection         = "employee"
	MongoClient        *mongo.Client
)

func InitializeMongoClient(mongoDBURI string) (*mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	return client, nil
}

func GetCollection() *mongo.Collection {
	return MongoClient.Database(db).Collection(collection)
}
