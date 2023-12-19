package main

import (
		"github.com/cloudputation/service-factory/packages/cli"
		"github.com/cloudputation/service-factory/packages/config"
		l "github.com/cloudputation/service-factory/packages/logger"
		"github.com/cloudputation/service-factory/packages/stats"
)


func main() {
	err := config.LoadConfiguration()
	if err != nil {
			l.Fatal("Failed to load configuration: %v", err)
	}

  err = stats.InitMetrics()
  if err != nil {
			l.Fatal("Failed to initialize metrics service: %v", err)
  }

	err = l.InitLogger(config.RootDir + "/" + config.AppConfig.LogDir)
  if err != nil {
			l.Fatal("Error initializing logs: %v", err)
  }
  defer l.CloseLogger()

  rootCmd := cli.SetupRootCommand()
  if err := rootCmd.Execute(); err != nil {
			l.Fatal("Error executing command: %v", err)
  }
}
