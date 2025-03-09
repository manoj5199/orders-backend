package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var CustomerCollection *mongo.Collection
var ProductCollection *mongo.Collection
var OrderCollection *mongo.Collection

func ConnectDB(uri, dbName, customerCollectionName, productCollectionName, orderCollectionName string) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	Client = client
	CustomerCollection = client.Database(dbName).Collection(customerCollectionName)
	ProductCollection = client.Database(dbName).Collection(productCollectionName)
	OrderCollection = client.Database(dbName).Collection(orderCollectionName)
}
