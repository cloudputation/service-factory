package server

import (
    "net/http"
    "io/ioutil"
)


func HealthHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    return
  }

  content, err := ioutil.ReadFile("/app/API_VERSION")
  if err != nil {
    http.Error(w, "Failed to read API version", http.StatusInternalServerError)
    return
  }

  response := string(content) + " OK"
  w.Write([]byte(response))
}
