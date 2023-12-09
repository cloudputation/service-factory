package v1

import (
    "fmt"
    "sync"
    "net/http"

    "github.com/cloudputation/service-factory/packages/config"
    l "github.com/cloudputation/service-factory/packages/logger"
  	"github.com/cloudputation/service-factory/packages/stats"
)


func RepoStatusHandlerWrapper(w http.ResponseWriter, r *http.Request) {
  RepoStatusHandler(w, r)
}

func RepoStatusHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
      err := http.StatusMethodNotAllowed
      l.Error("Received an invalid request method: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
      return
  }
  stats.RepoStatusEndpointCounter.Add(r.Context(), 1)

  var wg sync.WaitGroup

  wg.Add(1)
  go func() {
      defer wg.Done()
      ProcessRepoChecks(w, r)
  }()

  wg.Wait()
}

func ProcessRepoChecks(w http.ResponseWriter, r *http.Request){
  repoID        := r.URL.Query().Get("repoID")
  repoOwner     := r.URL.Query().Get("repoOwner")
  serviceName   := r.URL.Query().Get("serviceName")
  repoProvider  := r.URL.Query().Get("repoProvider")
  repoToken     := config.AppConfig.Repo.Gitlab.AccessToken


  var url string
  switch repoProvider {
  case "github":
      url = fmt.Sprintf("https://api.github.com/repos")
  case "gitlab":
      url = fmt.Sprintf("https://gitlab.com/api/v4/projects/%s", repoID)
  default:
      l.Error("Invalid repository provider", http.StatusBadRequest)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Invalid repository provider", http.StatusBadRequest)
      return
  }

  client := &http.Client{
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
          if repoProvider == "gitlab" && len(via) >= 10 {
              l.Error("Too many redirects encountered: %s", http.ErrUseLastResponse)
              return http.ErrUseLastResponse
          }
          return nil
      },
  }

  req, err := http.NewRequest("HEAD", url, nil)
  if err != nil {
      l.Error("Error creating request: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Internal server error", http.StatusInternalServerError)
      return
  }

  if repoToken != "" {
      req.Header.Add("PRIVATE-TOKEN", repoToken)
  }

  resp, err := client.Do(req)
  if err != nil {
      l.Error("Error making request: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Internal server error", http.StatusInternalServerError)
      return
  }
  defer resp.Body.Close()

  if resp.StatusCode == http.StatusOK {
      w.WriteHeader(http.StatusOK)
      response := fmt.Sprintf("Repository exists at: %s.com/%s/%s\n", repoProvider, repoOwner, serviceName)
      l.Info("[CONSUL CHECK] Repository exists at: %s.com/%s/%s\n", repoProvider, repoOwner, serviceName)
      w.Write([]byte(response))
  } else if resp.StatusCode == http.StatusNotFound {
      errorMsg := fmt.Sprintf("Repository not found at: %s.com/%s/%s\n", repoProvider, repoOwner, serviceName)
      l.Info("[CONSUL CHECK] Repository not found at: %s.com/%s/%s\n", repoProvider, repoOwner, serviceName)
      l.Error("%s: %s", errorMsg, http.StatusNotFound)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, errorMsg, http.StatusNotFound)
  } else {
      // 500 error if not found and not 404
      l.Error("Error checking repository: %s", http.StatusInternalServerError)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Error checking repository", http.StatusInternalServerError)
  }
}
