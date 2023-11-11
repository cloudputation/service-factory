package server

import (
  "fmt"
  "net/http"
  "encoding/json"
  "log"

  "github.com/cloudputation/service-factory/packages/config"
  "github.com/cloudputation/service-factory/packages/storage"
  "github.com/cloudputation/service-factory/packages/network"
)

type ServiceStatusResponseBody struct {
  Message string `json:"message"`
}

func ServiceStatusHandlerWrapper(serviceName string, w http.ResponseWriter, r *http.Request) {
  ServiceStatusHandler(serviceName, w, r)
}

func ServiceStatusHandler(serviceName string, w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    return
  }

  serviceName = r.URL.Query().Get("serviceName")
  if serviceName == "" {
    http.Error(w, "Service name must be provided", http.StatusBadRequest)
    return
  }

  err := network.InitConsul(config.AppConfig.Consul.ConsulHost)
  if err != nil {
    log.Printf("Failed to initialize Consul: %v", err)
    http.Error(w, "Failed to initialize Consul", http.StatusInternalServerError)
    return
  }

  if network.ConsulClient == nil {
    http.Error(w, "Consul client has not been initialized", http.StatusInternalServerError)
    return
  }

  keyPath := fmt.Sprintf("service-factory/data/services/%s", serviceName)
  consulServiceStatus, err := storage.ConsulStoreGet(network.ConsulClient, keyPath)
  if err != nil {
    http.Error(w, "Failed to fetch factory state: "+err.Error(), http.StatusInternalServerError)
    return
  }

  jsonData, err := json.MarshalIndent(consulServiceStatus, "", "    ")
  if err != nil {
    log.Printf("Error marshaling JSON: %v", err)
    http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(jsonData)
}
