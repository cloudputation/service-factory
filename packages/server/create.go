package server

import (
  "net/http"
  "encoding/json"
  "fmt"
  "log"
  "path/filepath"

  "github.com/cloudputation/service-factory/packages/utils"
  "github.com/cloudputation/service-factory/packages/config"
  "github.com/cloudputation/service-factory/packages/stats"
  "github.com/cloudputation/service-factory/packages/service"
)


type CreateResponseBody struct {
	Message  string `json:"message"`
}


func CreateHandlerWrapper(w http.ResponseWriter, r *http.Request) {
  CreateHandler(w, r)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    return
  }

  var body service.RequestBody
  err := json.NewDecoder(r.Body).Decode(&body)
  if err != nil {
    http.Error(w, "Error reading request body", http.StatusBadRequest)
    return
  }


  dataDir := config.AppConfig.DataDir
  if dataDir == "" {
      http.Error(w, "Data directory not configured", http.StatusInternalServerError)
      return
  }


  serviceBaseDir := filepath.Join(dataDir, "services", body.Service.Name)
  serviceRepoDir := filepath.Join(serviceBaseDir, "repo")
  serviceTerraformDir := filepath.Join(serviceBaseDir, "terraform")

  err = utils.CreateDir(serviceRepoDir)
  if err != nil {
    http.Error(w, "Failed to create service directory: "+err.Error(), http.StatusInternalServerError)
    return
  }

  err = utils.CreateDir(serviceTerraformDir)
  if err != nil {
    http.Error(w, "Failed to create service directory: "+err.Error(), http.StatusInternalServerError)
    return
  }


  repoProvider := body.Service.Template.TemplateURL
	template := body.Service.Template.TemplateName

  err = utils.GitClone(repoProvider, template, serviceRepoDir)
  if err != nil {
    http.Error(w, "Failed to clone the repository: "+err.Error(), http.StatusInternalServerError)
    return
  }

  cookieCutterFiles, err := service.GetCookieCutterFiles(serviceRepoDir)
  if err != nil {
      http.Error(w, "Failed to get CookieCutter files: "+err.Error(), http.StatusInternalServerError)
      return
  }

  for tplFile, outputFile := range cookieCutterFiles {
    fmt.Println("\nthe template file:", tplFile)
    fmt.Println("the rendered file:", outputFile)
      err = utils.RenderTemplate(tplFile, outputFile, body.ServiceSpecs)
      if err != nil {
          http.Error(w, "Failed to render the template: "+err.Error(), http.StatusInternalServerError)
          return
      }
  }

  fmt.Println("using terrafomr dir:", config.AppConfig.Terraform.TerraformDir)


  srcBasePath := config.AppConfig.Terraform.TerraformDir
	dstBasePath := serviceTerraformDir

  srcPath := filepath.Join(srcBasePath, ".terraform")
	dstPath := filepath.Join(dstBasePath, ".terraform")

	if err := utils.CopyDir(srcPath, dstPath); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Directory copied successfully")
	}


  filesToCopy := []string{".terraform.lock.hcl", "terraform.tfstate", "gitlab.tf"}

	for _, fileName := range filesToCopy {
    srcPath := filepath.Join(srcBasePath, fileName)
    dstPath := filepath.Join(dstBasePath, fileName)

    if err := utils.CopyFile(srcPath, dstPath); err != nil {
			fmt.Printf("Error copying file %s: %v\n", fileName, err)
		} else {
			fmt.Printf("Successfully copied %s\n", fileName)
		}
	}


  terraformCmd := "apply"
  err = utils.RunTerraform(serviceTerraformDir, terraformCmd, body.Service.Name, body.Service.Repo.RunnerID, body.Service.Repo.Token)
  if err != nil {
    http.Error(w, "Failed to run terraform: "+err.Error(), http.StatusInternalServerError)
    return
  }

  err = stats.GenerateState()
	if err != nil {
		log.Fatalf("Failed to get factory info: %v", err)
	}

  err = service.GenerateManifest(body.Service.Name, body.Service.Repo.RunnerID, body.Service.Repo.Token)
  if err != nil {
    log.Fatalf("Failed to generate manifest: %v", err)
  }



  json.NewEncoder(w).Encode(CreateResponseBody{Message: "Operation completed successfully"})
}
