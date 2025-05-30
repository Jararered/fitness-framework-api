package constants

import "strings"

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
)

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

func IsValidEquipment(name string) bool {
	for _, validName := range AllEquipmentNames {
		if strings.EqualFold(validName, name) {
			return true
		}
	}
	return false
}
