package v1

import (
    "os"
    "fmt"
    "sync"
    "strings"
    "os/exec"
    "net/http"
    "path/filepath"
    "encoding/json"
    "github.com/google/uuid"

    "github.com/cloudputation/service-factory/packages/config"
    l "github.com/cloudputation/service-factory/packages/logger"
    "github.com/cloudputation/service-factory/packages/stats"
    "github.com/cloudputation/service-factory/packages/service"
    "github.com/cloudputation/service-factory/packages/repository"
    "github.com/cloudputation/service-factory/packages/consul"
    "github.com/cloudputation/service-factory/packages/terraform"
)


type ApplyResponseBody struct {
    Message string   `json:"message"`
    AppliedServices []string `json:"appliedServices"`
}

type ServicesWrapper struct {
    Services []service.ServiceSpecs `json:"services"`
}


var spec service.ServiceSpecs
var terraformToApply []string
var runTerraformMutex sync.Mutex


func ApplyHandlerWrapper(w http.ResponseWriter, r *http.Request) {
  ApplyHandler(w, r)
}

func ApplyHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
      err := http.StatusMethodNotAllowed
      l.Error("Received an invalid request method: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
      return
  }
  stats.ApplyEndpointCounter.Add(r.Context(), 1)

  var wrapper ServicesWrapper
  if err := json.NewDecoder(r.Body).Decode(&wrapper); err != nil {
      err := http.StatusMethodNotAllowed
      l.Error("Received an invalid request method: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Error reading request body", http.StatusBadRequest)
      return
  }

  appliedServices := make([]string, 0)
  var mutex sync.Mutex
  jobs := make(chan service.ServiceSpecs, len(wrapper.Services))
  errors := make(chan error, len(wrapper.Services))
  var wg sync.WaitGroup

  for w := 0; w < config.MaxWorkers; w++ {
      wg.Add(1)
      go worker(&wg, jobs, &appliedServices, &mutex, errors)
  }

  for _, serviceSpec := range wrapper.Services {
      jobs <- serviceSpec
  }
  close(jobs)

  wg.Wait()
  close(errors)

  // Check if any errors occurred
  var hadError bool
  for err := range errors {
      if err != nil {
          l.Error("Failed to process service spec: %v", err)
          stats.ErrorCounter.Add(r.Context(), 1)
          hadError = true
      }
  }

  if hadError {
      http.Error(w, "Failed to process some service specifications", http.StatusInternalServerError)
      return
  }

  json.NewEncoder(w).Encode(ApplyResponseBody{
      Message: "Operation completed successfully",
      AppliedServices: appliedServices,
  })
}

func worker(wg *sync.WaitGroup, jobs <-chan service.ServiceSpecs, appliedServices *[]string, mutex *sync.Mutex, errors chan<- error) {
  defer wg.Done()
  for spec := range jobs {
      if err := processServiceSpec(spec, appliedServices, mutex); err != nil {
          errors <- err
      }
  }
}

func processServiceSpec(spec service.ServiceSpecs, appliedServices *[]string, mutex *sync.Mutex) error {
  var (
      UUID                  = uuid.New()
      dataDir               = config.AppConfig.DataDir
      SFHost                = config.AppConfig.Server.ServerAddress
      SFPort                = config.AppConfig.Server.ServerPort
      terraformDir          = config.AppConfig.Terraform.TerraformDir
      rootDir               = config.RootDir
      datastoreDir          = config.DatastoreDir
      serviceName           = spec.Service.Name
      serviceTags           = spec.Service.Tags
      templateURL           = spec.Service.Template.TemplateURL
      templateName          = spec.Service.Template.TemplateName
      repoOwner             = spec.Service.Repo.RepositoryOwner
      repoProvider          = spec.Service.Repo.Provider
      runnerID              = spec.Service.Repo.RunnerID
      serviceID             = fmt.Sprintf("%s-%s", serviceName, UUID)
      serviceBaseDir        = filepath.Join(dataDir, "services", serviceName)
      serviceRepoDir        = filepath.Join(serviceBaseDir, "repo")
      terraformProviderDir  = filepath.Join(rootDir, terraformDir, "providers", repoProvider)
      terraformTemplateDir  = filepath.Join(terraformProviderDir, "service-template")
      terraformServiceDir   = filepath.Join(dataDir, "services", serviceName, "terraform")
      terraformModuleDir    = filepath.Join(terraformServiceDir, "terraform_module")
  )


  if dataDir == "" {
      return fmt.Errorf("Data directory not configured")
  }

  // cmd := exec.Command("mkdir", "-p", serviceRepoDir)
  // err := cmd.Run()
  // if err != nil {
  //     return fmt.Errorf("Failed to create service directory: %v", err)
  // }
  //
  err := os.MkdirAll(serviceRepoDir, 0755)
  if err != nil {
      return fmt.Errorf("Failed to create service directory '%s': %v", serviceRepoDir, err)
  }

  parts := strings.Split(templateName, "/")
  templateRepoOwner := parts[0]
  templateRepoName := parts[1]
  parts = strings.Split(templateURL, ".com")
  templateRepoProvider := parts[0]

  l.Info("Got template repo provider: %s\n", templateRepoProvider)
  l.Info("Got template repo namespace: %s\n", templateRepoOwner)
  l.Info("Got template repo name: %s\n", templateRepoName)

  repoDatastorePath := filepath.Join(datastoreDir, templateRepoProvider, templateRepoOwner, templateRepoName)
  repoDatastoreParentPath := filepath.Join(datastoreDir, templateRepoProvider, templateRepoOwner)

  l.Info("Got repo at datastore path: %s\n", repoDatastorePath)
  l.Info("Got repo at datastore parent path: %s\n", repoDatastoreParentPath)

  err = repository.DownloadRepoToDatastore(repoDatastorePath, repoDatastoreParentPath, templateURL, templateName)
  if err != nil {
      return fmt.Errorf("Failed download repository to datastore: %v", err)
  }


  cmd := exec.Command("cp", "-r", repoDatastorePath+"/.", serviceRepoDir)
  err = cmd.Run()
  if err != nil {
      return fmt.Errorf("Failed to copy repo template to service directory: %v", err)
  }

  gitRepoToDelete := filepath.Join(serviceRepoDir, ".git")
  cmd = exec.Command("rm", "-fr", gitRepoToDelete)
  err = cmd.Run()
  if err != nil {
      return fmt.Errorf("Failed to remove old .git directory in service repo: %v", err)
  }

  cookieCutterFiles, err := service.GetCookieCutterFiles(serviceRepoDir)
  if err != nil {
      return fmt.Errorf("Failed to get CookieCutter files: %v", err)
  }

  for tplFile, outputFile := range cookieCutterFiles {
      err = service.RenderTemplate(tplFile, outputFile, spec)
      if err != nil {
          return fmt.Errorf("Failed to render the template: %v", err)
      }
  }
  l.Info("Rending done for service: %s", serviceName)

  cmd = exec.Command("cp", "-r", terraformTemplateDir+"/.", terraformServiceDir)
  err = cmd.Run()
  if err != nil {
      return fmt.Errorf("Failed to copy terraform directory: %v", err)
  }

  templatePath := filepath.Join(terraformModuleDir, "template.tf.tmpl")
  outputPath := filepath.Join(terraformModuleDir, "main.tf")

  err = terraform.RenderTerraform(serviceID, serviceName, repoProvider, repoOwner, runnerID, templatePath, outputPath, serviceTags)
  if err != nil {
      return fmt.Errorf("Error: %v", err)
  }
  l.Info("Terraform template processed successfully.")

  terraform.GenerateTerraformConfig(terraformServiceDir, serviceName)
  terraformCmd := "apply"
  err = terraform.RunTerraform(terraformServiceDir, terraformCmd, serviceID, serviceName)
  if err != nil {
      return fmt.Errorf("Failed to run terraform: %v", err)
  }

  consulServiceDir := config.ConsulServicesDataDir
  keyPath := fmt.Sprintf("%s/%s", consulServiceDir, serviceName)
  consulServiceData, err := consul.ConsulStoreGet(keyPath)
  if err != nil {
      return fmt.Errorf("Failed to fetch factory state: %v", err)
  }

  var (
      repoID, _  = consulServiceData["repo-id"].(string)
  )

  repoAddress := fmt.Sprintf("https://%s.com/%s/%s", repoProvider, repoOwner, serviceName)
  httpCheck := fmt.Sprintf("http://"+
  		"%s:%s/v1/repo/status?repoProvider=%s&repoID=%s&repoOwner=%s&serviceName=%s",
  		SFHost,
  		SFPort,
      repoProvider,
  		repoID,
  		repoOwner,
  		serviceName,
	)
  err = consul.RegisterRepo(serviceID, serviceName, repoAddress, httpCheck, serviceTags)
  if err != nil {
      return fmt.Errorf("Failed to register repository to Consul: %v", err)
  }

  err = stats.GenerateState()
  if err != nil {
      return fmt.Errorf("Failed to get factory info: %v", err)
  }
  mutex.Lock()
  *appliedServices = append(*appliedServices, spec.Service.Name)
  mutex.Unlock()


  return nil
}
