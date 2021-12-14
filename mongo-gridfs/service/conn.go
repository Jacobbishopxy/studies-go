package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 连接 MongoDB
func InitiateMongoClient() *mongo.Client {
	var err error
	var client *mongo.Client

	uri := "mongodb://admin:password@localhost:27017"
	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetMaxPoolSize(5)

	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		fmt.Println(err.Error())
	}

	return client
}
