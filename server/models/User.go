package models

import "gopkg.in/mgo.v2/bson"

//User Struct
type USER struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Name       string        `json:"name"`
	Email      string        `json:"email"`
	IsApproved bool          `json:"isapproved"`
}
