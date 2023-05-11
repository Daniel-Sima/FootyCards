package dbrepo

import (
	"context"
	"log"
	"server/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
type DBRepo struct {
	DB *mongo.Database
}

/*--------------------------------------------------------------------------------------------------------*/
/** Variable stockant le temps d'attente afin l'arret d'une requete. */
const dbTimeout = time.Second * 3

/*--------------------------------------------------------------------------------------------------------*/
/** Getteur de la base de donnes. */
func (db *DBRepo) Connection() *mongo.Database {
	return db.DB
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction qui retourne l'user avec comme email 'email'. */
func (db *DBRepo) GetUserByEmail(email string) (*models.User, error) {
	// Setting context
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := db.DB.Collection("real_users")
	var user models.User

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction qui retourne l'user avec le id 'id'. */
func (db *DBRepo) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	// Setting context
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := db.DB.Collection("real_users")
	var user models.User

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

/*--------------------------------------------------------------------------------------------------------*/
/** Getting all players from the data base */
func (db *DBRepo) GetAllPlayersDB() ([]models.PlayerCard, error) {
	log.Println("Get all players from the Data Base")

	// Setting context
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := db.DB.Collection("all_players")

	// Compte le nombre de documents dans la collection
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	log.Printf("Size 'all_players' collection: %d", count)

	var players []models.PlayerCard

	// if err = cursor.All(ctx, &players); err != nil {
	// 	return nil, err
	// }

	for cursor.Next(ctx) {
		var player models.PlayerCard
		err := cursor.Decode(&player)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
