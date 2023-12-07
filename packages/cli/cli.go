package cli

import (
		"github.com/spf13/cobra"

		"github.com/cloudputation/service-factory/packages/bootstrap"
    l "github.com/cloudputation/service-factory/packages/logger"
		"github.com/cloudputation/service-factory/packages/server"
)


func SetupRootCommand() *cobra.Command {
	var serviceName  string
	var serviceFiles []string
	var serviceNames []string

	var rootCmd = &cobra.Command{
		Use:   "factory",
		Short: "Service Factory is a comprehensive tool for managing service onboarding within your organization",
	}
	rootCmd.CompletionOptions.HiddenDefaultCmd = true


	var cmdConfig = &cobra.Command{
		Use:		"config",
		Short:	"Check the configuration",
		Args:		cobra.ExactArgs(0),
		Run:	func(cmd *cobra.Command, args []string) {
			err := CheckConfig()
			if err != nil {
				l.Error("Failed to parse configuration file: %v", err)
			}
		},
	}

	var cmdApply = &cobra.Command{
		Use:		"apply",
		Short:	"Apply a service configuration",
		Run:		func(cmd *cobra.Command, args []string) {
			err := ApplyServiceSpecs(serviceFiles)
			if err != nil {
				l.Error("Failed to parse service file: %v", err)
			}
		},
	}
	cmdApply.Flags().StringArrayVarP(&serviceFiles, "service-file", "f", []string{}, "Path to service file")

	var cmdDestroy = &cobra.Command{
		Use:		"destroy",
		Short:	"Decommission a service",
		Run:	func(cmd *cobra.Command, args []string) {
			err := DestroyService(serviceNames)
			if err != nil {
				l.Error("Failed to decommission services: %v", err)
			}
		},
	}
	cmdDestroy.Flags().StringArrayVarP(&serviceNames, "service", "s", []string{}, "Name of the service to destroy")

	var cmdAgent = &cobra.Command{
		Use:		"agent",
		Short:	"Start the service agent",
		Run:	func(cmd *cobra.Command, args []string) {
			err := bootstrap.BootstrapFactory()
      if err != nil {
          l.Fatal("Failed to bootstrap the factory: %v", err)
      }
			server.StartServer()
		},
	}

	var cmdService = &cobra.Command{
		Use:   "service",
		Short: "Commands related to service operations",
	}

	var cmdServiceStatus = &cobra.Command{
		Use:		"status [serviceName]",
		Short:	"Check status of a service",
  	Args:		cobra.ExactArgs(1),
		Run:	func(cmd *cobra.Command, args []string) {
			serviceName = args[0]
			err := GetServiceStatus(serviceName)
			if err != nil {
				l.Error("Failed to check service status: %v", err)
			}
		},
	}
	cmdService.AddCommand(cmdServiceStatus)

	var cmdSystem = &cobra.Command{
		Use:   "system",
		Short: "Commands related to system operations",
	}

	var cmdSystemStatus = &cobra.Command{
		Use:		"status",
		Short:	"Check the system status",
  	Args:		cobra.ExactArgs(0),
		Run:	func(cmd *cobra.Command, args []string) {
			err := GetFactoryStatus()
			if err != nil {
				l.Error("Failed to check server status %v", err)
			}
		},
	}
	cmdSystem.AddCommand(cmdSystemStatus)



	var commands = []*cobra.Command{
			cmdConfig,
			cmdApply,
			cmdDestroy,
			cmdAgent,
			cmdService,
			cmdSystem,
	}

	for _, cmd := range commands {
		rootCmd.AddCommand(cmd)
	}

  return rootCmd

}
