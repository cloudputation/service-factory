package server

import (
    "net/http"
    "io/ioutil"

    l "github.com/cloudputation/service-factory/packages/logger"
)


func HealthHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
      err := http.StatusMethodNotAllowed
      l.Error("Received an invalid request method: %v", err)
      http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
      return
  }

  content, err := ioutil.ReadFile("./API_VERSION")
  if err != nil {
      l.Error("Failed to read API version", http.StatusInternalServerError)
      http.Error(w, "Failed to read API version", http.StatusInternalServerError)
      return
  }

  response := string(content) + "OK\n"
  w.Write([]byte(response))
}
