package server

import (
  "log"
  "fmt"
  "net/http"
  "encoding/json"
  "path/filepath"

  "github.com/cloudputation/service-factory/packages/stats"
  "github.com/cloudputation/service-factory/packages/utils"
  "github.com/cloudputation/service-factory/packages/config"
  "github.com/cloudputation/service-factory/packages/storage"
  "github.com/cloudputation/service-factory/packages/network"
)


type DeleteResponseBody struct {
	Message  string `json:"message"`
}


func DeleteHandlerWrapper(serviceName string, w http.ResponseWriter, r *http.Request) {
  DeleteHandler(serviceName, w, r)
}

func DeleteHandler(serviceName string, w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    return
  }

  serviceName = r.URL.Query().Get("serviceName")
  if serviceName == "" {
    http.Error(w, "Service name must be provided", http.StatusBadRequest)
    return
  }

  dataDir := config.AppConfig.DataDir
  if dataDir == "" {
    http.Error(w, "Data directory not configured", http.StatusInternalServerError)
    return
  }

  serviceBaseDir := filepath.Join(dataDir, "services", serviceName)
  serviceTerraformDir := filepath.Join(serviceBaseDir, "terraform")

  log.Printf("Got service: %s", serviceName)

	err := network.InitConsul(config.AppConfig.Consul.ConsulHost)
  if err != nil {
    return
  }

  if network.ConsulClient == nil {
    http.Error(w, "Consul client has not been initialized", http.StatusInternalServerError)
    return
  }

  // Retrieve data from Consul
	consulFactoryData, err := storage.ConsulStoreGet(network.ConsulClient, "service-factory/data/stats")
	if err != nil {
		http.Error(w, "Failed to fetch factory state: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the service name exists in the 'services' slice
	services, ok := consulFactoryData["services"].([]interface{})
	if !ok {
		http.Error(w, "Malformed factory state data", http.StatusInternalServerError)
		return
	}

	serviceExists := false
	for _, service := range services {
		if service.(string) == serviceName {
			serviceExists = true
			break
		}
	}

	if serviceExists {
    log.Printf("Service is present in factory: %s. Deleting.", serviceName)

    keyPath := fmt.Sprintf("service-factory/data/services/%s", serviceName)
  	consulServiceData, err := storage.ConsulStoreGet(network.ConsulClient, keyPath)
  	if err != nil {
  		http.Error(w, "Failed to fetch factory state: "+err.Error(), http.StatusInternalServerError)
  		return
  	}

    runnerID, _ := consulServiceData["ci_runner"].(string)
    repoToken, _ := consulServiceData["repo_token"].(string)
    log.Printf("runnerID: %s", runnerID)
    log.Printf("repoToken: %s", repoToken)

		terraformCmd := "destroy"
		err = utils.RunTerraform(serviceTerraformDir, terraformCmd, serviceName, runnerID, repoToken)
		if err != nil {
			http.Error(w, "Failed to run terraform: "+err.Error(), http.StatusInternalServerError)
			return
		}

    dirName := filepath.Join(dataDir, "services", serviceName)
    err = utils.DeleteDir(dirName)
    if err != nil {
      log.Fatalf("Could not delete service directory: %v", err)
    }

    err = storage.ConsulStoreDelete(network.ConsulClient, keyPath)
  	if err != nil {
  		log.Fatalf("Failed to get factory info: %v", err)
  	}

    err = stats.GenerateState()
  	if err != nil {
  		log.Fatalf("Failed to get factory info: %v", err)
  	}


	} else {
		http.Error(w, "Service does not exist in factory state, cannot destroy", http.StatusForbidden)
		return
	}


  json.NewEncoder(w).Encode(DeleteResponseBody{Message: "Operation completed successfully"})
}
