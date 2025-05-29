package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"fitness-framework-api/internal/database"
	"fitness-framework-api/internal/models"
)

const (
	allowedOrigin = "http://ff.jarare.red"
)

// API holds dependencies for handlers, e.g., database connection and version info
type API struct {
	DB      *sql.DB
	ApiInfo *models.ApiInfo
}

// NewAPI creates a new API instance with the given database connection and version info.
func NewAPI(db *sql.DB, apiInfo *models.ApiInfo) *API {
	return &API{DB: db, ApiInfo: apiInfo}
}

func (api *API) GetExercisesHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	equipmentFilters := r.URL.Query()["equipment"]
	musclesFilters := r.URL.Query()["muscles"]

	allExercises, err := database.GetExercises(api.DB)
	if err != nil {
		log.Printf("Error getting all exercises: %v", err)
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

func (api *API) GetEquipmentOptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	equipment, err := database.GetUniqueEquipment(api.DB)
	if err != nil {
		log.Printf("Error getting unique equipment: %v", err)
		http.Error(w, "Failed to fetch equipment options: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipment)
}

func (api *API) GetMusclesOptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	muscles, err := database.GetUniqueMuscles(api.DB)
	if err != nil {
		log.Printf("Error getting unique muscles: %v", err)
		http.Error(w, "Failed to fetch muscle options: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(muscles)
}

func (api *API) GetVersionHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if api.ApiInfo == nil {
		log.Println("Version information is nil, returning internal server error.")
		http.Error(w, "Version information not available", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.ApiInfo)
}
