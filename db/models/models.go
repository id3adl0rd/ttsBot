package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"firstname,omitempty" validate:"required,alpha"`
}

type File struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Text     string             `json:"text,omitempty" bson:"text,omitempty" validate:"required,alpha"`
	Filepath string             `json:"filepath,omitempty" bson:"filepath,omitempty" validate:"required,alpha"`
}
