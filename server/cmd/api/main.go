package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"server/cmd/api/controllers"
	"server/cmd/api/database"
	"server/internal/repository"
	"server/internal/repository/dbrepo"
	"time"

	"github.com/rs/cors"
)

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/** Port d'ecoute. */
// const port = 8082

/*--------------------------------------------------------------------------------------------------------*/
type application struct {
	DB           repository.DatabaseRepo
	Domain       string
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

/*--------------------------------------------------------------------------------------------------------*/
func main() {

	// set application config
	var app application

	// read from command line
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", ".example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", ".example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", ".example.com", "domain")
	flag.Parse()

	// connect to the database
	clientDB, db := database.ConnectDB()
	app.DB = &dbrepo.DBRepo{DB: db}
	defer clientDB.Disconnect(context.Background())

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "__Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Starting  application on port", port)

	mux := http.NewServeMux()

	// Get all users from API once
	// controllers.GetAllPlayersAPI()

	// Pull player performance from API weekly and update users score and players rating
	// Execute UserRanking one time per week
	// Go routine to do in parallele of listening handlers
	go func() {
		for {
			log.Println("User Ranking update")
			controllers.UsersRaking()
			// 10080 minutes in one week
			time.Sleep(time.Minute * 10080) // 7 jours
		}
	}()

	mux.HandleFunc("/team", HandleUserTeam)        // POST/DELETE
	mux.HandleFunc("/ranking", HandleUsersRanking) // GET

	mux.HandleFunc("/authenticate", app.authenticate) // POST
	mux.HandleFunc("/refresh", app.refreshToken)      // GET
	mux.HandleFunc("/logout", app.logout)             // GET
	mux.HandleFunc("/register", app.register)         // POST

	mux.HandleFunc("/allPlayers", app.publish_market)          // GET
	mux.HandleFunc("/userInfoCookies", app.user_infos_cookies) // GET
	mux.HandleFunc("/buyingPlayer", app.buy_player)            // POST
	mux.HandleFunc("/sellingPlayer", app.selling_player)       // POST
	mux.HandleFunc("/userInfo", app.user_infos)                // POST

	// start a web server
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://euphonious-kitsune-652034.netlify.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handler))
}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
