package main

import (
	"log"
	"net/http"

	"fitness-framework-api/internal/database"
	"fitness-framework-api/internal/handlers"
	"fitness-framework-api/internal/version"
)

const (
	port = ":9001"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	apiInfo, err := version.LoadVersionInfo()
	if err != nil {
		log.Fatalf("Failed to load version information: %v", err)
	}

	// 3. Create an instance of your API handlers, passing dependencies
	apiHandlers := handlers.NewAPI(db, apiInfo)

	http.HandleFunc("/api/version", apiHandlers.GetVersionHandler)

	http.HandleFunc("/api/exercises", apiHandlers.GetExercisesHandler)
	http.HandleFunc("/api/equipment-options", apiHandlers.GetEquipmentOptionsHandler)
	http.HandleFunc("/api/muscles-options", apiHandlers.GetMusclesOptionsHandler)

	log.Println("Server starting on ", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
