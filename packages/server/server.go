package server

import (
    "net/http"
    "log"
    "fmt"
)


func StartServer() {
  port := 48840
  serverPort := fmt.Sprintf(": %d", port)
  logMessage := fmt.Sprintf("Server started listenning to port%s", serverPort)
  log.Println(logMessage)

  http.HandleFunc("/health", HealthHandler)
  http.HandleFunc("/create", CreateHandler)
  http.ListenAndServe(serverPort, nil)

}
