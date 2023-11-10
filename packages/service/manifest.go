package service

import (
  "fmt"
  "encoding/json"

  "github.com/cloudputation/service-factory/packages/storage"
)


type ManifestData struct {
	ServiceName  string   `json:"service_name"`
	CIRunner     string   `json:"ci_runner"`
	RepoToken    string   `json:"repo_token"`
}

var serviceSpecs RequestBody


func GenerateManifest(serviceName string, repoRunnerID string, repoToken string) error {
  jsonData, err := GetManifestData(serviceName, repoRunnerID, repoToken)
  if err != nil {
    return fmt.Errorf("Failed to get factory info: %v", err)
  }

  keyPath := fmt.Sprintf("service-factory/data/services/%s", serviceName)
  err = storage.ConsulStorePut(jsonData, keyPath)
  if err != nil {
    return fmt.Errorf("Failed to put data in Consul KV: %v", err)
  }

  fmt.Printf(keyPath)
  fmt.Printf(jsonData)
  return nil
}


func GetManifestData(serviceName string, repoRunnerID string, repoToken string) (string, error) {

	// Create a ManifestData object
	manifestData := ManifestData{
		ServiceName:  serviceName,
		CIRunner:     repoRunnerID,
		RepoToken:    repoToken,
	}

	// Marshal the ManifestData object to JSON
	jsonData, err := json.MarshalIndent(manifestData, "", "  ")
	if err != nil {
		return "", err
	}

	// Return the JSON data as a string
	return string(jsonData), nil
}
