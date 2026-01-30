package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"-" bson:"password"`
	UserCart []ProductUser      `json:"usercart" bson:"usercart"`
}

type MongoDb struct {
	Client     mongo.Client
	Collection mongo.Collection
}
