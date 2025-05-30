package version

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"fitness-framework-api/internal/models"
)

const (
	VersionFilePath = "./data/version.json"
)

func LoadVersionInfo() (*models.ApiInfo, error) {
	absPath, err := filepath.Abs(VersionFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for %s: %w", VersionFilePath, err)
	}

	slog.Info("Attempting to load version info from", "path", absPath)

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read version file %s: %w", VersionFilePath, err)
	}

	var apiInfos []models.ApiInfo
	if err := json.Unmarshal(data, &apiInfos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal version data from %s: %w", VersionFilePath, err)
	}

	if len(apiInfos) == 0 {
		return nil, fmt.Errorf("version file %s is empty or contains no API info objects", VersionFilePath)
	}

	slog.Info("Successfully loaded version", "version", apiInfos[0].Version, "buildType", apiInfos[0].BuildType)
	return &apiInfos[0], nil
}
