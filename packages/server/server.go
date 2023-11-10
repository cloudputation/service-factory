package server

import (
    "net/http"
    "log"
    "fmt"
)


func StartServer() {
  var serviceName string

  port := 48840
  serverPort := fmt.Sprintf(": %d", port)
  logMessage := fmt.Sprintf("Agent started on port%s", serverPort)
  log.Println(logMessage)

  http.HandleFunc("/health", HealthHandler)
  http.HandleFunc("/create", CreateHandlerWrapper)
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		DeleteHandlerWrapper(serviceName, w, r)
	})
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		StatusHandlerWrapper(serviceName, w, r)
	})
	http.HandleFunc("/service/status", func(w http.ResponseWriter, r *http.Request) {
		ServiceStatusHandlerWrapper(serviceName, w, r)
	})
  http.ListenAndServe(serverPort, nil)

}
