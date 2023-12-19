package config

import (
    "os"
    "fmt"

    "github.com/spf13/viper"
    "github.com/hashicorp/hcl/v2/gohcl"
    "github.com/hashicorp/hcl/v2/hclparse"
)


type Configuration struct {
    LogDir      string      `hcl:"log_dir"`
    DataDir     string      `hcl:"data_dir"`
    Server      Server      `hcl:"server,block"`
    Consul      Consul      `hcl:"consul,block"`
    Repository  Repository  `hcl:"repository,block"`
}

type Server struct {
    ServerPort    string `hcl:"port"`
    ServerAddress string `hcl:"address"`
}

type Consul struct {
    ConsulHost  string `hcl:"consul_host"`
    ConsulToken string `hcl:"consul_token"`
}

type Repository struct {
    Gitlab  *Gitlab `hcl:"gitlab,block"`
    Github  *Github `hcl:"github,block"`
}

type Gitlab struct {
    AccessToken string `hcl:"access_token"`
}

type Github struct {
    AccessToken string `hcl:"access_token"`
}


var AppConfig Configuration
var ConfigPath string
var RootDir string
const MaxWorkers = 10

var (
    DatastoreDir                = "sf-data/repositories"
    ConsulDataDir               = "sf::Data"
    ConsulTerraformDataDir      = ConsulDataDir + "/factory/terraform"
    ConsulFactoryDataDir        = ConsulDataDir + "/factory"
    ConsulServicesDataDir       = ConsulDataDir + "/services"
    ConsulServiceSummaryDataDir = ConsulServicesDataDir + "/summary"
    SFRepoURL                   = "https://github.com/cloudputation/service-factory.git"
)


func LoadConfiguration() error {
  viper.SetDefault("ConfigPath", "/etc/service-factory/config.hcl")
  viper.BindEnv("ConfigPath", "SF_CONFIG_FILE_PATH")

  ConfigPath = viper.GetString("ConfigPath")

  var err error
  RootDir, err = os.Getwd()
  if err != nil {
      return fmt.Errorf("Failed to get service root directory: %v", err)
  }

  // Read the HCL file
  data, err := os.ReadFile(ConfigPath)
  if err != nil {
      return fmt.Errorf("Failed to read configuration file: %v", err)
  }

  // Parse the HCL file
  parser := hclparse.NewParser()
  file, diags := parser.ParseHCL(data, ConfigPath)
  if diags.HasErrors() {
      return fmt.Errorf("Failed to parse configuration: %v", diags)
  }

  // Populate the Config struct
  diags = gohcl.DecodeBody(file.Body, nil, &AppConfig)
  if diags.HasErrors() {
      return fmt.Errorf("Failed to apply configuration: %v", diags)
  }

  return nil
}

func GetConfigPath() string {
  return ConfigPath
}
