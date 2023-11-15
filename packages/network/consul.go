package network

import (
	"log"

	"github.com/hashicorp/consul/api"
)


var ConsulClient *api.Client

func InitConsul(consulAddress string) error {
	// Initialize Consul client
	consulConfig := api.DefaultConfig()
  consulPort := ":8500"
	consulConfig.Address = consulAddress+consulPort

	var err error
	ConsulClient, err = api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Consul client: %v", err)
	}

	kv := ConsulClient.KV()

	// Try to read a key-value pair from Consul
	pair, _, err := kv.Get("service-factory/data/stats", nil)
	if err != nil {
		log.Fatalf("Failed to initiate Consul connection: %v", err)
	}

	// If key-value pair doesn't exist, create one
	if pair == nil {
		writeOptions := &api.WriteOptions{}
		p := &api.KVPair{Key: "service-factory/data/status", Value: []byte("initialized.")}
		_, err = kv.Put(p, writeOptions)
		if err != nil {
			log.Fatalf("Failed to initialize data store on Consul: %v", err)
		}
		log.Println("Data store initialized successfully.")
	} else {
		log.Println("Data store is initialized.")
	}

	return nil
}

// func AddEntryToJSON(client *api.Client, key string, newEntryKey string, newEntryValue interface{}) error {
// 	jsonData, err := RetrieveJSONData(client, key)
// 	if err != nil {
// 		return err
// 	}
//
// 	jsonData[newEntryKey] = newEntryValue
//
// 	newJSONValue, err := json.Marshal(jsonData)
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = ConsulClient.KV().Put(&api.KVPair{Key: key, Value: newJSONValue}, nil)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func RemoveEntryFromJSON(client *api.Client, key string, entryKeyToRemove string) error {
// 	jsonData, err := RetrieveJSONData(client, key)
// 	if err != nil {
// 		return err
// 	}
//
// 	delete(jsonData, entryKeyToRemove)
//
// 	newJSONValue, err := json.Marshal(jsonData)
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = ConsulClient.KV().Put(&api.KVPair{Key: key, Value: newJSONValue}, nil)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
