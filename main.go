package main

import (
	"context" // NEW for client disconnect
	"log"
	"net/http"

	"fitness-framework-api/internal/handlers"
	"fitness-framework-api/internal/mongodb" // UPDATED: Import mongodb package
	"fitness-framework-api/internal/version" // Assuming version info is still used
)

const (
	port = ":9001"
)

func main() {
	// 1. Initialize the MongoDB connection
	db, err := mongodb.InitDB() // UPDATED: Call mongodb.InitDB
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// MongoDB client needs to be disconnected gracefully. Get the client from the database.
	defer func() {
		if db != nil {
			err := db.Client().Disconnect(context.Background())
			if err != nil {
				log.Printf("Error disconnecting from MongoDB: %v", err)
			} else {
				log.Println("MongoDB client disconnected.")
			}
		}
	}()

	// 2. Load application version information
	apiInfo, err := version.LoadVersionInfo()
	if err != nil {
		log.Fatalf("Failed to load version information: %v", err)
	}

	// 3. Create an instance of your API handlers, passing dependencies
	apiHandlers := handlers.NewAPI(db, apiInfo) // Pass *mongo.Database

	// 4. Register HTTP handlers
	http.HandleFunc("/api/version", apiHandlers.GetVersionHandler)
	http.HandleFunc("/api/exercises", apiHandlers.GetExercisesHandler)
	http.HandleFunc("/api/equipment-options", apiHandlers.GetEquipmentOptionsHandler)
	http.HandleFunc("/api/muscles-options", apiHandlers.GetMusclesOptionsHandler)

	// 5. Start the HTTP server
	log.Println("Server starting on ", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
