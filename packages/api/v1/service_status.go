package v1

import (
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/cloudputation/service-factory/packages/config"
    "github.com/cloudputation/service-factory/packages/consul"
    "github.com/cloudputation/service-factory/packages/stats"
    l "github.com/cloudputation/service-factory/packages/logger"
)


type ServiceStatusResponseBody struct {
  Message string `json:"message"`
}

func ServiceStatusHandlerWrapper(w http.ResponseWriter, r *http.Request) {
  ServiceStatusHandler(w, r)
}

func ServiceStatusHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
      err := http.StatusMethodNotAllowed
      l.Error("Received an invalid request method: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
      return
  }
  stats.ServiceStatusEndpointCounter.Add(r.Context(), 1)

  serviceName := r.URL.Query().Get("serviceName")
  if serviceName == "" {
      l.Error("Failed to read API version", http.StatusBadRequest)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Service name must be provided", http.StatusBadRequest)
      return
  }

  consulServiceDir := config.ConsulServicesDataDir
  keyPath := fmt.Sprintf("%s/%s", consulServiceDir, serviceName)
  consulServiceStatus, err := consul.ConsulStoreGet(keyPath)
  if err != nil {
      l.Error("Error marshaling JSON: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Failed to fetch factory state: "+err.Error(), http.StatusInternalServerError)
      return
  }

  jsonData, err := json.MarshalIndent(consulServiceStatus, "", "    ")
  if err != nil {
      l.Error("Error marshaling JSON: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
      return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(jsonData)
}
