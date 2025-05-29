package constants

import "strings"

// Equipment types as Go constants for type safety and validation
const (
	EquipmentBarbell             = "Barbell"
	EquipmentSquatRack           = "Squat Rack"
	EquipmentNone                = "None"
	EquipmentPullupBar           = "Pullup Bar"
	EquipmentDumbbells           = "Dumbbells"
	EquipmentWeightPlates        = "Weight Plates"
	EquipmentLatPulldownMachine  = "Lat Pulldown Machine"
	EquipmentSmithMachine        = "Smith Machine"
	EquipmentCableMachine        = "Cable Machine"
	EquipmentFlatBench           = "Flat Bench"
	EquipmentDeclineBench        = "Decline Bench"
	EquipmentInclineBench        = "Incline Bench"
	EquipmentChestPressMachine   = "Chest Press Machine"
	EquipmentLegCurlMachine      = "Leg Curl Machine"
	EquipmentLegExtensionMachine = "Leg Extension Machine"
	EquipmentLegPressMachine     = "Leg Press Machine"
	EquipmentEZBar               = "EZ Bar"
	// Add any other equipment you use
)

// AllEquipmentNames provides a slice of all valid equipment names.
// This is used for validation and for populating the `equipment` database table.
var AllEquipmentNames = []string{
	EquipmentBarbell,
	EquipmentSquatRack,
	EquipmentNone,
	EquipmentPullupBar,
	EquipmentDumbbells,
	EquipmentWeightPlates,
	EquipmentLatPulldownMachine,
	EquipmentSmithMachine,
	EquipmentCableMachine,
	EquipmentFlatBench,
	EquipmentDeclineBench,
	EquipmentInclineBench,
	EquipmentChestPressMachine,
	EquipmentLegCurlMachine,
	EquipmentLegExtensionMachine,
	EquipmentLegPressMachine,
	EquipmentEZBar,
}

// IsValidEquipment checks if a given string is a valid equipment name (case-insensitive).
func IsValidEquipment(name string) bool {
	for _, validName := range AllEquipmentNames {
		if strings.EqualFold(validName, name) {
			return true
		}
	}
	return false
}
