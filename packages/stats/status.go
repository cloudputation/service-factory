package stats

import (
    "fmt"
  	"io/ioutil"
  	"encoding/json"

    "github.com/cloudputation/service-factory/packages/config"
    "github.com/cloudputation/service-factory/packages/consul"
)


type FactoryInfo struct {
	FactoryState string   `json:"factory-state"`
	Services     []string `json:"services"`
}

func GenerateState() error {
  dirName := config.AppConfig.DataDir + "/services"
  jsonData, err := GetFactoryInfo(dirName)
	if err != nil {
		return fmt.Errorf("Failed to get factory info: %v", err)
	}

	keyPath := config.ConsulServiceSummaryDataDir
	err = consul.ConsulStorePut(jsonData, keyPath)
	if err != nil {
		return fmt.Errorf("Failed to put data in Consul KV: %v", err)
	}


  return nil
}

func GetFactoryInfo(dirPath string) (string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("Failed to read factory files in %s: %v", dirPath, err)
	}

	var dirNames []string
	for _, file := range files {
		if file.IsDir() {
			dirNames = append(dirNames, file.Name())
		}
	}

	factoryInfo := FactoryInfo{
		FactoryState: "running",
		Services:     dirNames,
	}

	jsonData, err := json.MarshalIndent(factoryInfo, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Failed to process factory info: %v", err)
	}


	return string(jsonData), nil
}
