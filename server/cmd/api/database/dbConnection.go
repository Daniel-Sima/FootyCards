package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*Connection to mongoDB*/
func ConnectDB() (*mongo.Client, *mongo.Database) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://footycard:1234@cluster0.nrojts7.mongodb.net/test"))
	if err != nil {
		log.Fatal(err)
	}
	fcDatabase := client.Database("footycardDB")

	log.Println("Connnected to the Data Base!")

	return client, fcDatabase
}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
