package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"fitness-framework-api/internal/models"
)

const (
	VersionFilePath = "./data/version.json"
)

// LoadVersionInfo reads the version information from the JSON file.
// It now expects the JSON to be an array and takes the first element.
func LoadVersionInfo() (*models.ApiInfo, error) {
	absPath, err := filepath.Abs(VersionFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for %s: %w", VersionFilePath, err)
	}

	log.Printf("Attempting to load version info from: %s", absPath)

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read version file %s: %w", VersionFilePath, err)
	}

	// *** CRUCIAL CHANGE HERE: Unmarshal into a slice of ApiInfo ***
	var apiInfos []models.ApiInfo
	if err := json.Unmarshal(data, &apiInfos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal version data from %s: %w", VersionFilePath, err)
	}

	// Check if the slice is empty (e.g., if JSON was "[]")
	if len(apiInfos) == 0 {
		return nil, fmt.Errorf("version file %s is empty or contains no API info objects", VersionFilePath)
	}

	// *** Return the first element of the slice ***
	log.Printf("Successfully loaded version: %s (buildType: %s)", apiInfos[0].Version, apiInfos[0].BuildType)
	return &apiInfos[0], nil
}
