package consul

import (
		"fmt"
		"encoding/json"

		"github.com/hashicorp/consul/api"
		"github.com/cloudputation/service-factory/packages/config"
	  l "github.com/cloudputation/service-factory/packages/logger"
)


type FactoryStatus struct {
	Status string `json:"factory_status"`
}

var ConsulClient *api.Client
var err error
var statusDir = config.ConsulFactoryDataDir + "/status"


func InitConsul() error {
	consulConfig := api.DefaultConfig()
  consulHost := config.AppConfig.Consul.ConsulHost
	consulPort := ":8500"
	consulConfig.Address = consulHost + consulPort

	ConsulClient, err = api.NewClient(consulConfig)
	if err != nil {
			return fmt.Errorf("Failed to initialize Consul client: %v", err)
	}
	l.Info("Consul client initialized successfully.")


	return nil
}

func BootstrapConsul() error {
	kv := ConsulClient.KV()

	// Try to read a key-value pair from Consul
	l.Info("Checking if data store is initialized.")
	pair, _, err := kv.Get(statusDir, nil)
	if err != nil {
			return fmt.Errorf("Failed to initiate Consul connection: %v", err)
	}

	if pair == nil {
		l.Info("Consul data store is not initialized. Initializing..")
		// Create an instance of FactoryStatus with the status set to "initialized"
		status := FactoryStatus{Status: "initialized"}

		statusJSON, err := json.Marshal(status)
		if err != nil {
				return fmt.Errorf("Failed to marshal JSON: %v", err)
		}

		writeOptions := &api.WriteOptions{}
		p := &api.KVPair{Key: statusDir, Value: statusJSON}
		_, err = kv.Put(p, writeOptions)
		if err != nil {
				return fmt.Errorf("Failed to initialize data store on Consul: %v", err)
		}
		l.Info("Data store initialized successfully.")
	} else {
		l.Info("Data store is already initialized.")
	}


	return nil
}

func ConsulStoreGet(key string) (map[string]interface{}, error) {
  kv := ConsulClient.KV()

  kvPair, _, err := kv.Get(key, nil)
	if err != nil {
			return nil, fmt.Errorf("Failed to query key on Consul: %v", err)
	}
	if kvPair == nil {
			return nil, fmt.Errorf("Key not found: %s", key)
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(kvPair.Value, &jsonData)
	if err != nil {
			return nil, fmt.Errorf("Failed to parse Consul key %s: %v", key, err)
	}


	return jsonData, nil
}

func ConsulStorePut(jsonData, keyPath string) error {
	kv := ConsulClient.KV()

	writeOptions := &api.WriteOptions{}
	p := &api.KVPair{Key: keyPath, Value: []byte(jsonData)}
	_, err = kv.Put(p, writeOptions)
	if err != nil {
			return fmt.Errorf("Failed to upload key: %s, error: %v", keyPath, err)
	}


	return nil
}

func ConsulStoreDelete(keyPath string) error {
  kv := ConsulClient.KV()

	_, err := kv.Delete(keyPath, nil)
	if err != nil {
			return fmt.Errorf("Failed to delete key: %s, error: %v", keyPath, err)
	}


	return nil
}

func RegisterRepo(serviceID, serviceName, repoAddress, httpCheck string, serviceTags []string) error {
	serviceRegistration := &api.AgentServiceRegistration{
	    ID:      serviceID,
	    Name:    "service-repository",
	    Tags:    serviceTags,
	    Address: repoAddress,
	    Meta: map[string]string{
	        "resident-service": serviceName,
	    },
	    Check: &api.AgentServiceCheck{
					Name:			"Service Repository Alive",
	        HTTP:     httpCheck,
	        Interval: "5m",
	        Timeout:  "10s",
	    },
	}

  client := ConsulClient
  err := client.Agent().ServiceRegister(serviceRegistration)
  if err != nil {
      return fmt.Errorf("Error registering service: %v", err)
  }
  l.Info("Service registered successfully on Consul")


	return nil
}

func DeregisterRepo(serviceID string) error {
  client := ConsulClient
  err := client.Agent().ServiceDeregister(serviceID)
  if err != nil {
      return fmt.Errorf("Error registering service: %v", err)
  }
  l.Info("Service deregistered successfully from Consul")


	return nil
}
