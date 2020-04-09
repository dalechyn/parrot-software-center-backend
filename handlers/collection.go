package handlers

import (
	"go.mongodb.org/mongo-driver/mongo"
	"parrot-software-center-backend/db"
)

var packageCollection *mongo.Collection

func init() {
	// Initializing dbClient for further package usage
	dbClient := db.MongoInit().Database("parrot_software_center")
	packageCollection = dbClient.Collection("packages")
}