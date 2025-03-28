// services/mongo_service.go
package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ricardomussett/gotest/models"
)

type MongoService struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewMongoService(uri, dbName, collection string) (*MongoService, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collectionObj := client.Database(dbName).Collection(collection)
	return &MongoService{
		Client:     client,
		Collection: collectionObj,
	}, nil
}

func (ms *MongoService) SaveGPSData(data *models.GPSData) error {
	_, err := ms.Collection.InsertOne(context.Background(), data)
	return err
}
