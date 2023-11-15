package main

import (
	"os"
	"log"

	"github.com/cloudputation/service-factory/packages/bootstrap"
	"github.com/cloudputation/service-factory/packages/cli"
	"github.com/cloudputation/service-factory/packages/config"
	"github.com/cloudputation/service-factory/packages/server"
	"github.com/cloudputation/service-factory/packages/service"
)


func main() {
  if len(os.Args) < 2 {
      cli.HelperMessage()
      return
  }
	err := config.LoadConfiguration()
	if err != nil {
			log.Fatalf("Failed to load config: %v", err)
	}

  switch os.Args[1] {
  case "config":
      cli.CheckConfig()
  case "apply":
      service.ParsePayload("service.hcl")
  case "agent":
			err = bootstrap.BootstrapFactory()
			if err != nil {
					log.Fatalf("Failed to bootstrap Service Factory: %v", err)
			}
      server.StartServer()
	case "help":
			cli.HelperMessage()
  default:
      cli.HelperMessage()
  }
}
