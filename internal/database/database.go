package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"fitness-framework-api/internal/models"
)

const (
	DatabaseFileName = "./workout.db"
	JSONFilePath     = "./data/exercises.json"
)

// InitDB initializes the database connection and creates the exercises table.
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", DatabaseFileName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Database connection established.")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS exercises (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		equipment TEXT,
		muscles TEXT
	);`
	_, err = db.Exec(createTableSQL)

	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error creating table: %w", err)
	}
	log.Println("Database table 'exercises' checked/created successfully.")

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM exercises").Scan(&count)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error counting exercises: %w", err)
	}

	if count == 0 {
		log.Println("Database is empty. Attempting to populate initial data from JSON...")
		initialExercises, err := loadExercisesFromJSON(JSONFilePath)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to load initial exercises from JSON: %w", err)
		}

		tx, err := db.Begin()
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback()

		// *** IMPORTANT: Update INSERT statement to match new table and arguments ***
		stmt, err := tx.Prepare(`INSERT INTO exercises (id, name, equipment, muscles) VALUES (?, ?, ?, ?)`)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		for _, ex := range initialExercises {
			ex.ID = uuid.New().String()

			// Marshal slices to JSON strings for database storage
			equipmentJSON, err := json.Marshal(ex.Equipment)
			if err != nil {
				log.Printf("Warning: Could not marshal equipment for %s: %v", ex.Name, err)
				equipmentJSON = []byte("[]")
			}
			musclesJSON, err := json.Marshal(ex.Muscles) // Marshal the new 'Muscles' field
			if err != nil {
				log.Printf("Warning: Could not marshal muscles for %s: %v", ex.Name, err)
				musclesJSON = []byte("[]")
			}

			_, err = stmt.Exec(
				ex.ID,
				ex.Name,
				string(equipmentJSON),
				string(musclesJSON),
			)
			if err != nil {
				db.Close()
				return nil, fmt.Errorf("error inserting initial exercise %s: %w", ex.Name, err)
			}
		}

		err = tx.Commit()
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
		log.Printf("Successfully populated %d exercises from JSON.", len(initialExercises))
	} else {
		log.Printf("Database already contains %d exercises. Skipping initial data population.", count)
	}

	return db, nil
}

func loadExercisesFromJSON(filePath string) ([]models.Exercise, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for %s: %w", filePath, err)
	}

	log.Printf("Attempting to load exercises from: %s", absPath)

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file %s: %w", filePath, err)
	}

	var exercises []models.Exercise
	err = json.Unmarshal(data, &exercises)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data from %s: %w", filePath, err)
	}

	return exercises, nil
}

// GetExercises fetches all exercises from the database.
func GetExercises(db *sql.DB) ([]models.Exercise, error) {
	rows, err := db.Query("SELECT id, name, equipment, muscles FROM exercises")
	if err != nil {
		return nil, fmt.Errorf("failed to query exercises: %w", err)
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var ex models.Exercise
		var equipmentJSON string
		var musclesJSON string

		err := rows.Scan(&ex.ID, &ex.Name, &equipmentJSON, &musclesJSON)
		if err != nil {
			log.Printf("Error scanning exercise row: %v", err)
			continue
		}

		// Unmarshal JSON strings back into Go string slices
		err = json.Unmarshal([]byte(equipmentJSON), &ex.Equipment)
		if err != nil {
			log.Printf("Warning: Could not unmarshal equipment for %s (ID: %s): %v", ex.Name, ex.ID, err)
			ex.Equipment = []string{}
		}
		// Unmarshal musclesJSON into ex.Muscles
		err = json.Unmarshal([]byte(musclesJSON), &ex.Muscles)
		if err != nil {
			log.Printf("Warning: Could not unmarshal muscles for %s (ID: %s): %v", ex.Name, ex.ID, err)
			ex.Muscles = []string{}
		}

		exercises = append(exercises, ex)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during exercise iteration: %w", err)
	}

	return exercises, nil
}

// GetUniqueExerciseNames fetches all distinct exercise names from the database.
func GetUniqueExerciseNames(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT name FROM exercises ORDER BY name ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to query distinct exercise names: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Printf("Error scanning exercise name row: %v", err)
			continue
		}
		names = append(names, name)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during exercise name iteration: %w", err)
	}
	return names, nil
}

// GetUniqueMuscles fetches all unique muscle components from the database.
func GetUniqueMuscles(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT muscles FROM exercises")
	if err != nil {
		return nil, fmt.Errorf("failed to query muscles: %w", err)
	}
	defer rows.Close()

	uniqueMusclesMap := make(map[string]bool)
	for rows.Next() {
		var musclesJSON string
		if err := rows.Scan(&musclesJSON); err != nil {
			log.Printf("Error scanning muscles JSON row: %v", err)
			continue
		}

		var musclesList []string
		if err := json.Unmarshal([]byte(musclesJSON), &musclesList); err != nil {
			log.Printf("Warning: Could not unmarshal muscles JSON '%s': %v", musclesJSON, err)
			continue
		}

		for _, item := range musclesList {
			uniqueMusclesMap[item] = true
		}
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during muscles iteration: %w", err)
	}

	var uniqueMuscles []string
	for item := range uniqueMusclesMap {
		uniqueMuscles = append(uniqueMuscles, item)
	}
	sort.Strings(uniqueMuscles)

	return uniqueMuscles, nil
}

// GetUniqueEquipment fetches all unique equipment types from the database (remains unchanged).
func GetUniqueEquipment(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT equipment FROM exercises")
	if err != nil {
		return nil, fmt.Errorf("failed to query equipment: %w", err)
	}
	defer rows.Close()

	uniqueEquipmentMap := make(map[string]bool)
	for rows.Next() {
		var equipmentJSON string
		if err := rows.Scan(&equipmentJSON); err != nil {
			log.Printf("Error scanning equipment JSON row: %v", err)
			continue
		}

		var equipmentList []string
		if err := json.Unmarshal([]byte(equipmentJSON), &equipmentList); err != nil {
			log.Printf("Warning: Could not unmarshal equipment JSON '%s': %v", equipmentJSON, err)
			continue
		}

		for _, item := range equipmentList {
			uniqueEquipmentMap[item] = true
		}
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during equipment iteration: %w", err)
	}

	var uniqueEquipment []string
	for item := range uniqueEquipmentMap {
		uniqueEquipment = append(uniqueEquipment, item)
	}
	sort.Strings(uniqueEquipment)

	return uniqueEquipment, nil
}
