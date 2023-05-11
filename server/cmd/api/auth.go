package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

/** Methodes permettant de gerer les tokens JWT. */
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration // longer that TokenExpiry
	CookieDomain  string        // not readable in JS
	CookiePath    string        // not readable in JS
	CookieName    string        // not readable in JS
}

// Minimal information to issue a token
type jwtUser struct {
	ID     string `json:"id"`
	Pseudo string `json:"Pseudo"`
}

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.RegisteredClaims
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction permettant de generer un token pair. */
func (j *Auth) GenerateTokenPair(user *jwtUser) (TokenPairs, error) {
	log.Printf("GenerateTokenPair")
	// Create a token
	token := jwt.New(jwt.SigningMethodHS256) // most common

	// Set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprint(user.Pseudo)
	claims["sub"] = fmt.Sprint(user.ID) // subject
	// claims["sub"] = fmt.Sprint(strconv.Atoi(user.ID.Hex()))
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix() // Issued at
	claims["typ"] = "JWT"

	log.Printf("==> user.ID: %s", user.ID)
	// Set the expiry for JWT
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// Create a signed token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}
	log.Printf("signedAccessToken ok")

	// Create a refresh token and set claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID) // subject/ID to refresh in the database
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()

	// Set the expiry for the refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

	// Create signed refresh token
	signedRefreshToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	log.Printf("signedRefreshToken ok")

	// Create TokenPairs and populate with signed token
	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// Retrun TokenPairs
	return tokenPairs, nil
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction permettant de creer un refresh cookie. */
func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Domain:   j.CookieDomain,
		Value:    refreshToken,
		Expires:  time.Now().Add(j.RefreshExpiry),
		MaxAge:   int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode, // limited only to this site
		HttpOnly: true,                    //  cookie not readable by JavaScript
		Secure:   true,
	}
}

/*--------------------------------------------------------------------------------------------------------*/
/** Function to delete cookies by setting its max age to -1 and time Unix to 0 */
func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Domain:   j.CookieDomain,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction verifiant le token afin de d'empecher/autoriser l'acces a des routes non autorisees. */
func (j *Auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	w.Header().Add("Vary", "Authorization")

	// Get auth header
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}

	// Split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid auth header")
	}

	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("invalid auth header")
	}

	token := headerParts[1]
	claims := &Claims{}

	// Parsing the token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}

	// Check if the issuer is the right one
	if claims.Issuer != j.Issuer {
		return "", nil, errors.New("invalid issuer")
	}

	return token, claims, nil

}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
