package stats

import (
  "log"

  "github.com/cloudputation/service-factory/packages/config"
  "github.com/cloudputation/service-factory/packages/storage"
)


func GenerateState() error {
  dirName := config.AppConfig.DataDir + "/services"
  jsonData, err := GetFactoryInfo(dirName)
	if err != nil {
		log.Fatalf("Failed to get factory info: %v", err)
	}

	keyPath := "service-factory/data/stats"
	err = storage.ConsulStorePut(jsonData, keyPath)
	if err != nil {
		log.Fatalf("Failed to put data in Consul KV: %v", err)
	}


  return nil
}
