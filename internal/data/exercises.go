package data

import (
	"fitness-framework-api/internal/constants" // Import your constants package
	"fitness-framework-api/internal/models"
)

// AllRawExercises is the canonical, type-safe list of exercises defined directly in Go.
// You can modify this list directly to add, remove, or change exercises.
var AllRawExercises = []models.RawExercise{
	{
		Name:      "Barbell Bent-Over Row",
		Equipment: []string{constants.EquipmentBarbell},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Barbell Deadlift",
		Equipment: []string{constants.EquipmentBarbell}, // Assuming Barbell is enough, otherwise add WeightPlates
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Dumbbell Bent-Over Row",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Dumbbell Single-Arm Row",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Lat Pulldown",
		Equipment: []string{constants.EquipmentLatPulldownMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Lat Pulldown Behind-the-Head",
		Equipment: []string{constants.EquipmentLatPulldownMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Lat Pulldown Narrow-Grip",
		Equipment: []string{constants.EquipmentLatPulldownMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Lat Pulldown Reverse-Grip",
		Equipment: []string{constants.EquipmentLatPulldownMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Lat Pulldown Single-Arm",
		Equipment: []string{constants.EquipmentLatPulldownMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Lat Pulldown Wide-Grip",
		Equipment: []string{constants.EquipmentLatPulldownMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Smith Machine Bent-Over Row",
		Equipment: []string{constants.EquipmentSmithMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Pullup",
		Equipment: []string{constants.EquipmentPullupBar},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Cable Lat Pulldown",
		Equipment: []string{constants.EquipmentCableMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Seated Cable Row",
		Equipment: []string{constants.EquipmentCableMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Standing Cable Pullover",
		Equipment: []string{constants.EquipmentCableMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Cable Face Pull",
		Equipment: []string{constants.EquipmentCableMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Straight Arm Cable Pulldown",
		Equipment: []string{constants.EquipmentCableMachine},
		Muscles:   []string{constants.MuscleGroupBack},
	},
	{
		Name:      "Barbell Curl",
		Equipment: []string{constants.EquipmentBarbell},
		Muscles:   []string{constants.MuscleGroupBiceps},
	},
	{
		Name:      "Dumbbell Bicep Curl",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupBiceps},
	},
	{
		Name:      "Dumbbell Hammer Curl",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupBiceps},
	},
	{
		Name:      "EZ Bar Close-Grip Curl",
		Equipment: []string{constants.EquipmentEZBar},
		Muscles:   []string{constants.MuscleGroupBiceps},
	},
	{
		Name:      "EZ Bar Curl",
		Equipment: []string{constants.EquipmentEZBar},
		Muscles:   []string{constants.MuscleGroupBiceps},
	},
	{
		Name:      "EZ Bar Wide-Grip Curl",
		Equipment: []string{constants.EquipmentEZBar},
		Muscles:   []string{constants.MuscleGroupBiceps},
	},
	{
		Name:      "Barbell Bench Press",
		Equipment: []string{constants.EquipmentBarbell, constants.EquipmentFlatBench},
		Muscles:   []string{constants.MuscleGroupChest},
	},
	{
		Name:      "Barbell Decline Bench Press",
		Equipment: []string{constants.EquipmentBarbell, constants.EquipmentDeclineBench},
		Muscles:   []string{constants.MuscleGroupChest},
	},
	{
		Name:      "Barbell Incline Bench Press",
		Equipment: []string{constants.EquipmentBarbell, constants.EquipmentInclineBench},
		Muscles:   []string{constants.MuscleGroupChest},
	},
	{
		Name:      "Machine Chest Press",
		Equipment: []string{constants.EquipmentChestPressMachine},
		Muscles:   []string{constants.MuscleGroupChest},
	},
	{
		Name:      "Machine Incline Chest Press",
		Equipment: []string{constants.EquipmentChestPressMachine},
		Muscles:   []string{constants.MuscleGroupChest},
	},
	{
		Name:      "Smith Machine Bench Press",
		Equipment: []string{constants.EquipmentSmithMachine, constants.EquipmentFlatBench},
		Muscles:   []string{constants.MuscleGroupChest},
	},
	{
		Name:      "Smith Machine Close-Grip Bench Press",
		Equipment: []string{constants.EquipmentSmithMachine, constants.EquipmentFlatBench},
		Muscles:   []string{constants.MuscleGroupChest},
	},
	{
		Name:      "Smith Machine Incline Bench Press",
		Equipment: []string{constants.EquipmentSmithMachine, constants.EquipmentInclineBench},
		Muscles:   []string{constants.MuscleGroupChest},
	},
	{
		Name:      "Barbell Front Squat",
		Equipment: []string{constants.EquipmentBarbell},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Barbell Lunge",
		Equipment: []string{constants.EquipmentBarbell},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Barbell Romanian Deadlift",
		Equipment: []string{constants.EquipmentBarbell},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Barbell Squat",
		Equipment: []string{constants.EquipmentBarbell},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Barbell Sumo Deadlift",
		Equipment: []string{constants.EquipmentBarbell},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Dumbbell Lunges",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Leg Curl",
		Equipment: []string{constants.EquipmentLegCurlMachine},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Leg Curl + Isometric Hold",
		Equipment: []string{constants.EquipmentLegCurlMachine},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Leg Extension",
		Equipment: []string{constants.EquipmentLegExtensionMachine},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Leg Extension + Isometric Hold",
		Equipment: []string{constants.EquipmentLegExtensionMachine},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Leg Press",
		Equipment: []string{constants.EquipmentLegPressMachine},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Smith Machine Lunges",
		Equipment: []string{constants.EquipmentSmithMachine},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Smith Machine Squat",
		Equipment: []string{constants.EquipmentSmithMachine},
		Muscles:   []string{constants.MuscleGroupLegs},
	},
	{
		Name:      "Barbell Overhead Press",
		Equipment: []string{constants.EquipmentBarbell},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbbell Front Raise",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbbell Lateral Raise",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbbell Shoulder Press",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbbell Shrugs",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbbell Upright Row",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Machine Shoulder Press",
		Equipment: []string{constants.EquipmentChestPressMachine}, // Assuming same machine for chest and shoulder press
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Smith Machine Shoulder Press",
		Equipment: []string{constants.EquipmentSmithMachine},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Smith Machine Shrugs",
		Equipment: []string{constants.EquipmentSmithMachine},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbell Seated Military Press",
		Equipment: []string{constants.EquipmentDumbbells, constants.EquipmentFlatBench},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbbell Standing Military Press",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbbell Seated Arnold Press",
		Equipment: []string{constants.EquipmentDumbbells, constants.EquipmentFlatBench},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbbell Seated Reverse Arnold Press",
		Equipment: []string{constants.EquipmentDumbbells, constants.EquipmentFlatBench},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbbell Seated Arnold Rotations",
		Equipment: []string{constants.EquipmentDumbbells, constants.EquipmentFlatBench},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbbell Hammer Shrug",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupShoulders},
	},
	{
		Name:      "Dumbbell Skull Crushers",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupTriceps},
	},
	{
		Name:      "Dumbbell Tricep Extension",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupTriceps},
	},
	{
		Name:      "EZ Bar Skull Crusher",
		Equipment: []string{constants.EquipmentEZBar, constants.EquipmentFlatBench},
		Muscles:   []string{constants.MuscleGroupTriceps},
	},
	{
		Name:      "Rope Pulldown",
		Equipment: []string{constants.EquipmentCableMachine},
		Muscles:   []string{constants.MuscleGroupTriceps},
	},
	{
		Name:      "Dumbbell Side Oblique Crunch",
		Equipment: []string{constants.EquipmentDumbbells},
		Muscles:   []string{constants.MuscleGroupObliques},
	},
	{
		Name:      "Cable Straight Arm Oblique Twist",
		Equipment: []string{constants.EquipmentCableMachine},
		Muscles:   []string{constants.MuscleGroupObliques},
	},
	{
		Name:      "Cable Bent Over Oblique Dig",
		Equipment: []string{constants.EquipmentCableMachine},
		Muscles:   []string{constants.MuscleGroupObliques},
	},
	{
		Name:      "Rope Crunch",
		Equipment: []string{constants.EquipmentCableMachine},
		Muscles:   []string{constants.MuscleGroupAbs},
	},
	{
		Name:      "Hanging Leg Raise",
		Equipment: []string{constants.EquipmentPullupBar},
		Muscles:   []string{constants.MuscleGroupAbs},
	},
	{
		Name:      "Hanging Knee Raise",
		Equipment: []string{constants.EquipmentPullupBar},
		Muscles:   []string{constants.MuscleGroupAbs},
	},
	{
		Name:      "Weighted Hanging Knee Raise",
		Equipment: []string{constants.EquipmentPullupBar},
		Muscles:   []string{constants.MuscleGroupAbs},
	},
	// Add other exercises here, using constants.EquipmentX and constants.MuscleGroupY
}
