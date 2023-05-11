package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/cmd/api/database"
	"server/internal/models"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/** Recuperer l'ensemble des joueurs au debut du projet pour le mettre a jour */
func GetAllPlayersAPI() {
	callAPI(1)
}

/*--------------------------------------------------------------------------------------------------------*/
func callAPI(paging int32) {

	clientDB, db := database.ConnectDB()
	// defer -> executed when the main func return
	defer clientDB.Disconnect(context.TODO())
	playerColl := db.Collection("all_players2") // 2 pour pas ecraser l'originale

	var url string
	if paging != 1 {
		url = "https://v3.football.api-sports.io/players?league=61&season=2022&page=" + strconv.FormatInt(int64(paging), 10)
	} else {
		url = "https://v3.football.api-sports.io/players?league=61&season=2022"
	}
	// &paging=" + strconv.FormatInt(int64(paging), 10)
	fmt.Println(url)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	// req.Header.Add("x-rapidapi-key", "9cfa46fb2a5914ad838fe1dec5a5d9f9")
	req.Header.Add("x-rapidapi-key", "653e109f8c7a6ff2115193ae42617519")
	req.Header.Add("x-rapidapi-host", "v3.football.api-sports.io")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data models.RequestJSON
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	var players []interface{}
	for _, elem := range data.Response {
		var rating float64
		var price int32
		if elem.Statistics[0].Games.Rating != "" {
			rating, err = strconv.ParseFloat(elem.Statistics[0].Games.Rating, 32)
			price = int32(rating * 100)
			if err != nil {
				panic(err)
			}
		} else {
			rating = -1
			price = int32(150)
		}
		player := models.PlayerCard{
			ID_API:       int32(elem.Player.Id),
			Name:         elem.Player.Name,
			RatingGlobal: float32(rating),
			LastRating:   float32(rating),
			Image:        elem.Player.Photo,
			Price:        price,
			Club:         elem.Statistics[0].Team.Name,
			Position:     elem.Statistics[0].Games.Position,
			Exemplars:    10, // par defaut
			LastUpdate:   time.Now(),
		}
		players = append(players, player)
	}

	playerColl.InsertMany(context.TODO(), players)

	if data.Paging.Current == 1 {
		// Repeat if is not the last page
		for i := 2; int32(i) <= data.Paging.Total; i++ {
			callAPI(int32(i))
		}
	}
}

/*--------------------------------------------------------------------------------------------------------*/
// Update user score and players depending on players IRL performance
func UsersRaking() {

	fixturesIds := getFixtures()

	playersRating := make(map[int32]string)
	for _, elem := range fixturesIds {
		playersFixtureRating := getPlayerFixtureRating(int32(elem))
		for k, v := range playersFixtureRating {
			playersRating[k] = v
		}
	}

	// Update users and players in users team
	users := GetAllUsers()
	// For all users
	for id, _ := range users {
		// For all players in user team
		for id_playerTeam, playerTeam := range users[id].Team {
			// For all players that have played this week
			for idPlayersRating, rating := range playersRating {
				if playerTeam.ID_API == idPlayersRating {
					// If the player have a rating
					if rating != "" {
						if r, err := strconv.ParseFloat(rating, 32); err == nil {
							if users[id].LastUpdate.Day() != time.Now().Day() {
								users[id].Score = float32(r)
								users[id].LastUpdate = time.Now()
							} else {
								users[id].Score += float32(r)
							}
							users[id].Team[id_playerTeam].LastRating = float32(r)
							users[id].Team[id_playerTeam].RatingGlobal = (users[id].Team[id_playerTeam].RatingGlobal + float32(r)) / 2
							users[id].Team[id_playerTeam].Price = int32(users[id].Team[id_playerTeam].RatingGlobal * 100)
							users[id].Team[id_playerTeam].LastUpdate = time.Now()
						}
					} else {
						if users[id].LastUpdate.Day() != time.Now().Day() {
							users[id].Score = float32(0)
							users[id].LastUpdate = time.Now()
						}
					}
				} else {
					if users[id].LastUpdate.Day() != time.Now().Day() {
						users[id].Score = float32(0)
						users[id].LastUpdate = time.Now()
					}
				}
			}
		}
	}

	// Update player db
	players := GetPlayers()
	for id, player := range players {
		for idPlayerRating, rating := range playersRating {
			if player.ID_API == idPlayerRating {
				if r, err := strconv.ParseFloat(rating, 32); err == nil {
					players[id].LastRating = float32(r)
					players[id].RatingGlobal = (players[id].RatingGlobal + float32(r)) / 2
					players[id].Price = int32(players[id].RatingGlobal * 100)
					players[id].LastUpdate = time.Now()
				}
			}
		}
	}

	// Update users in DB
	client, db := database.ConnectDB()
	// defer -> executed when the main func return
	defer client.Disconnect(context.TODO())
	userColl := db.Collection("real_users")

	for _, user := range users {
		for _, player := range user.Team {
			filter := bson.D{{"_id", user.ID}}
			update := bson.D{
				{"$set", bson.D{{"score", user.Score}}},
				{"$set", bson.D{{"LastUpdate", user.LastUpdate}}},
				{"$set", bson.D{{"club.$[element].lastRating", player.LastRating}}},
				{"$set", bson.D{{"club.$[element].ratingGlobal", player.RatingGlobal}}},
				{"$set", bson.D{{"club.$[element].price", player.Price}}},
				{"$set", bson.D{{"club.$[element].lastupdate", player.LastUpdate}}},
				{"$set", bson.D{{"team.$[element].lastRating", player.LastRating}}},
				{"$set", bson.D{{"team.$[element].ratingGlobal", player.RatingGlobal}}},
				{"$set", bson.D{{"team.$[element].price", player.Price}}},
				{"$set", bson.D{{"team.$[element].lastupdate", player.LastUpdate}}},
			}
			arrayFiltersOpt := options.ArrayFilters{Filters: []interface{}{bson.D{{"element._id", player.ID}}}}
			opts := options.Update().SetArrayFilters(arrayFiltersOpt)
			userColl.UpdateOne(context.TODO(), filter, update, opts)
		}
	}

	// playerColl := db.Collection("player")
	playerColl := db.Collection("all_players")

	for _, player := range players {
		filter := bson.D{{"_id", player.ID}}
		update := bson.D{
			{"$set", bson.D{{"lastRating", player.LastRating}}},
			{"$set", bson.D{{"ratingGlobal", player.RatingGlobal}}},
			{"$set", bson.D{{"price", player.Price}}},
			{"$set", bson.D{{"lastupdate", player.LastUpdate}}},
		}
		playerColl.UpdateOne(context.TODO(), filter, update)
	}

	fmt.Println("Ranked and Players updated")
}

/*--------------------------------------------------------------------------------------------------------*/
// Get last 10 fixtures from API
func getFixtures() []int32 {

	url := "https://v3.football.api-sports.io/fixtures?league=61&last=10"

	fmt.Println(url)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	// req.Header.Add("x-rapidapi-key", "653e109f8c7a6ff2115193ae42617519")
	req.Header.Add("x-rapidapi-key", "653e109f8c7a6ff2115193ae42617519")
	req.Header.Add("x-rapidapi-host", "v3.football.api-sports.io")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data models.Fixtures
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	var ids []int32
	for _, elem := range data.Response {
		ids = append(ids, elem.Fixture.Id)
	}
	fmt.Println("getFixtures done")
	return ids
}

/*--------------------------------------------------------------------------------------------------------*/
// Get players stats for the last 10 fixtures from API
func getPlayerFixtureRating(id int32) map[int32]string {
	url := "https://v3.football.api-sports.io/fixtures/players?fixture=" + strconv.FormatInt(int64(id), 10)

	fmt.Println(url)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	// req.Header.Add("x-rapidapi-key", "653e109f8c7a6ff2115193ae42617519")
	req.Header.Add("x-rapidapi-key", "653e109f8c7a6ff2115193ae42617519")
	req.Header.Add("x-rapidapi-host", "v3.football.api-sports.io")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data models.FixturePlayer
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	playersFixtureRating := make(map[int32]string)
	for _, elem := range data.Response {
		for _, player := range elem.Players {
			playersFixtureRating[int32(player.Player.Id)] = player.Statistics[0].Games.Rating
		}
	}
	fmt.Println("getPlayerFixtureRating done")
	return playersFixtureRating
}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
