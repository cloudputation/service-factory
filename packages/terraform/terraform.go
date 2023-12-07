package terraform

import (
    "fmt"
    "os"
    "text/template"
    "path/filepath"
    "os/exec"
    "strings"

    "github.com/cloudputation/service-factory/packages/config"
    l "github.com/cloudputation/service-factory/packages/logger"
)


type TemplateData struct {
  SFHost            string
  SFPort            string
  DataDir           string
  ServiceID         string
  SerialServiceID   string
  ServiceName       string
  ServiceTags       []string
  RepoProvider      string
  RepoOwner         string
  RunnerID          string
  ConsulServiceData string
}


func RunTerraform(
  terraformDir,
  terraformCmd,
  ServiceID,
  serviceName   string,
  ) error {
  args := []string{
    terraformCmd,
    "-auto-approve",
    "--terragrunt-non-interactive",
    "--terragrunt-working-dir", terraformDir,
    "-var", fmt.Sprintf("data_dir=%s/%s", config.RootDir, config.AppConfig.DataDir),
    "-var", fmt.Sprintf("service_id=%s", ServiceID),
    "-var", fmt.Sprintf("repo_name=%s", serviceName),
  }

  commandString := "terragrunt " + strings.Join(args, " ")
  l.Info("Running terraform command: %s\n", commandString)

  cmd := exec.Command("terragrunt", args...)

  output, err := cmd.CombinedOutput()
  if err != nil {
      return fmt.Errorf("Terraform failed with the following output:\n%s\n%v", string(output), err)
  }
  fmt.Printf("Terraform succeeded with the following output:\n%s\n", string(output))


  return nil
}

func GenerateTerraformConfig(configPath, serviceName string) error {
  consulTerraformDir  := config.ConsulFactoryDataDir
  consulToken         := config.AppConfig.Consul.ConsulToken
  consulAddress       := config.AppConfig.Consul.ConsulHost
  gitlabToken         := config.AppConfig.Repo.Gitlab.AccessToken

  configPath = fmt.Sprintf("%s/config.yaml", configPath)

  file, err := os.Create(configPath)
  if err != nil {
    return fmt.Errorf("Error creating file:", err)
  }
  defer file.Close()

  // Prepare config file
  content := fmt.Sprintf(
      "consul_factory_data_dir: \"%s\"\n"+
      "consul_access_token: \"%s\"\n"+
      "consul_address: \"%s\"\n"+
      "gitlab_api_token: \"%s\"\n"+
      "service_name: \"%s\"\n",
      consulTerraformDir, consulToken, consulAddress, gitlabToken, serviceName,
  )

  _, err = file.WriteString(content)
  if err != nil {
    return fmt.Errorf("Error writing to file:", err)
  }
  l.Info("Terraform config written successfully to: %s", configPath)


  return nil
}

func RenderTerraform(
  serviceID,
  serviceName,
  repoProvider,
  repoOwner,
  runnerID,
  templatePath,
  outputPath  string,
  ServiceTags []string,
  ) error {
  serialServiceID := strings.ReplaceAll(serviceID, "-", "_")

  // Define custom function map
  funcMap := template.FuncMap{
      "sub": func(a, b int) int {
          return a - b
      },
  }

  dataDir := filepath.Join(config.RootDir, config.AppConfig.DataDir)
  data := TemplateData{
      SFHost:             config.AppConfig.Server.ServerAddress,
      SFPort:             config.AppConfig.Server.ServerPort,
      DataDir:            dataDir,
      ServiceID:          serviceID,
      SerialServiceID:    serialServiceID,
      ServiceName:        serviceName,
      ServiceTags:        ServiceTags,
      RepoProvider:       repoProvider,
      RepoOwner:          repoOwner,
      RunnerID:           runnerID,
      ConsulServiceData:  config.ConsulServicesDataDir,
  }

  // Render the template
  err := RenderTemplate(templatePath, outputPath, data, funcMap)
    if err != nil {
      return fmt.Errorf("Failed to render template %s: %v", templatePath, err)
  }

  // Delete the source template
  return os.Remove(templatePath)
}

func RenderTemplate(templatePath, outputPath string, data TemplateData, funcMap template.FuncMap) error {
    tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
    if err != nil {
        return fmt.Errorf("Failed to parse template files: %v", err)
    }

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("Failed to create destination files", err)
    }
    defer outputFile.Close()

    return tmpl.Execute(outputFile, data)
}
