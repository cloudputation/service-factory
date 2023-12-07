package bootstrap

import (
    "fmt"
    "os/exec"

    "github.com/cloudputation/service-factory/packages/config"
    "github.com/cloudputation/service-factory/packages/consul"
    l "github.com/cloudputation/service-factory/packages/logger"
    "github.com/cloudputation/service-factory/packages/stats"
)


func BootstrapFactory() error {
  l.Info("Starting Service Factory agent.. Bootstrapping factory.")
  dataDir := config.AppConfig.DataDir
  l.Info("Loaded configuration file: %s", config.ConfigPath)

  serviceDataDir := fmt.Sprintf("%s/services", dataDir)
  cmd := exec.Command("mkdir", "-p", serviceDataDir)
  err := cmd.Run()
  if err != nil {
      return fmt.Errorf("Failed to create directory: %v", err)
  }

  err = consul.InitConsul()
  if err != nil {
      return fmt.Errorf("Could not initialize Consul: %v", err)
  }

  err = consul.BootstrapConsul()
  if err != nil {
      return fmt.Errorf("Could not bootstrap factory on Consul: %v", err)
  }

  l.Info("Refreshing factory state.")
  err = stats.GenerateState()
  if err != nil {
      return fmt.Errorf("Failed to get factory info: %v", err)
  }
  l.Info("Factory state created successfully!")


  return nil
}
