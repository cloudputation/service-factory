package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type OpenAPI struct {
	Paths      map[string]map[string]PathDetails `json:"paths"`
	Components struct {
		Schemas map[string]Schema `json:"schemas"`
	} `json:"components"`
}

type PathDetails struct {
	Summary     string `json:"summary"`
	OperationID string `json:"operationId"`
	RequestBody struct {
		Content map[string]struct {
			Schema struct {
				Ref string `json:"$ref"`
			} `json:"schema"`
		} `json:"content"`
		Required bool `json:"required"`
	} `json:"requestBody"`
}

type Schema struct {
	Properties map[string]map[string]string `json:"properties"`
}

func main() {
	jsonFile, err := os.Open("api.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var openAPI OpenAPI

	json.Unmarshal(byteValue, &openAPI)

	for path, methods := range openAPI.Paths {
		for method, details := range methods {
			if details.RequestBody.Required {
				for _, schema := range details.RequestBody.Content {
					schemaName := strings.TrimPrefix(schema.Schema.Ref, "#/components/schemas/")
					actualSchema := openAPI.Components.Schemas[schemaName]
					var properties []string
					for propName := range actualSchema.Properties {
						properties = append(properties, propName)
					}
					fmt.Printf("Path: %s, Method: %s, Summary: %s\n", path, method, details.Summary)
					fmt.Printf("Schema Properties: %v\n\n", properties)
				}
			}
		}
	}
}
