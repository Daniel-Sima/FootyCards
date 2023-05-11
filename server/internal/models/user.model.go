package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
type User struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"password" bson:"password"` // hash of a password actually
	Pseudo     string             `json:"pseudo" bson:"pseudo" `
	ClubName   string             `json:"clubName" bson:"clubName"`
	Coins      int64              `json:"coins" bson:"coins"`
	Score      float32            `json:"score" bson:"score"`
	Classement int32              `json:"classement" bson:"classement"`
	Club       []PlayerCard       `json:"club" bson:"club,omitempty" `
	Team       []PlayerCard       `json:"team" bson:"team,omitempty" `
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	LastUpdate time.Time          `json:"LastUpdate" bson:"LastUpdate"`
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction qui verifie les mot de passes entre celui fourni et celui de la base de donnees. */
func (u *User) PasswordMatches(plainText string) (bool, error) {
	// Compares the hash password from the database 'plainText'

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction permettant de genere une hash d'un mot de passe (pour les tests).*/
func (u *User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
