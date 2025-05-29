package mongodb

import (
	"context" // For context.Context
	"fmt"
	"log"
	"sort"
	"time" // For context timeout

	"go.mongodb.org/mongo-driver/bson"           // For BSON marshalling/unmarshalling and queries
	"go.mongodb.org/mongo-driver/bson/primitive" // For ObjectID
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" // For client options

	"fitness-framework-api/internal/constants" // Import constants package
	"fitness-framework-api/internal/data"      // Import data package
	"fitness-framework-api/internal/models"
)

const (
	MongoURI       = "mongodb://localhost:27017" // MongoDB connection string
	DatabaseName   = "workout_app"
	CollectionName = "exercises"
)

// InitDB initializes the MongoDB connection and populates initial data
// if the exercises collection is empty.
// It returns the *mongo.Database connection.
func InitDB() (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(MongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	// Ping the primary to verify connection
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		client.Disconnect(context.Background()) // Disconnect if ping fails
		return nil, fmt.Errorf("error pinging MongoDB: %w", err)
	}

	log.Println("MongoDB connection established.")
	db := client.Database(DatabaseName)
	collection := db.Collection(CollectionName)

	// Check if exercises collection is empty
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, bson.D{}) // Empty filter counts all documents
	if err != nil {
		client.Disconnect(context.Background())
		return nil, fmt.Errorf("error counting documents in exercises collection: %w", err)
	}

	if count == 0 {
		log.Println("Exercises collection is empty. Attempting to populate initial data from hardcoded Go data...")

		// --- Populate unique equipment and muscles from constants ---
		// We'll insert these into their own collections for easier querying later,
		// though MongoDB's distinct() could get them directly from exercises collection too.
		// This keeps your unique lists canonical.
		equipmentCollection := db.Collection("equipment_options")
		muscleCollection := db.Collection("muscles_options")

		// Drop existing option collections to ensure freshness
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := equipmentCollection.Drop(ctx); err != nil {
			log.Printf("Warning: Could not drop equipment_options collection: %v", err)
		}
		if err := muscleCollection.Drop(ctx); err != nil {
			log.Printf("Warning: Could not drop muscles_options collection: %v", err)
		}

		// Prepare documents for bulk insert
		var equipmentDocs []interface{}
		for _, name := range constants.AllEquipmentNames {
			equipmentDocs = append(equipmentDocs, bson.D{{"_id", primitive.NewObjectID()}, {"name", name}})
		}
		var muscleDocs []interface{}
		for _, name := range constants.AllMuscleGroupNames {
			muscleDocs = append(muscleDocs, bson.D{{"_id", primitive.NewObjectID()}, {"name", name}})
		}

		if len(equipmentDocs) > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if _, err := equipmentCollection.InsertMany(ctx, equipmentDocs); err != nil {
				client.Disconnect(context.Background())
				return nil, fmt.Errorf("error inserting equipment options: %w", err)
			}
			log.Printf("Inserted %d unique equipment items from constants.", len(equipmentDocs))
		}
		if len(muscleDocs) > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if _, err := muscleCollection.InsertMany(ctx, muscleDocs); err != nil {
				client.Disconnect(context.Background())
				return nil, fmt.Errorf("error inserting muscle options: %w", err)
			}
			log.Printf("Inserted %d unique muscle items from constants.", len(muscleDocs))
		}

		// --- Insert exercises ---
		var documents []interface{}
		for _, rawEx := range data.AllRawExercises {
			// Validate against constants before adding to documents, like before
			validEquipment := []string{}
			for _, eqName := range rawEx.Equipment {
				if constants.IsValidEquipment(eqName) {
					validEquipment = append(validEquipment, eqName)
				} else {
					log.Printf("Warning: Invalid equipment '%s' for exercise '%s'. Skipping.", eqName, rawEx.Name)
				}
			}
			validMuscles := []string{}
			for _, muscleName := range rawEx.Muscles {
				if constants.IsValidMuscleGroup(muscleName) {
					validMuscles = append(validMuscles, muscleName)
				} else {
					log.Printf("Warning: Invalid muscle group '%s' for exercise '%s'. Skipping.", muscleName, rawEx.Name)
				}
			}

			// Create a models.Exercise struct and then convert it to BSON document
			// MongoDB automatically handles primitive.ObjectID for _id
			doc := models.Exercise{
				ID:        primitive.NewObjectID(), // MongoDB generates _id automatically if not provided, but explicit is good
				Name:      rawEx.Name,
				Equipment: validEquipment,
				Muscles:   validMuscles,
			}
			documents = append(documents, doc)
		}

		if len(documents) > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second) // Longer timeout for bulk insert
			defer cancel()
			_, err := collection.InsertMany(ctx, documents)
			if err != nil {
				client.Disconnect(context.Background())
				return nil, fmt.Errorf("error inserting initial exercises: %w", err)
			}
			log.Printf("Successfully populated %d exercises.", len(documents))
		}

	} else {
		log.Printf("Exercises collection already contains %d documents. Skipping initial data population.", count)
	}

	return db, nil // Return the MongoDB database connection
}

// GetExercises fetches all exercises from the MongoDB collection.
func GetExercises(db *mongo.Database) ([]models.Exercise, error) {
	collection := db.Collection(CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{}) // Empty filter to find all documents
	if err != nil {
		return nil, fmt.Errorf("failed to find exercises: %w", err)
	}
	defer cursor.Close(ctx)

	var exercises []models.Exercise
	if err = cursor.All(ctx, &exercises); err != nil { // Decode all documents into the slice
		return nil, fmt.Errorf("failed to decode exercises: %w", err)
	}

	// For consistency with SQL GROUP_CONCAT, sort slices within each exercise
	// (MongoDB doesn't guarantee order for arrays unless you explicitly sort them on insert/query)
	for i := range exercises {
		sort.Strings(exercises[i].Equipment)
		sort.Strings(exercises[i].Muscles)
	}

	return exercises, nil
}

// GetUniqueExerciseNames fetches all distinct exercise names from the database.
func GetUniqueExerciseNames(db *mongo.Database) ([]string, error) {
	collection := db.Collection(CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use Distinct method to get unique names
	distinctNames, err := collection.Distinct(ctx, "name", bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to get distinct exercise names: %w", err)
	}

	var names []string
	for _, doc := range distinctNames {
		if name, ok := doc.(string); ok {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names, nil
}

// GetUniqueMuscles fetches all unique muscle names from the dedicated 'muscles_options' collection.
func GetUniqueMuscles(db *mongo.Database) ([]string, error) {
	collection := db.Collection("muscles_options") // Query dedicated options collection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{"name", 1}}))
	if err != nil {
		return nil, fmt.Errorf("failed to find unique muscles: %w", err)
	}
	defer cursor.Close(ctx)

	var muscles []string
	for cursor.Next(ctx) {
		var result struct {
			Name string `bson:"name"`
		}
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Warning: Failed to decode muscle option: %v", err)
			continue
		}
		muscles = append(muscles, result.Name)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error during unique muscles iteration: %w", err)
	}
	// No need to sort if SetSort is used in query.
	return muscles, nil
}

// GetUniqueEquipment fetches all unique equipment names from the dedicated 'equipment_options' collection.
func GetUniqueEquipment(db *mongo.Database) ([]string, error) {
	collection := db.Collection("equipment_options") // Query dedicated options collection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{"name", 1}}))
	if err != nil {
		return nil, fmt.Errorf("failed to find unique equipment: %w", err)
	}
	defer cursor.Close(ctx)

	var equipment []string
	for cursor.Next(ctx) {
		var result struct {
			Name string `bson:"name"`
		}
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Warning: Failed to decode equipment option: %v", err)
			continue
		}
		equipment = append(equipment, result.Name)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error during unique equipment iteration: %w", err)
	}
	// No need to sort if SetSort is used in query.
	return equipment, nil
}
