package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/cmd/api/database"
	"server/internal/models"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func GetUserPlayers(w http.ResponseWriter, r *http.Request) {
// 	middleware.EnableCors(&w)
// 	client, db := database.ConnectDB()
// 	// defer -> executed when the main func return
// 	defer client.Disconnect(context.TODO())
// 	// playerColl := db.Collection("player")
// 	userColl := db.Collection("user")

// 	/* Find user */
// 	usrId, err := primitive.ObjectIDFromHex("644ec501ba123d29bfb751cc")
// 	if err != nil {
// 		panic(err)
// 	}

// 	filter := bson.D{{"_id", usrId}}
// 	opts := options.FindOne().SetProjection(bson.M{"Club": 1, "_id": 0})

// 	var club bson.M
// 	if err = userColl.FindOne(context.TODO(), filter, opts).Decode(&club); err != nil {
// 		log.Fatal(err)
// 	}

// 	// var user models.User
// 	// userColl.FindOne(context.TODO(), bson.D{{"_id", usrId}}).Decode(&user)

// 	// fmt.Println("here")

// 	// /* Find all user players */
// 	// var club []models.PlayerCard
// 	// for _, elem := range user.Club {
// 	// 	var player models.PlayerCard
// 	// 	playerColl.FindOne(context.TODO(), bson.D{{"_id", elem}}).Decode(&player)
// 	// 	club = append(club, player)
// 	// }

// 	out, err := json.Marshal(club)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(out)
// }

// func GetUserTeam(w http.ResponseWriter, r *http.Request) {
// 	middleware.EnableCors(&w)
// 	client, db := database.ConnectDB()
// 	// defer -> executed when the main func return
// 	defer client.Disconnect(context.TODO())
// 	// playerColl := db.Collection("player")
// 	userColl := db.Collection("user")

// 	/* Find user */
// 	usrId, err := primitive.ObjectIDFromHex("644ec501ba123d29bfb751cc")
// 	if err != nil {
// 		panic(err)
// 	}
// 	filter := bson.D{{"_id", usrId}}
// 	opts := options.FindOne().SetProjection(bson.M{"Team": 1, "_id": 0})

// 	var team bson.M
// 	if err = userColl.FindOne(context.TODO(), filter, opts).Decode(&team); err != nil {
// 		log.Fatal(err)
// 	}

// 	out, err := json.Marshal(team)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(out)
// }

// func GetUser(w http.ResponseWriter, r *http.Request) {
// 	client, db := database.ConnectDB()
// 	// defer -> executed when the main func return
// 	defer client.Disconnect(context.TODO())
// 	userColl := db.Collection("real_users")

// 	/* Find user */
// 	usrId, err := primitive.ObjectIDFromHex("6457cb803e2c827cc0a44ecc")
// 	if err != nil {
// 		panic(err)
// 	}
// 	filter := bson.D{{"_id", usrId}}

// 	var user models.User
// 	if err = userColl.FindOne(context.TODO(), filter).Decode(&user); err != nil {
// 		log.Fatal(err)
// 	}

// 	out, err := json.Marshal(user)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(out)
// 	fmt.Println("GetUser ended")
// }

// Get all users from the DB
func GetAllUsers() []models.User {
	client, db := database.ConnectDB()
	// defer -> executed when the main func return
	defer client.Disconnect(context.TODO())
	userColl := db.Collection("real_users")

	var results []models.User
	cur, err := userColl.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem models.User
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return results
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler POST permettant d'ajouter un jouer dans la Team */
func AddPlayerTeam(w http.ResponseWriter, r *http.Request) {
	client, db := database.ConnectDB()
	// defer -> executed when the main func return
	defer client.Disconnect(context.TODO())
	userColl := db.Collection("real_users")

	var req models.UserClubReq
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatal(err)
	}

	req.Player.InTeam = true
	filter := bson.D{{"_id", req.UserID}}
	update := bson.D{
		{"$push", bson.D{{"team", req.Player}}},
		{"$set", bson.D{{"club.$[element]", req.Player}}},
	}
	arrayFiltersOpt := options.ArrayFilters{Filters: []interface{}{bson.D{{"element._id", req.Player.ID}}}}
	opts := options.Update().SetArrayFilters(arrayFiltersOpt)

	var usr models.User
	userColl.FindOne(context.TODO(), filter).Decode(&usr)
	if len(usr.Team) < 4 {
		userColl.UpdateOne(context.TODO(), filter, update, opts)
		w.WriteHeader(http.StatusOK)
		fmt.Println("AddPlayerTeam OK")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("AddPlayerTeam Bad request")
	}
	w.Header().Set("Content-Type", "application/json")
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler DELETE permettant d'enlever un jouer de la Team */
func DeletePlayerTeam(w http.ResponseWriter, r *http.Request) {
	client, db := database.ConnectDB()
	// defer -> executed when the main func return
	defer client.Disconnect(context.TODO())
	userColl := db.Collection("real_users")

	var req models.UserClubReq
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatal(err)
	}

	// req.Player.InTeam = false
	filter := bson.D{{"_id", req.UserID}}
	update := bson.D{
		{"$pull", bson.D{{"team", req.Player}}},
		{"$set", bson.D{{"club.$[element].inteam", false}}},
	}
	arrayFiltersOpt := options.ArrayFilters{Filters: []interface{}{bson.D{{"element._id", req.Player.ID}}}}
	opts := options.Update().SetArrayFilters(arrayFiltersOpt)

	userColl.UpdateOne(context.TODO(), filter, update, opts)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("DeletePlayerTeam ended")
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler GET permettant de recuperer les Users tries par en fonction du score */
func GetRanking(w http.ResponseWriter, r *http.Request) {

	users := GetAllUsers()

	sort.Slice(users, func(i, j int) bool { return float32(users[i].Score) > float32(users[j].Score) }) // triage par score

	// Get the 10 first player ranked
	out, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
	fmt.Println("GetRanking ended")

}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
