package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Exercise struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Equipment []string           `json:"equipment" bson:"equipment"`
	Muscles   []string           `json:"muscles" bson:"muscles"`
}

type RawExercise struct {
	Name      string   `json:"name"`
	Equipment []string `json:"equipment"`
	Muscles   []string `json:"muscles"`
}
