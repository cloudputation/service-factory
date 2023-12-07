package main

import (
		"github.com/cloudputation/service-factory/packages/cli"
    "github.com/cloudputation/service-factory/packages/config"
    l "github.com/cloudputation/service-factory/packages/logger"
)


func main() {
	err := config.LoadConfiguration()
	if err != nil {
			l.Fatal("Failed to load configuration: %v", err)
	}

	err = l.InitLogger(config.AppConfig.LogDir)
  if err != nil {
      l.Error("Error initializing logger: %v", err)
  }
  defer l.CloseLogger()

  rootCmd := cli.SetupRootCommand()
  if err := rootCmd.Execute(); err != nil {
      l.Fatal("Error executing command: %v", err)
  }
}
