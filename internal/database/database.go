package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"fitness-framework-api/internal/models"
)

const (
	DatabaseFileName = "./workout.db"
	JSONFilePath     = "./data/exercises.json"
)

// InitDB initializes the database connection and creates the exercises table.
// It returns the *sql.DB connection.
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", DatabaseFileName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Database connection established.")

	// *** NEW: Create all 5 tables ***
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS exercises (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE
	);
	CREATE TABLE IF NOT EXISTS equipment (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE
	);
	CREATE TABLE IF NOT EXISTS muscles (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE
	);
	CREATE TABLE IF NOT EXISTS exercise_equipment (
		exercise_id TEXT NOT NULL,
		equipment_id TEXT NOT NULL,
		PRIMARY KEY (exercise_id, equipment_id),
		FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE,
		FOREIGN KEY (equipment_id) REFERENCES equipment(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS exercise_muscles (
		exercise_id TEXT NOT NULL,
		muscle_id TEXT NOT NULL,
		PRIMARY KEY (exercise_id, muscle_id),
		FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE,
		FOREIGN KEY (muscle_id) REFERENCES muscles(id) ON DELETE CASCADE
	);`
	_, err = db.Exec(createTablesSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error creating tables: %w", err)
	}
	log.Println("Database tables checked/created successfully.")

	// Check if exercises table is empty (proxy for checking if data needs to be populated)
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

		tx, err := db.Begin() // Start a transaction for bulk insertion
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback() // Rollback on error, Commit later on success

		// --- Collect and insert unique equipment and muscles first ---
		equipmentNameToID := make(map[string]string)
		muscleNameToID := make(map[string]string)

		// Collect unique equipment
		for _, ex := range initialExercises {
			for _, eq := range ex.Equipment {
				if _, exists := equipmentNameToID[eq]; !exists {
					equipmentNameToID[eq] = uuid.New().String()
				}
			}
			// Collect unique muscles (using the temporary InitialJSONMuscles field)
			for _, m := range ex.InitialJSONMuscles {
				if _, exists := muscleNameToID[m]; !exists {
					muscleNameToID[m] = uuid.New().String()
				}
			}
		}

		// Insert unique equipment
		stmtInsertEquipment, err := tx.Prepare(`INSERT INTO equipment (id, name) VALUES (?, ?)`)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to prepare equipment insert statement: %w", err)
		}
		defer stmtInsertEquipment.Close()
		for name, id := range equipmentNameToID {
			_, err = stmtInsertEquipment.Exec(id, name)
			if err != nil {
				db.Close()
				return nil, fmt.Errorf("error inserting equipment '%s': %w", name, err)
			}
		}
		log.Printf("Inserted %d unique equipment items.", len(equipmentNameToID))

		// Insert unique muscles
		stmtInsertMuscles, err := tx.Prepare(`INSERT INTO muscles (id, name) VALUES (?, ?)`)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to prepare muscles insert statement: %w", err)
		}
		defer stmtInsertMuscles.Close()
		for name, id := range muscleNameToID {
			_, err = stmtInsertMuscles.Exec(id, name)
			if err != nil {
				db.Close()
				return nil, fmt.Errorf("error inserting muscle '%s': %w", name, err)
			}
		}
		log.Printf("Inserted %d unique muscle items.", len(muscleNameToID))

		// --- Insert exercises and their associations ---
		stmtInsertExercise, err := tx.Prepare(`INSERT INTO exercises (id, name) VALUES (?, ?)`)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to prepare exercise insert statement: %w", err)
		}
		defer stmtInsertExercise.Close()

		stmtInsertExerciseEquipment, err := tx.Prepare(`INSERT INTO exercise_equipment (exercise_id, equipment_id) VALUES (?, ?)`)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to prepare exercise_equipment insert statement: %w", err)
		}
		defer stmtInsertExerciseEquipment.Close()

		stmtInsertExerciseMuscles, err := tx.Prepare(`INSERT INTO exercise_muscles (exercise_id, muscle_id) VALUES (?, ?)`)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to prepare exercise_muscles insert statement: %w", err)
		}
		defer stmtInsertExerciseMuscles.Close()

		for _, ex := range initialExercises {
			ex.ID = uuid.New().String() // Generate UUID for each exercise

			_, err = stmtInsertExercise.Exec(ex.ID, ex.Name)
			if err != nil {
				db.Close()
				return nil, fmt.Errorf("error inserting exercise '%s': %w", ex.Name, err)
			}

			// Insert into exercise_equipment junction table
			for _, eqName := range ex.Equipment {
				eqID, ok := equipmentNameToID[eqName]
				if !ok {
					log.Printf("Warning: Equipment '%s' not found for exercise '%s'. Skipping association.", eqName, ex.Name)
					continue
				}
				_, err = stmtInsertExerciseEquipment.Exec(ex.ID, eqID)
				if err != nil {
					db.Close()
					return nil, fmt.Errorf("error inserting exercise_equipment for '%s' and '%s': %w", ex.Name, eqName, err)
				}
			}

			// Insert into exercise_muscles junction table (using InitialJSONMuscles)
			for _, muscleName := range ex.InitialJSONMuscles {
				muscleID, ok := muscleNameToID[muscleName]
				if !ok {
					log.Printf("Warning: Muscle '%s' not found for exercise '%s'. Skipping association.", muscleName, ex.Name)
					continue
				}
				_, err = stmtInsertExerciseMuscles.Exec(ex.ID, muscleID)
				if err != nil {
					db.Close()
					return nil, fmt.Errorf("error inserting exercise_muscles for '%s' and '%s': %w", ex.Name, muscleName, err)
				}
			}
		}

		err = tx.Commit() // Commit the transaction
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
		log.Printf("Successfully populated %d exercises and their associations.", len(initialExercises))
	} else {
		log.Printf("Database already contains %d exercises. Skipping initial data population.", count)
	}

	return db, nil // Return the opened database connection
}

// loadExercisesFromJSON now specifically unmarshals into models.Exercise,
// where the "muscles" key from JSON goes into InitialJSONMuscles.
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

	// Define a temporary struct that matches the exact input JSON structure
	type TempExercise struct {
		Name               string   `json:"name"`
		Equipment          []string `json:"equipment"`
		InitialJSONMuscles []string `json:"muscles"` // This matches the "muscles" key in your input JSON
	}

	var tempExercises []TempExercise
	err = json.Unmarshal(data, &tempExercises)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data from %s: %w", filePath, err)
	}

	// Convert TempExercise to models.Exercise
	var exercises []models.Exercise
	for _, te := range tempExercises {
		exercises = append(exercises, models.Exercise{
			Name:               te.Name,
			Equipment:          te.Equipment,
			InitialJSONMuscles: te.InitialJSONMuscles, // Store for processing in InitDB
		})
	}

	return exercises, nil
}

// GetExercises fetches all exercises from the database,
// including their associated equipment and muscles via JOINs.
func GetExercises(db *sql.DB) ([]models.Exercise, error) {
	// Query to get all exercises along with their associated equipment and muscles.
	// We use GROUP_CONCAT to aggregate the names into a single comma-separated string,
	// which is then split in Go. This is a common pattern for SQLite.
	query := `
	SELECT
		e.id,
		e.name,
		GROUP_CONCAT(DISTINCT eq.name) AS equipment_names,
		GROUP_CONCAT(DISTINCT m.name) AS muscle_names
	FROM
		exercises e
	LEFT JOIN
		exercise_equipment ee ON e.id = ee.exercise_id
	LEFT JOIN
		equipment eq ON ee.equipment_id = eq.id
	LEFT JOIN
		exercise_muscles em ON e.id = em.exercise_id
	LEFT JOIN
		muscles m ON em.muscle_id = m.id
	GROUP BY
		e.id, e.name
	ORDER BY
		e.name ASC;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query exercises: %w", err)
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var ex models.Exercise
		var equipmentNamesSQL, muscleNamesSQL sql.NullString // Use sql.NullString for columns that might be NULL due to LEFT JOIN

		err := rows.Scan(&ex.ID, &ex.Name, &equipmentNamesSQL, &muscleNamesSQL)
		if err != nil {
			log.Printf("Error scanning exercise row: %v", err)
			continue
		}

		// Split the comma-separated strings back into slices
		if equipmentNamesSQL.Valid {
			ex.Equipment = strings.Split(equipmentNamesSQL.String, ",")
		} else {
			ex.Equipment = []string{}
		}
		if muscleNamesSQL.Valid {
			ex.Muscles = strings.Split(muscleNamesSQL.String, ",")
		} else {
			ex.Muscles = []string{}
		}

		// Clean up empty strings that might result from splitting empty or null GROUP_CONCAT results
		// Example: if an exercise has no equipment, GROUP_CONCAT returns NULL, which becomes "" after NullString.String.
		// strings.Split("", ",") results in [""]
		if len(ex.Equipment) == 1 && ex.Equipment[0] == "" {
			ex.Equipment = []string{}
		}
		if len(ex.Muscles) == 1 && ex.Muscles[0] == "" {
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

// GetUniqueMuscles fetches all unique muscle names from the dedicated 'muscles' table.
func GetUniqueMuscles(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM muscles ORDER BY name ASC") // Simpler query now
	if err != nil {
		return nil, fmt.Errorf("failed to query unique muscles: %w", err)
	}
	defer rows.Close()

	var muscles []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Printf("Error scanning muscle name row: %v", err)
			continue
		}
		muscles = append(muscles, name)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during unique muscles iteration: %w", err)
	}
	return muscles, nil
}

// GetUniqueEquipment fetches all unique equipment names from the dedicated 'equipment' table.
func GetUniqueEquipment(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM equipment ORDER BY name ASC") // Simpler query now
	if err != nil {
		return nil, fmt.Errorf("failed to query unique equipment: %w", err)
	}
	defer rows.Close()

	var equipment []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Printf("Error scanning equipment name row: %v", err)
			continue
		}
		equipment = append(equipment, name)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during unique equipment iteration: %w", err)
	}
	return equipment, nil
}
