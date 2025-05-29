package models

import "go.mongodb.org/mongo-driver/bson/primitive" // NEW import

// Exercise represents a single workout exercise as it's returned by the API
// and how it's stored/assembled from MongoDB.
type Exercise struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"` // MongoDB uses _id. omitempty for _id in API output.
	Name      string             `json:"name" bson:"name"`
	Equipment []string           `json:"equipment" bson:"equipment"` // MongoDB stores arrays directly
	Muscles   []string           `json:"muscles" bson:"muscles"`     // MongoDB stores arrays directly
}

// RawExercise represents an exercise as it's defined directly in Go code (from internal/data).
// This struct is primarily for unmarshaling the initial data.
// It matches the shape of data/exercises.json or internal/data/exercises.go.
type RawExercise struct {
	Name      string   `json:"name"`
	Equipment []string `json:"equipment"`
	Muscles   []string `json:"muscles"`
}
