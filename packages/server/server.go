package server

import (
  "net/http"
  "log"
  "fmt"
)


func StartServer() {
  var serviceName string

  port := 48840
  serverPort := fmt.Sprintf(":%d", port)
  log.Printf("Starting server on port %s", serverPort)

  http.HandleFunc("/health", HealthHandler)
  http.HandleFunc("/create", CreateHandlerWrapper)
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		DeleteHandlerWrapper(serviceName, w, r)
	})
	http.HandleFunc("/system/status", func(w http.ResponseWriter, r *http.Request) {
		StatusHandlerWrapper(serviceName, w, r)
	})
	http.HandleFunc("/service/status", func(w http.ResponseWriter, r *http.Request) {
		ServiceStatusHandlerWrapper(serviceName, w, r)
	})
  err := http.ListenAndServe(serverPort, nil)
  if err != nil {
      log.Fatalf("HTTP server failed: %v", err)
  }
}
