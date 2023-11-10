package bootstrap

import (
  "log"

  "github.com/cloudputation/service-factory/packages/config"
  "github.com/cloudputation/service-factory/packages/stats"
  "github.com/cloudputation/service-factory/packages/network"
)


func BootstrapFactory() error {
	err := network.InitConsul(config.AppConfig.Consul.ConsulHost)
  if err != nil {
    return err
  }

  err = stats.GenerateState()
	if err != nil {
		log.Fatalf("Failed to get factory info: %v", err)
	}



  return nil
}
