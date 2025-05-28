package models

// Exercise represents a single workout exercise
type Exercise struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Equipment []string `json:"equipment"`
	Muscles   []string `json:"muscles"`
}
