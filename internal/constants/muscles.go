package constants

import "strings"

// MuscleGroup types as Go constants for type safety and validation
const (
	MuscleGroupBack      = "Back"
	MuscleGroupBiceps    = "Biceps"
	MuscleGroupChest     = "Chest"
	MuscleGroupLegs      = "Legs"
	MuscleGroupShoulders = "Shoulders"
	MuscleGroupTriceps   = "Triceps"
	MuscleGroupObliques  = "Obliques"
	MuscleGroupAbs       = "Abs"
	MuscleGroupFullBody  = "Full Body" // From your Deadlift example
	// Add any other muscle groups you use
)

// AllMuscleGroupNames provides a slice of all valid muscle group names.
// This is used for validation and for populating the `muscles` database table.
var AllMuscleGroupNames = []string{
	MuscleGroupBack,
	MuscleGroupBiceps,
	MuscleGroupChest,
	MuscleGroupLegs,
	MuscleGroupShoulders,
	MuscleGroupTriceps,
	MuscleGroupObliques,
	MuscleGroupAbs,
	MuscleGroupFullBody,
}

// IsValidMuscleGroup checks if a given string is a valid muscle group name (case-insensitive).
func IsValidMuscleGroup(name string) bool {
	for _, validName := range AllMuscleGroupNames {
		if strings.EqualFold(validName, name) {
			return true
		}
	}
	return false
}
