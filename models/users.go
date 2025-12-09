package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user document
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Email     string             `bson:"email" json:"email"`
	Age       int                `bson:"age" json:"age"`
}
