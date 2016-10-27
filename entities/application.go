package entities

import "gopkg.in/mgo.v2/bson"

type Application struct {
	Id   bson.ObjectId `json:"id" bson:"_id"`
	Name string        `json:"name" bson:"name"`
}

type Applications []Application
