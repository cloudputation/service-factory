package stats

import (
	"encoding/json"
	"io/ioutil"
)


// FactoryInfo struct holds factory state and services list
type FactoryInfo struct {
	FactoryState string   `json:"factory_state"`
	Services     []string `json:"services"`
}

// GetFactoryInfo scans the given directory for sub-directories
// and returns FactoryInfo as a JSON string
func GetFactoryInfo(dirPath string) (string, error) {
	// Read the directory
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	// Initialize a slice to store the names of directories
	var dirNames []string

	// Loop through each file and check if it's a directory
	for _, file := range files {
		if file.IsDir() {
			dirNames = append(dirNames, file.Name())
		}
	}

	// Create a FactoryInfo object
	factoryInfo := FactoryInfo{
		FactoryState: "running",
		Services:     dirNames,
	}

	// Marshal the FactoryInfo object to JSON
	jsonData, err := json.MarshalIndent(factoryInfo, "", "  ")
	if err != nil {
		return "", err
	}

	// Return the JSON data as a string
	return string(jsonData), nil
}
