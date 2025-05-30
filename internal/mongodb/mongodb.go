package mongodb

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fitness-framework-api/internal/constants"
	"fitness-framework-api/internal/data"
	"fitness-framework-api/internal/models"
)

const (
	MongoURI       = "mongodb://localhost:27017"
	CollectionName = "exercises"
)

func InitDB(databaseName string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(MongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		client.Disconnect(context.Background())
		return nil, fmt.Errorf("error pinging MongoDB: %w", err)
	}

	slog.Info("MongoDB connection established.")
	db := client.Database(databaseName)
	collection := db.Collection(CollectionName)

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		client.Disconnect(context.Background())
		return nil, fmt.Errorf("error counting documents in exercises collection: %w", err)
	}

	if count == 0 {
		slog.Info("Exercises collection is empty. Attempting to populate initial data from hardcoded Go data...")

		equipmentCollection := db.Collection("equipment_options")
		muscleCollection := db.Collection("muscles_options")

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := equipmentCollection.Drop(ctx); err != nil {
			slog.Error("Warning: Could not drop equipment_options collection", "error", err)
		}

		if err := muscleCollection.Drop(ctx); err != nil {
			slog.Error("Warning: Could not drop muscles_options collection", "error", err)
		}

		var equipmentDocs []interface{}
		for _, name := range constants.AllEquipmentNames {
			equipmentDocs = append(equipmentDocs, bson.D{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "name", Value: name}})
		}

		var muscleDocs []interface{}
		for _, name := range constants.AllMuscleGroupNames {
			muscleDocs = append(muscleDocs, bson.D{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "name", Value: name}})
		}

		if len(equipmentDocs) > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if _, err := equipmentCollection.InsertMany(ctx, equipmentDocs); err != nil {
				client.Disconnect(context.Background())
				return nil, fmt.Errorf("error inserting equipment options: %w", err)
			}
			slog.Info("Inserted equipment items from constants", "count", len(equipmentDocs))
		}

		if len(muscleDocs) > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if _, err := muscleCollection.InsertMany(ctx, muscleDocs); err != nil {
				client.Disconnect(context.Background())
				return nil, fmt.Errorf("error inserting muscle options: %w", err)
			}
			slog.Info("Inserted muscle items from constants", "count", len(muscleDocs))
		}

		var documents []interface{}
		for _, rawEx := range data.AllRawExercises {
			validEquipment := []string{}
			for _, eqName := range rawEx.Equipment {
				if constants.IsValidEquipment(eqName) {
					validEquipment = append(validEquipment, eqName)
				} else {
					slog.Error("Warning: Invalid equipment '%s' for exercise '%s'. Skipping", eqName, rawEx.Name)
				}
			}

			validMuscles := []string{}
			for _, muscleName := range rawEx.Muscles {
				if constants.IsValidMuscleGroup(muscleName) {
					validMuscles = append(validMuscles, muscleName)
				} else {
					slog.Error("Warning: Invalid muscle group '%s' for exercise '%s'. Skipping", muscleName, rawEx.Name)
				}
			}

			doc := models.Exercise{
				ID:        primitive.NewObjectID(),
				Name:      rawEx.Name,
				Equipment: validEquipment,
				Muscles:   validMuscles,
			}
			documents = append(documents, doc)
		}

		if len(documents) > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			_, err := collection.InsertMany(ctx, documents)
			if err != nil {
				client.Disconnect(context.Background())
				return nil, fmt.Errorf("error inserting initial exercises: %w", err)
			}
			slog.Info("Successfully populated exercises", "count", len(documents))
		}

	} else {
		slog.Info("Exercises collection already contains documents", "count", count)
	}

	return db, nil
}

func GetExercises(db *mongo.Database) ([]models.Exercise, error) {
	collection := db.Collection(CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to find exercises: %w", err)
	}
	defer cursor.Close(ctx)

	var exercises []models.Exercise
	if err = cursor.All(ctx, &exercises); err != nil {
		return nil, fmt.Errorf("failed to decode exercises: %w", err)
	}

	for i := range exercises {
		sort.Strings(exercises[i].Equipment)
		sort.Strings(exercises[i].Muscles)
	}

	return exercises, nil
}

func GetUniqueExerciseNames(db *mongo.Database) ([]string, error) {
	collection := db.Collection(CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

func GetUniqueMuscles(db *mongo.Database) ([]string, error) {
	collection := db.Collection("muscles_options")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
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
			slog.Error("Warning: Failed to decode muscle option", "error", err)
			continue
		}
		muscles = append(muscles, result.Name)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error during unique muscles iteration: %w", err)
	}

	return muscles, nil
}

func GetUniqueEquipment(db *mongo.Database) ([]string, error) {
	collection := db.Collection("equipment_options")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
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
			slog.Error("Warning: Failed to decode equipment option", "error", err)
			continue
		}
		equipment = append(equipment, result.Name)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error during unique equipment iteration: %w", err)
	}

	return equipment, nil
}
