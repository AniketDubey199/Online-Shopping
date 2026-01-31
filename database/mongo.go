package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database

func InitializeDb() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uri := os.Getenv("MONGO_URI")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientoptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, clientoptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Cannot connect to the database")
	}
	Db = client.Database("OnlineShopping")

	fmt.Printf("Connected to mongodb")

	return Db
}

func UserData(client *mongo.Client, usercollection string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("OnlineShopping").Collection(usercollection)
	return collection
}

func ProdData(client *mongo.Client, productCollection string) *mongo.Collection {
	var prodcollection *mongo.Collection = client.Database("OnlineShopping").Collection(productCollection)
	return prodcollection
}
