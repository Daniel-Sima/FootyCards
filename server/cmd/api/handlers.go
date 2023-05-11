package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"server/cmd/api/controllers"
	"server/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/** Handler POST permettant de faire l'authentification sur la page de login. */
func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("Authenticate")
	EnableCors(&w)
	// Read JSON payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid email"), http.StatusBadRequest)
		return
	}

	// Check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid password"), http.StatusBadRequest)
		return
	}

	log.Printf("authenticate user.ID:: %d", user.ID)
	log.Printf("authenticate user.Pseudo:: %s", user.Pseudo)
	log.Printf("authenticate user.Email:: %s", user.Email)

	// Create a JWT user
	u := jwtUser{
		ID:     user.ID.Hex(),
		Pseudo: user.Pseudo,
		// LastName:  user.LastName,
	}

	// Generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	cookieHeader := w.Header().Get("Set-Cookie")
	if cookieHeader == "" {
		log.Println("Failed to set cookie")
		// Traitement en cas d'erreur de définition du cookie
	} else {
		log.Printf("Cookie set: %s", cookieHeader)
		// Traitement en cas de succès de définition du cookie
	}

	app.writeJSON(w, http.StatusAccepted, tokens)
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction permettant d'afficher quelque chose sur le path "/". */
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Footy Cards is running well!",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler GET refreshing the cookie to maintain the connection by updating */
func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Refreshing cookies")
	EnableCors(&w)
	for _, cookie := range r.Cookies() {
		log.Printf("aololo cookies size=> %d", len(r.Cookies()))
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims
			// userID, err := strconv.Atoi(claims.Subject)
			// if err != nil {
			// 	app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
			// 	return
			// }
			userID := claims.Subject
			idHex, err := primitive.ObjectIDFromHex(userID)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(idHex)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user id"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:     user.ID.Hex(),
				Pseudo: user.Pseudo,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			app.writeJSON(w, http.StatusOK, tokenPairs)

		}
	}
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler GET removes the cookie when the user logout */
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler POST permettant de faire inscription sur la page de sign-up. */
func (app *application) register(w http.ResponseWriter, r *http.Request) {
	log.Println("Register")
	EnableCors(&w)
	// Read JSON payload
	var requestPayload struct {
		Pseudo   string `json:"pseudo"`
		ClubName string `json:"clubName"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check user against database
	_, err = app.DB.GetUserByEmail(requestPayload.Email)
	if err == nil {
		app.errorJSON(w, errors.New("already an account with this email"), http.StatusBadRequest)
		return
	}

	// Insert in the data base
	dbTimeout := time.Second * 3
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := app.DB.Connection().Collection("real_users")

	passwordHash, err := HashPassword(requestPayload.Password)
	if err != nil {
		return
	}

	var playerVide []models.PlayerCard

	newUser := models.User{
		Email:    requestPayload.Email,
		Password: passwordHash,
		Pseudo:   requestPayload.Pseudo,
		ClubName: requestPayload.ClubName,
		Coins:    4000, // car vient de s'inscrire
		Score:    0,    // car vient de s'inscrire
		// Classement: -1,   // TODO faire fct qui met a jour les classements
		Club:       playerVide,
		Team:       playerVide,
		CreatedAt:  time.Now().UTC(),
		LastUpdate: time.Now().UTC(),
	}

	// Convertit la struct User en un document BSON
	userDoc, err := bson.Marshal(newUser)
	if err != nil {
		app.errorJSON(w, errors.New("converting bson error"), http.StatusBadRequest)
		return
	}

	// Insère le document BSON dans la collection users
	_, err = collection.InsertOne(ctx, userDoc)
	if err != nil {
		app.errorJSON(w, errors.New("inserting bson error"), http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusAccepted, userDoc) // XXX
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction permettant de generer une hash d'un mot de passe.*/
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler GET permettant de recuperer tous les joueurs de la base de donnees. */
func (app *application) publish_market(w http.ResponseWriter, r *http.Request) {
	log.Println("Publish Market")
	EnableCors(&w)

	// Get all players from the database
	players, err := app.DB.GetAllPlayersDB()
	if err != nil {
		app.errorJSON(w, errors.New("getting all players from database error"), http.StatusBadRequest)
		return
	}

	log.Printf("Number of players harvested :: %d", len(players))

	app.writeJSON(w, http.StatusAccepted, players)
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler GET les information d'un user connecte (cookies).  */
func (app *application) user_infos_cookies(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting User Infos Cookies")
	EnableCors(&w)

	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims
			userID := claims.Subject
			// if err != nil {
			// 	app.errorJSON(w, errors.New("unknown user 00"), http.StatusUnauthorized)
			// 	return
			// }

			idHex, err := primitive.ObjectIDFromHex(userID)
			if err != nil {
				app.errorJSON(w, errors.New("object id from hex error"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(idHex)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user id"), http.StatusUnauthorized)
				return
			}

			log.Printf("User : %s has %d coins !", user.Pseudo, user.Coins)

			app.writeJSON(w, http.StatusAccepted, user)
		}
	}
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler POST pour recuperer les informations d'un user connecte avec l'email dans le localStorage si le cookies
ne fonctionnent pas.  */
func (app *application) user_infos(w http.ResponseWriter, r *http.Request) {
	log.Println("Get User Infos (using localStorage)")
	EnableCors(&w)
	// Read JSON payload
	var requestPayload struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid email"), http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusAccepted, user)
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler POST de l'achat d'un joueur sur la page de Market. */
func (app *application) buy_player(w http.ResponseWriter, r *http.Request) {
	log.Println("Buying a player")
	EnableCors(&w)

	// Read JSON user
	// var user models.User

	type RequestData struct {
		User   models.User       `json:"user"`
		Player models.PlayerCard `json:"player"`
	}

	var requestData RequestData

	err := app.readJSON(w, r, &requestData)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user := requestData.User
	player := requestData.Player

	// Check user against database
	userDB, err := app.DB.GetUserByEmail(user.Email)
	if err != nil {
		app.errorJSON(w, errors.New("no user with this email"), http.StatusBadRequest)
		return
	}

	// Check if the user hasn't already this player
	if contains(userDB.Club, player) {
		app.errorJSON(w, errors.New(player.Name+" is already in your club"), http.StatusBadRequest)
		return
	}

	// Check user coins
	if userDB.Coins < int64(player.Price) {
		app.errorJSON(w, errors.New("not enough money to buy this player"), http.StatusBadRequest)
		return
	}

	// Check is enough exemplars
	if player.Exemplars == 0 {
		app.errorJSON(w, errors.New("no more exemplars, wait for someone to sell one"), http.StatusBadRequest)
		return
	}

	// Insert in the data base
	dbTimeout := time.Second * 3
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := app.DB.Connection().Collection("real_users")

	// Copie du player pour avoir une nouvelle instance mongoDB
	newPlayer := models.PlayerCard{
		ID:           primitive.NewObjectID(),
		ID_API:       player.ID_API,
		Name:         player.Name,
		RatingGlobal: player.RatingGlobal,
		LastRating:   player.LastRating,
		Image:        player.Image,
		Price:        player.Price,
		Club:         player.Club,
		Position:     player.Position,
		LastUpdate:   time.Now(),
	}

	userMAJ, err := collection.UpdateOne(
		ctx,
		bson.M{"email": user.Email},
		bson.M{
			"$set":  bson.M{"coins": user.Coins - int64(player.Price)},
			"$push": bson.M{"club": newPlayer},
		},
	)
	if err != nil {
		app.errorJSON(w, errors.New("update coins error"), http.StatusBadRequest)
		return
	}

	// Diminution du nombre d'exemplaires disponibles
	collectionAllPlayers := app.DB.Connection().Collection("all_players")
	_, err = collectionAllPlayers.UpdateOne(
		ctx,
		bson.M{"_id": player.ID},
		bson.M{
			"$set": bson.M{"exemplars": player.Exemplars - 1},
		},
	)
	if err != nil {
		app.errorJSON(w, errors.New("update exemplar numbers error"), http.StatusBadRequest)
		return
	}

	// log.Printf("Player buyed: %s", player.Name)

	// Convertit la struct User en un document BSON
	userDoc, err := bson.Marshal(userMAJ)
	if err != nil {
		app.errorJSON(w, errors.New("converting bson error"), http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusAccepted, userDoc) // XXX
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler POST de la vente d'un joueur sur la page de Club. */
func (app *application) selling_player(w http.ResponseWriter, r *http.Request) {
	log.Println("Selling a player")
	EnableCors(&w)

	// Read JSON user and player
	type RequestData struct {
		User   models.User       `json:"user"`
		Player models.PlayerCard `json:"player"`
	}

	var requestData RequestData

	err := app.readJSON(w, r, &requestData)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user := requestData.User
	player := requestData.Player

	// Check user against database
	userDB, err := app.DB.GetUserByEmail(user.Email)
	if err != nil {
		app.errorJSON(w, errors.New("no user with this email"), http.StatusBadRequest)
		return
	}

	// Check if the user has this player
	if !contains(userDB.Club, player) {
		app.errorJSON(w, errors.New(player.Name+" is not in your club"), http.StatusBadRequest)
		return
	}

	// Update data base
	dbTimeout := time.Second * 3
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := app.DB.Connection().Collection("real_users")

	userMAJ, err := collection.UpdateOne(
		ctx,
		bson.M{"email": user.Email},
		bson.M{
			"$set":  bson.M{"coins": user.Coins + int64(player.Price)},
			"$pull": bson.M{"club": bson.M{"name": player.Name}},
		},
	)
	if err != nil {
		app.errorJSON(w, errors.New("update coins or add player to error"), http.StatusBadRequest)
		return
	}

	// Augmentation du nombre d'exemplaires disponibles
	collectionAllPlayers := app.DB.Connection().Collection("all_players")
	_, err = collectionAllPlayers.UpdateOne(
		ctx,
		bson.M{"idapi": player.ID_API},
		bson.M{
			"$inc": bson.M{"exemplars": 1},
		},
	)
	if err != nil {
		app.errorJSON(w, errors.New("update exemplar numbers error"), http.StatusBadRequest)
		return
	}

	log.Printf("Player sold: %s", player.Name)

	// Convertir la struct User en un document BSON
	userDoc, err := bson.Marshal(userMAJ)
	if err != nil {
		app.errorJSON(w, errors.New("converting bson error"), http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusAccepted, userDoc) // XXX
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction intermediaire pour verifier qu'un PlayerCard appartient a un tableau de PlayerCards. */
func contains(tab []models.PlayerCard, val models.PlayerCard) bool {
	for _, valCurrent := range tab {
		if valCurrent.ID_API == val.ID_API {
			return true
		}
	}
	return false
}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
func HandleUserPlayers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	err := r.ParseForm()
	log.Printf("User \tip %s\n", r.RemoteAddr)
	if err != nil {
		// handleProblem(w, r)
		panic(err)
	}
	switch r.Method {
	// case "GET":
	// 	controllers.GetUserPlayers(w, r)
	// case "POST":
	// 	user.AddUser(w, r)
	// case "DELETE":
	// 	user.DeleteUser(w, r)
	// default:
	// 	handleProblem(w, r)
	}
}

/*--------------------------------------------------------------------------------------------------------*/
// func HandleUser(w http.ResponseWriter, r *http.Request) {
// 	EnableCors(&w)
// 	err := r.ParseForm()
// 	log.Printf("User \tip %s\n", r.RemoteAddr)
// 	if err != nil {
// 		// handleProblem(w, r)
// 		panic(err)
// 	}
// 	switch r.Method {
// 	case "GET":
// 		controllers.GetUser(w, r)
// 	}
// }

/*--------------------------------------------------------------------------------------------------------*/
/** Handler POST/DELETE permettant de gerer l'ajout/retrait d'un joueur de Team */
func HandleUserTeam(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	err := r.ParseForm()
	log.Printf("User \tip %s\n", r.RemoteAddr)
	if err != nil {
		// handleProblem(w, r)
		panic(err)
	}
	fmt.Println(r.Method)
	switch r.Method {
	case "POST":
		controllers.AddPlayerTeam(w, r)
	case "DELETE":
		controllers.DeletePlayerTeam(w, r)
	}
}

/*--------------------------------------------------------------------------------------------------------*/
/** Handler GET permettant de recuperer les utilisteurs avec le cumul des scores en ordre decroissant */
func HandleUsersRanking(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	err := r.ParseForm()
	log.Printf("User \tip %s\n", r.RemoteAddr)
	if err != nil {
		panic(err)
	}
	switch r.Method {
	case "GET":
		controllers.GetRanking(w, r)
	}
}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
