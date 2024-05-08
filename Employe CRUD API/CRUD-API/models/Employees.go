package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Employee struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Mobile   string             `json:"mobile" bson:"mobile"`
	Company  string             `json:"company" bson:"company"`
	Location string             `json:"location" bson:"location"`
}
