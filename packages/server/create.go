package server

import (
  "net/http"
  "encoding/json"
  "fmt"

  "github.com/cloudputation/service-factory/packages/utils"
  "github.com/cloudputation/service-factory/packages/config"
  "github.com/cloudputation/service-factory/packages/service"
)


type ResponseBody struct {
	Message  string `json:"message"`
}

var terraformCommand string


func CreateHandler(w http.ResponseWriter, r *http.Request, serviceSpecs service.ServiceSpecs) {
  if r.Method != http.MethodPost {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    return
  }

  var body service.RequestBody
  err := json.NewDecoder(r.Body).Decode(&body)
  if err != nil {
    http.Error(w, "Error reading request body",
      http.StatusBadRequest)
    return
  }

  var serviceDir = body.ServiceName
  templateDir := fmt.Sprintf("%s/template", serviceDir)

  cookieCutterFiles := service.GetCookieCutterFiles(serviceDir, templateDir)


    err = utils.CreateServiceDir(body.ServiceSpecs)
  if err != nil {
    http.Error(w, "Failed to create service directory: "+err.Error(),
      http.StatusInternalServerError)
    return
  }

  err = utils.GitClone()
  if err != nil {
    http.Error(w, "Failed to clone the repository: "+err.Error(),
      http.StatusInternalServerError)
    return
  }

  for tplFile, outputFile := range cookieCutterFiles {
    err = utils.RenderTemplate(tplFile, outputFile, body.ServiceSpecs)
    if err != nil {
      http.Error(w, "Failed to render the template: "+err.Error(),
        http.StatusInternalServerError)
      return
    }
  }

  var terraformCommand = "apply"
  err = utils.RunTerraform(config.terraformDir, body.ServiceSpecs, terraformCommand)
  if err != nil {
    http.Error(w, "Failed to run terraform: "+err.Error(),
      http.StatusInternalServerError)
    return
  }

  err = utils.CleanServiceDir()
  if err != nil {
    http.Error(w, "Failed to clean service directory: "+err.Error(),
      http.StatusInternalServerError)
    return
  }

  json.NewEncoder(w).Encode(ResponseBody{Message: "Operation completed successfully"})
}
