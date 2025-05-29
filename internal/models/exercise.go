package models

// Exercise represents a single workout exercise
// In the normalized model, this struct represents the data *as returned by the API*,
// which is a "denormalized" view of the underlying database tables.
type Exercise struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Equipment []string `json:"equipment"`
	Muscles   []string `json:"muscles"`

	// Temporary field for unmarshaling the initial JSON input only
	// Will not be persisted directly, but used to populate junction tables.
	InitialJSONMuscles []string `json:"muscles_from_json,omitempty"`
}
