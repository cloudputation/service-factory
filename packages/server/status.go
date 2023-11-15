package server

import (
  "net/http"
  "encoding/json"
  "log"

  "github.com/cloudputation/service-factory/packages/config"
  "github.com/cloudputation/service-factory/packages/storage"
  "github.com/cloudputation/service-factory/packages/network"
)


type StatusResponseBody struct {
  Message string `json:"message"`
}

func StatusHandlerWrapper(serviceName string, w http.ResponseWriter, r *http.Request) {
  StatusHandler(serviceName, w, r)
}

func StatusHandler(serviceName string, w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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

  // Retrieve data from Consul
  consulFactoryStatus, err := storage.ConsulStoreGet(network.ConsulClient, "service-factory/data/stats")
  if err != nil {
    http.Error(w, "Failed to fetch factory state: "+err.Error(), http.StatusInternalServerError)
    return
  }

  jsonData, err := json.MarshalIndent(consulFactoryStatus, "", "    ")
  if err != nil {
    log.Printf("Error marshaling JSON: %v", err)
    http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(jsonData)
}
