package v1

import (
    "net/http"
    "io/ioutil"

    l "github.com/cloudputation/service-factory/packages/logger"
    "github.com/cloudputation/service-factory/packages/stats"
)


func HealthHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
      err := http.StatusMethodNotAllowed
      l.Error("Received an invalid request method: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
      return
  }
  stats.HealthEndpointCounter.Add(r.Context(), 1)


  content, err := ioutil.ReadFile("./VERSION")
  if err != nil {
      l.Error("Failed to read API version", http.StatusInternalServerError)
      http.Error(w, "Failed to read API version", http.StatusInternalServerError)
      return
  }

  response := string(content) + "OK\n"
  w.Write([]byte(response))
}
