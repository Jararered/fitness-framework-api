package main

import (
	"context"
	"log/slog"
	"net/http"

	"fitness-framework-api/internal/handlers"
	"fitness-framework-api/internal/mongodb"
	"fitness-framework-api/internal/version"
)

const (
	PORT = ":9001"
)

func main() {
	db, err := mongodb.InitDB("workout_app")
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
	}
	defer func() {
		if db != nil {
			err := db.Client().Disconnect(context.Background())
			if err != nil {
				slog.Error("Error disconnecting from MongoDB", "error", err)
			} else {
				slog.Info("MongoDB client disconnected.")
			}
		}
	}()

	apiInfo, err := version.LoadVersionInfo()
	if err != nil {
		slog.Error("Failed to load version information", "error", err)
	}

	apiHandlers := handlers.NewAPI(db, apiInfo)

	http.HandleFunc("/api/version", apiHandlers.GetVersionHandler)
	http.HandleFunc("/api/exercises", apiHandlers.GetExercisesHandler)
	http.HandleFunc("/api/equipment-options", apiHandlers.GetEquipmentOptionsHandler)
	http.HandleFunc("/api/muscles-options", apiHandlers.GetMusclesOptionsHandler)

	slog.Info("Server starting on", "port", PORT)
	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
