package terraform

import (
    "fmt"
    "os"
    "path/filepath"
    "os/exec"
    "strings"

    "github.com/cloudputation/service-factory/packages/config"
    l "github.com/cloudputation/service-factory/packages/logger"
    "github.com/cloudputation/service-factory/packages/service"
)


func RunTerraform(
  terraformDir,
  terraformCmd,
  ServiceID,
  serviceName string,
  repoVars    map[string]string,
  ) error {
  args := []string{
    terraformCmd,
    "-auto-approve",
    "--terragrunt-non-interactive",
    "--terragrunt-working-dir", terraformDir,
    "-var", fmt.Sprintf("consul_path=%s", config.ConsulServicesDataDir),
    "-var", fmt.Sprintf("data_dir=%s/%s", config.RootDir, config.AppConfig.DataDir),
    "-var", fmt.Sprintf("service_id=%s", ServiceID),
    "-var", fmt.Sprintf("repo_name=%s", serviceName),
  }

  for key, value := range repoVars {
      if key != "namespace_id" && key != "provider" && key != "registry_token" {
        args = append(args, "-var", fmt.Sprintf("%s=%s", key, value))
      }
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

func GenerateTerraformConfig(configPath, serviceName , repoProvider string) error {
  var repoToken string

  consulTerraformDir  := config.ConsulFactoryDataDir
  consulToken         := config.AppConfig.Consul.ConsulToken
  consulAddress       := config.AppConfig.Consul.ConsulHost

  if repoProvider == "github" {
      repoToken = config.AppConfig.Repository.Github.AccessToken
  } else if repoProvider == "gitlab" {
      repoToken = config.AppConfig.Repository.Gitlab.AccessToken
  }


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
      "api_token: \"%s\"\n"+
      "service_name: \"%s\"\n",
      consulTerraformDir, consulToken, consulAddress, repoToken, serviceName,
  )

  _, err = file.WriteString(content)
  if err != nil {
    return fmt.Errorf("Error writing to file:", err)
  }
  l.Info("Terraform config written successfully to: %s", configPath)


  return nil
}

func ImportTerraform() error {
  dataDir       := config.AppConfig.DataDir
  repoURL       := config.SFRepoURL
  cloneDest     := dataDir + "/tmp/sf-repo"
  terraformPath := filepath.Join(cloneDest, "terraform")
  importDest    := dataDir + "/terraform"

  if _, err := os.Stat(importDest); os.IsNotExist(err) {
      l.Info("Terraform repository not present at: %s. Importing..", importDest)
      if err := exec.Command("git", "clone", repoURL, cloneDest).Run(); err != nil {
          return fmt.Errorf("Failed to clone Terraform directory: %v", err)
      }
      defer os.RemoveAll(cloneDest)
      if err := exec.Command("cp", "-r", terraformPath, importDest).Run(); err != nil {
          return fmt.Errorf("Failed to copy Terraform %s directory to Data directory %s: %v", terraformPath, importDest, err)
      }
  }
  l.Info("Terraform repository present at: %s", importDest)


  return nil
}

func GetRepoVars(repo service.Repository) map[string]string {
  repoVars := make(map[string]string)

  // Check if GitLab configuration is provided
  if repo.Provider == "gitlab" {
      gitlab := repo.RepoConfig
      repoVars["provider"] = "gitlab"

      if gitlab.NamespaceID != nil {
          repoVars["namespace_id"] = *gitlab.NamespaceID
      }
      if gitlab.RunnerID != nil {
          repoVars["runner_id"] = *gitlab.RunnerID
      }
      repoVars["registry_token"] = gitlab.RegistryToken
      repoVars["repository_owner"] = gitlab.RepositoryOwner
  }

  // Check if GitHub configuration is provided
  if repo.Provider == "github" {
      github := repo.RepoConfig
      repoVars["provider"] = "github"
      repoVars["registry_token"] = github.RegistryToken
      repoVars["repository_owner"] = github.RepositoryOwner
  }

  return repoVars
}
