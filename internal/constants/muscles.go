package constants

import "strings"

const (
	MuscleGroupBack      = "Back"
	MuscleGroupBiceps    = "Biceps"
	MuscleGroupChest     = "Chest"
	MuscleGroupLegs      = "Legs"
	MuscleGroupShoulders = "Shoulders"
	MuscleGroupTriceps   = "Triceps"
	MuscleGroupObliques  = "Obliques"
	MuscleGroupAbs       = "Abs"
	MuscleGroupFullBody  = "Full Body"
)

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

func IsValidMuscleGroup(name string) bool {
	for _, validName := range AllMuscleGroupNames {
		if strings.EqualFold(validName, name) {
			return true
		}
	}
	return false
}
