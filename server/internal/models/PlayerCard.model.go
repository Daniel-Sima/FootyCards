package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
type PlayerCard struct {
	ID           primitive.ObjectID `json:"_id"  bson:"_id,omitempty"`
	ID_API       int32              `json:"idapi"  bson:"idapi,omitempty"`
	Name         string             `json:"name"  bson:"name,omitempty"`
	RatingGlobal float32            `json:"ratingGlobal" bson:"ratingGlobal,omitempty"`
	LastRating   float32            `json:"lastRating" bson:"lastRating,omitempty"`
	Image        string             `json:"image" bson:"image,omitempty"`
	Price        int32              `json:"price" bson:"price,omitempty"`
	Club         string             `json:"club" bson:"club,omitempty"`
	InTeam       bool               `json:"inteam" bson:"inteam,omitempty"`
	Position     string             `json:"position" bson:"position,omitempty"`
	Exemplars    int32              `json:"exemplars" bson:"exemplars,omitempty"`
	LastUpdate   time.Time          `json:"lastupdate" bson:"lastupdate,omitempty"`
}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
