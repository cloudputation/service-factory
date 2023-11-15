package config

import (
  "os"

  "github.com/hashicorp/hcl/v2/gohcl"
  "github.com/hashicorp/hcl/v2/hclparse"
)


type Configuration struct {
  DataDir string        `hcl:"data_dir"`
  Server  Server        `hcl:"server,block"`
  Consul  Consul        `hcl:"consul,block"`
  Terraform  Terraform  `hcl:"terraform,block"`
}

type Server struct {
  ServerPort    string `hcl:"port"`
  ServerAddress string `hcl:"address"`

}

type DataDir struct {
  DataDir string `hcl:data_dir`
}

type Consul struct {
  ConsulHost string `hcl:"consul_host"`

}

type Terraform struct {
  TerraformDir string `hcl:"terraform_dir"`

}

var AppConfig Configuration




func LoadConfiguration() error {
  configPath := os.Getenv("SF_CONFIG_FILE_PATH")
  if configPath == "" {
      configPath = "/etc/service-factory/config.hcl" // Default path
  }

  // Read the HCL file
  data, err := os.ReadFile(configPath)
  if err != nil {
      return err
  }

  // Parse the HCL file
  parser := hclparse.NewParser()
  file, diags := parser.ParseHCL(data, "config.hcl")
  if diags.HasErrors() {
      return diags
  }

  // Decode the HCL file into your Config struct
  diags = gohcl.DecodeBody(file.Body, nil, &AppConfig)
  if diags.HasErrors() {
      return diags
  }

  return nil
}
