package controllers

import (
	"context"
	"server/cmd/api/database"
	"server/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/** Permet de recuperer tous les joueurs dans la base de donnees */
func GetPlayers() []models.PlayerCard {
	client, db := database.ConnectDB()
	// defer -> executed when the main func return
	defer client.Disconnect(context.TODO())
	playerColl := db.Collection("all_players")

	/* Find player */
	cur, currErr := playerColl.Find(context.TODO(), bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(context.TODO())

	var players []models.PlayerCard
	if err := cur.All(context.TODO(), &players); err != nil {
		panic(err)
	}

	return players
}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
