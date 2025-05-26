// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"       // For formatted errors/strings
	"io/ioutil" // For reading files (deprecated in newer Go, but simpler for this use case)
	"log"
	"net/http"
	"path/filepath" // For building file paths safely

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// Exercise represents a single workout exercise
type Exercise struct {
	ID                string   `json:"id"`                // Unique identifier
	Name              string   `json:"name"`              // e.g., "Squats"
	Equipment         []string `json:"equipment"`         // e.g., ["Barbell", "Squat Rack"]
	Difficulty        string   `json:"difficulty"`        // e.g., "Beginner", "Intermediate", "Advanced"
	MusclesWorked     []string `json:"musclesWorked"`     // e.g., ["Quadriceps", "Glutes", "Hamstrings"]
	MainBodyComponent string   `json:"mainBodyComponent"` // e.g., "Legs", "Chest", "Back"
}

// Global database connection variable
var db *sql.DB

const (
	databaseFileName = "./workout.db"
	jsonFilePath     = "./data/exercises.json" // Path to your JSON file
)

func main() {
	// Initialize the database connection and schema
	initDB()
	// Ensure the database connection is closed when the main function exits
	defer db.Close()

	// Register the handler for the /exercises endpoint
	http.HandleFunc("/exercises", getExercisesHandler)

	// Start the server on port 8080
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// initDB initializes the database connection, creates the exercises table if it doesn't exist,
// and populates initial data from JSON if the table is empty.
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", databaseFileName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Database connection established.")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS exercises (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		equipment TEXT,        -- Stored as JSON string
		difficulty TEXT,
		muscles_worked TEXT,   -- Stored as JSON string
		main_body_component TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	log.Println("Database table 'exercises' checked/created successfully.")

	// Check if the table is empty and populate with initial data from JSON
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM exercises").Scan(&count)
	if err != nil {
		log.Fatalf("Error counting exercises: %v", err)
	}

	if count == 0 {
		log.Println("Database is empty. Attempting to populate initial data from JSON...")
		initialExercises, err := loadExercisesFromJSON(jsonFilePath)
		if err != nil {
			log.Fatalf("Failed to load initial exercises from JSON: %v", err)
		}

		tx, err := db.Begin() // Start a transaction for bulk insertion
		if err != nil {
			log.Fatalf("Failed to begin transaction: %v", err)
		}
		defer tx.Rollback() // Rollback on error, Commit later on success

		stmt, err := tx.Prepare(`INSERT INTO exercises (id, name, equipment, difficulty, muscles_worked, main_body_component) VALUES (?, ?, ?, ?, ?, ?)`)
		if err != nil {
			log.Fatalf("Failed to prepare statement: %v", err)
		}
		defer stmt.Close()

		for _, ex := range initialExercises {
			// Generate a new UUID for each exercise
			ex.ID = uuid.New().String()

			// Marshal slices (Equipment, MusclesWorked) into JSON strings for database storage
			equipmentJSON, err := json.Marshal(ex.Equipment)
			if err != nil {
				log.Printf("Warning: Could not marshal equipment for %s: %v", ex.Name, err)
				equipmentJSON = []byte("[]")
			}
			musclesJSON, err := json.Marshal(ex.MusclesWorked)
			if err != nil {
				log.Printf("Warning: Could not marshal muscles for %s: %v", ex.Name, err)
				musclesJSON = []byte("[]")
			}

			_, err = stmt.Exec(ex.ID, ex.Name, string(equipmentJSON), ex.Difficulty, string(musclesJSON), ex.MainBodyComponent)
			if err != nil {
				log.Fatalf("Error inserting initial exercise %s: %v", ex.Name, err)
			}
		}

		err = tx.Commit() // Commit the transaction
		if err != nil {
			log.Fatalf("Failed to commit transaction: %v", err)
		}
		log.Printf("Successfully populated %d exercises from JSON.", len(initialExercises))
	} else {
		log.Printf("Database already contains %d exercises. Skipping initial data population.", count)
	}
}

// loadExercisesFromJSON reads exercise data from a JSON file.
func loadExercisesFromJSON(filePath string) ([]Exercise, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for %s: %w", filePath, err)
	}

	log.Printf("Attempting to load exercises from: %s", absPath)

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file %s: %w", filePath, err)
	}

	var exercises []Exercise
	err = json.Unmarshal(data, &exercises)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data from %s: %w", filePath, err)
	}

	return exercises, nil
}

// getExercisesHandler fetches all exercises from the database and returns them as JSON
func getExercisesHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	rows, err := db.Query("SELECT id, name, equipment, difficulty, muscles_worked, main_body_component FROM exercises")
	if err != nil {
		http.Error(w, "Failed to fetch exercises: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var exercises []Exercise
	for rows.Next() {
		var ex Exercise
		var equipmentJSON, musclesJSON string

		err := rows.Scan(&ex.ID, &ex.Name, &equipmentJSON, &ex.Difficulty, &musclesJSON, &ex.MainBodyComponent)
		if err != nil {
			log.Printf("Error scanning exercise row: %v", err)
			continue
		}

		err = json.Unmarshal([]byte(equipmentJSON), &ex.Equipment)
		if err != nil {
			log.Printf("Warning: Could not unmarshal equipment for %s (ID: %s): %v", ex.Name, ex.ID, err)
			ex.Equipment = []string{}
		}
		err = json.Unmarshal([]byte(musclesJSON), &ex.MusclesWorked)
		if err != nil {
			log.Printf("Warning: Could not unmarshal muscles for %s (ID: %s): %v", ex.Name, ex.ID, err)
			ex.MusclesWorked = []string{}
		}

		exercises = append(exercises, ex)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Error during exercise iteration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercises)
}
