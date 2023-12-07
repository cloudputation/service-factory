package server

import (
  "fmt"
  "net/http"

  "github.com/cloudputation/service-factory/packages/config"
  l "github.com/cloudputation/service-factory/packages/logger"
)


const maxWorkers = 10

func StartServer() {
  serverPort := fmt.Sprintf(":%s", config.AppConfig.Server.ServerPort)
  l.Info("Starting server on port %s", serverPort)

  http.HandleFunc("/health", HealthHandler)
  http.HandleFunc("/apply", ApplyHandlerWrapper)
	http.HandleFunc("/destroy", DestroyHandlerWrapper)
  http.HandleFunc("/repo/status", RepoStatusHandlerWrapper)
  http.HandleFunc("/service/status", ServiceStatusHandlerWrapper)
  http.HandleFunc("/system/status", SystemStatusHandlerWrapper)

  err := http.ListenAndServe(serverPort, nil)
  if err != nil {
      l.Fatal("HTTP server failed to start: %v", err)
  }
}
