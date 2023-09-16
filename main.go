package main

import (
		"os"

		"github.com/cloudputation/service-factory/packages/cli"
		"github.com/cloudputation/service-factory/packages/server"
		"github.com/cloudputation/service-factory/packages/service"
)


func main() {
  if len(os.Args) < 2 {
      cli.HelperMessage()
      return
  }

  switch os.Args[1] {
  case "config":
      cli.RunConfig()
  case "apply":
      service.ParseService("service.hcl")
  case "agent":
      server.StartServer()
	case "help":
			cli.HelperMessage()
  default:
      cli.HelperMessage()
  }
}
