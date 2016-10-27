package entities

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Release struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Timestamp   time.Time     `json:"timestamp" bson:"timestamp"`
	Application string        `json:"application" bson:"application"`
	Version     string        `json:"version" bson:"version"`
}

type Releases []Release
