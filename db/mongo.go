package db

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"time"
)

const mongoAddress = "localhost:27017"
const connectTimeout = 10 * time.Second

func MongoInit() *mongo.Client {
	log.WithField("ConnectTimeout", connectTimeout).Info("Attempting to connect to MongoDB")
	ctx, _ := context.WithTimeout(context.Background(), connectTimeout)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", mongoAddress)))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}

	log.Info("Connected to MongoDB")
	return client
}