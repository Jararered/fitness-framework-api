package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"fitness-framework-api/internal/models"
	"fitness-framework-api/internal/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

type API struct {
	DB          *mongo.Database
	VersionInfo *models.ApiInfo
}

func NewAPI(db *mongo.Database, versionInfo *models.ApiInfo) *API {
	return &API{DB: db, VersionInfo: versionInfo}
}

func (api *API) GetExercisesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	equipmentFilters := r.URL.Query()["equipment"]
	musclesFilters := r.URL.Query()["muscles"]

	allExercises, err := mongodb.GetExercises(api.DB)
	if err != nil {
		slog.Error("Error getting all exercises from MongoDB", "error", err)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	equipment, err := mongodb.GetUniqueEquipment(api.DB)
	if err != nil {
		slog.Error("Error getting unique equipment from MongoDB", "error", err)
		http.Error(w, "Failed to fetch equipment options: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipment)
}

func (api *API) GetMusclesOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	muscles, err := mongodb.GetUniqueMuscles(api.DB)
	if err != nil {
		slog.Error("Error getting unique muscles from MongoDB", "error", err)
		http.Error(w, "Failed to fetch muscle options: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(muscles)
}

func (api *API) GetVersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if api.VersionInfo == nil {
		slog.Error("Version information is nil, returning internal server error.")
		http.Error(w, "Version information not available", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.VersionInfo)
}
