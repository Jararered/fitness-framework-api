package models

// ApiInfo contains version information about the API
type ApiInfo struct {
	Version   string `json:"version"`
	BuildType string `json:"buildType"`
}
