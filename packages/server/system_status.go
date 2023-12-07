package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudputation/service-factory/packages/config"
	"github.com/cloudputation/service-factory/packages/consul"
	l "github.com/cloudputation/service-factory/packages/logger"
)

type FactoryStatus struct {
	Status string `json:"factory-status"`
}

type SystemStatusResponseBody struct {
	Message string `json:"message"`
}

func SystemStatusHandlerWrapper(w http.ResponseWriter, r *http.Request) {
	SystemStatusHandler(w, r)
}

func SystemStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
      err := http.StatusMethodNotAllowed
      l.Error("Received an invalid request method: %v", err)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
	}

	factoryDataPath := config.ConsulFactoryDataDir
	statusPath := fmt.Sprintf("%s/status", factoryDataPath)

	consulFactoryStatus, err := consul.ConsulStoreGet(statusPath)
	if err != nil {
			l.Error("Failed to fetch factory state: "+err.Error(), http.StatusInternalServerError)
			http.Error(w, "Failed to fetch factory state: "+err.Error(), http.StatusInternalServerError)
			return
	}

	consulFactoryStatusBytes, err := json.Marshal(consulFactoryStatus)
	if err != nil {
			l.Error("Error marshaling map to JSON: %v", err)
			http.Error(w, "Failed to marshal map to JSON", http.StatusInternalServerError)
			return
	}

	var status FactoryStatus
	err = json.Unmarshal(consulFactoryStatusBytes, &status)
	if err != nil {
			l.Error("Error unmarshaling JSON: %v", err)
			http.Error(w, "Failed to unmarshal JSON", http.StatusInternalServerError)
			return
	}


	w.Header().Set("Content-Type", "text/plain")
  fmt.Fprintf(w, status.Status)
}
