package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type GreenhouseData struct {
	ID           primitive.ObjectID `bson:"_id"`
	IdGreenhouse float64            `json:"idgreenhouse"`
	Temperature  float64            `json:"temperature"`
	Humidity     float64            `json:"humidity"`
	Flowrate     float64            `json:"flowrate"`
	DayMonth     primitive.DateTime `json:"timestamp"`
}
