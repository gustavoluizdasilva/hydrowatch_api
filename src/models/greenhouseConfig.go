package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type GreenHouseConfig struct {
	ID          primitive.ObjectID `bson:"_id"`
	IdSensor    float64            `json:"idsensor"`
	Greenhouse  string             `json:"greenhouse"`
	Cultivar    string             `json:"cultivar"`
	Temperature Temperature        `json:"temperature"`
	Humidity    Humidity           `json:"humidity"`
	Flowrate    Flowrate           `json:"flowrate"`
}

type Temperature struct {
	Max float64 `json:"max"`
	Min float64 `json:"min"`
}

type Humidity struct {
	Max float64 `json:"max"`
	Min float64 `json:"min"`
}

type Flowrate struct {
	Max float64 `json:"max"`
	Min float64 `json:"min"`
}
