package api

import (
    "fmt"
    "net/http"

    "github.com/prometheus/client_golang/prometheus/promhttp"

    "github.com/cloudputation/service-factory/packages/config"
    l "github.com/cloudputation/service-factory/packages/logger"
    "github.com/cloudputation/service-factory/packages/api/v1"
)


const MaxWorkers = 10

func StartServer() {
  serverPort := fmt.Sprintf(":%s", config.AppConfig.Server.ServerPort)
  l.Info("Starting server on port %s", serverPort)

  http.HandleFunc("/v1/health", v1.HealthHandler)
  http.HandleFunc("/v1/service/apply", v1.ApplyHandlerWrapper)
  http.HandleFunc("/v1/service/destroy", v1.DestroyHandlerWrapper)
  http.HandleFunc("/v1/service/status", v1.ServiceStatusHandlerWrapper)
  http.HandleFunc("/v1/repo/status", v1.RepoStatusHandlerWrapper)
  http.HandleFunc("/v1/system/status", v1.SystemStatusHandlerWrapper)
  http.Handle("/v1/system/metrics", promhttp.Handler())

  err := http.ListenAndServe(serverPort, nil)
  if err != nil {
      l.Fatal("HTTP server failed to start: %v", err)
  }
}
