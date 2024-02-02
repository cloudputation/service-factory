package main

import (
		"github.com/cloudputation/service-factory/packages/cli"
		"github.com/cloudputation/service-factory/packages/config"
		l "github.com/cloudputation/service-factory/packages/logger"
		"github.com/cloudputation/service-factory/packages/stats"
)

func main() {
	// Load main configuration file
	err := config.LoadConfiguration()
	if err != nil {
			l.Fatal("Failed to load configuration: %v", err)
	}

	// Initialize server metrics
  err = stats.InitMetrics()
  if err != nil {
			l.Fatal("Failed to initialize metrics service: %v", err)
  }

	// Initialize logging system
	err = l.InitLogger(config.RootDir + "/" + config.AppConfig.LogDir)
  if err != nil {
			l.Fatal("Error initializing logs: %v", err)
  }
  defer l.CloseLogger()

	// Run CLI
  rootCmd := cli.SetupRootCommand()
  if err := rootCmd.Execute(); err != nil {
			l.Fatal("Error executing command: %v", err)
  }
}
