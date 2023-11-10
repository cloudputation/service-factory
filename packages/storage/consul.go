package storage

import (
	"log"
	"fmt"
	"encoding/json"

	"github.com/hashicorp/consul/api"
)


func ConsulStorePut(jsonData string, keyPath string) error {
	// Initialize Consul client
	consulConfig := api.DefaultConfig()
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return err
	}

	// Get KV API client
	kv := client.KV()

	pair, _, err := kv.Get(keyPath, nil)
	if err != nil {
		return err
	}

	// If the key-value pair doesn't exist, create one
	if pair == nil {
		log.Println("No Factory state found. Generating.")
  } else {
		log.Println("Refreshing factory state.")
	}
	writeOptions := &api.WriteOptions{}
	p := &api.KVPair{Key: keyPath, Value: []byte(jsonData)}
	_, err = kv.Put(p, writeOptions)
	if err != nil {
		return err
	}
	log.Println("Factory state created successfully!")


	return nil
}

func ConsulStoreGet(ConsulClient *api.Client, key string) (map[string]interface{}, error) {
	kvPair, _, err := ConsulClient.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kvPair == nil {
		return nil, fmt.Errorf("Key not found: %s", key)
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(kvPair.Value, &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func ConsulStoreDelete(ConsulClient *api.Client, keyPath string) error {
	_, err := ConsulClient.KV().Delete(keyPath, nil)
	if err != nil {
		return fmt.Errorf("Failed to delete key: %s, error: %v", keyPath, err)
	}
	return nil
}
