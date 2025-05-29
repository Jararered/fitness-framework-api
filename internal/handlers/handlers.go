package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"fitness-framework-api/internal/models"
	"fitness-framework-api/internal/mongodb" // UPDATED: Import mongodb package

	"go.mongodb.org/mongo-driver/mongo" // NEW: Import mongo driver
)

// API holds dependencies for handlers, e.g., database connection and version info
type API struct {
	DB          *mongo.Database // UPDATED: Now holds *mongo.Database
	VersionInfo *models.ApiInfo // Assuming you're still using models.ApiInfo from the previous step
}

// NewAPI creates a new API instance with the given database connection and version info.
func NewAPI(db *mongo.Database, versionInfo *models.ApiInfo) *API { // UPDATED: Parameter type
	return &API{DB: db, VersionInfo: versionInfo}
}

// GetExercisesHandler (Logic remains the same, but calls mongodb.GetExercises)
func (api *API) GetExercisesHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	equipmentFilters := r.URL.Query()["equipment"]
	musclesFilters := r.URL.Query()["muscles"]

	allExercises, err := mongodb.GetExercises(api.DB) // UPDATED: Call mongodb package
	if err != nil {
		log.Printf("Error getting all exercises from MongoDB: %v", err)
		http.Error(w, "Failed to fetch exercises: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var filteredExercises []models.Exercise

	for _, ex := range allExercises {
		equipmentMatches := containsAnyCaseInsensitive(ex.Equipment, equipmentFilters)
		musclesMatches := containsAnyCaseInsensitive(ex.Muscles, musclesFilters)

		if equipmentMatches && musclesMatches {
			filteredExercises = append(filteredExercises, ex)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredExercises)
}

// containsAnyCaseInsensitive (no changes, still useful for in-memory filtering)
func containsAnyCaseInsensitive(haystack []string, needles []string) bool {
	if len(needles) == 0 {
		return true
	}
	for _, needle := range needles {
		for _, item := range haystack {
			if strings.EqualFold(item, needle) {
				return true
			}
		}
	}
	return false
}

// GetEquipmentOptionsHandler (calls mongodb.GetUniqueEquipment)
func (api *API) GetEquipmentOptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	equipment, err := mongodb.GetUniqueEquipment(api.DB) // UPDATED: Call mongodb package
	if err != nil {
		log.Printf("Error getting unique equipment from MongoDB: %v", err)
		http.Error(w, "Failed to fetch equipment options: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipment)
}

// GetMusclesOptionsHandler (calls mongodb.GetUniqueMuscles)
func (api *API) GetMusclesOptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	muscles, err := mongodb.GetUniqueMuscles(api.DB) // UPDATED: Call mongodb package
	if err != nil {
		log.Printf("Error getting unique muscles from MongoDB: %v", err)
		http.Error(w, "Failed to fetch muscle options: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(muscles)
}

// GetVersionHandler (no changes)
func (api *API) GetVersionHandler(w http.ResponseWriter, r *http.Request) {
	// ... (content remains the same as your current GetVersionHandler)
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if api.VersionInfo == nil {
		log.Println("Version information is nil, returning internal server error.")
		http.Error(w, "Version information not available", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.VersionInfo)
}
