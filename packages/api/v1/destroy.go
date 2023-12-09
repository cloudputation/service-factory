package v1

import (
    "encoding/json"
    "fmt"
    "net/http"
    "path/filepath"
    "os/exec"
    "sync"

    "github.com/cloudputation/service-factory/packages/config"
    "github.com/cloudputation/service-factory/packages/consul"
    l "github.com/cloudputation/service-factory/packages/logger"
    "github.com/cloudputation/service-factory/packages/stats"
    "github.com/cloudputation/service-factory/packages/terraform"
)

type DestroyResponseBody struct {
    Message          string   `json:"message"`
    DestroyedServices  []string `json:"destroyedServices"`
}

type DestroyRequestBody struct {
    ServiceNames []string `json:"service-names"`
}

func DestroyHandlerWrapper(w http.ResponseWriter, r *http.Request) {
  DestroyHandler(w, r)
}

func DestroyHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
      err := http.StatusMethodNotAllowed
      l.Error("Received an invalid request method: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
      return
  }
  stats.DestroyEndpointCounter.Add(r.Context(), 1)

  var requestBody DestroyRequestBody
  if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
      err := http.StatusMethodNotAllowed
      l.Error("Received an invalid request method: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Error reading request body", err)
      return
  }

  jobs := make(chan string, len(requestBody.ServiceNames))
  errors := make(chan error, len(requestBody.ServiceNames))
  var wg sync.WaitGroup
  destroyedServices := make([]string, 0)
  var mutex sync.Mutex

  for w := 0; w < config.MaxWorkers; w++ {
      wg.Add(1)
      go func() {
          defer wg.Done()
          for serviceName := range jobs {
              if err := ProcessDeletion(serviceName, &destroyedServices, &mutex); err != nil {
                  errors <- err
              }
          }
      }()
  }

  for _, service := range requestBody.ServiceNames {
      jobs <- service
  }
  close(jobs)

  wg.Wait()
  close(errors)

  // Check if any errors occurred
  var hadError bool
  for err := range errors {
      if err != nil {
          l.Error("Failed to delete service: %v", err)
          stats.ErrorCounter.Add(r.Context(), 1)
          hadError = true
      }
  }

  if hadError {
      http.Error(w, "Failed to delete some services", http.StatusInternalServerError)
      return
  }

  json.NewEncoder(w).Encode(DestroyResponseBody{
      Message: "Services deleted successfully",
      DestroyedServices: destroyedServices,
  })
}

func ProcessDeletion(serviceName string, destroyedServices *[]string, mutex *sync.Mutex) error {
  dataDir := config.AppConfig.DataDir
  if dataDir == "" {
      return fmt.Errorf("Data directory not configured")
  }

  consulServiceDir := config.ConsulServicesDataDir
  keyPath := fmt.Sprintf("%s/%s", consulServiceDir, serviceName)
  consulServiceData, err := consul.ConsulStoreGet(keyPath)
  if err != nil {
      return fmt.Errorf("Failed to fetch service manifest: %v", err)
  }

  var serviceID, _ = consulServiceData["service-id"].(string)
  terraformServiceDir := filepath.Join(dataDir, "services", serviceName, "terraform")
  terraformCmd := "destroy"
  if err := terraform.RunTerraform(terraformServiceDir, terraformCmd, serviceID, serviceName); err != nil {
      return fmt.Errorf("Failed to run terraform: %v", err)
  }

  if err := consul.DeregisterRepo(serviceID); err != nil {
      return fmt.Errorf("Failed to deregister repository service: %v", err)
  }

  if err := consul.ConsulStoreDelete(keyPath); err != nil {
      return fmt.Errorf("Failed to delete service manifest: %v", err)
  }

  dirName := filepath.Join(dataDir, "services", serviceName)
  cmd := exec.Command("rm", "-fr", dirName)
  err = cmd.Run()
  if err != nil {
      return fmt.Errorf("Failed to delete service directory: %v", err)
  }
  if err := stats.GenerateState(); err != nil {
      return fmt.Errorf("Failed to get factory info: %v", err)
  }

  mutex.Lock()
  *destroyedServices = append(*destroyedServices, serviceName)
  mutex.Unlock()

  return nil
}
