package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

type OpenAPI struct {
	Components Components `json:"components"`
}

type Components struct {
	Schemas map[string]Schema `json:"schemas"`
}

type Schema struct {
	Title      string              `json:"title"`
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Title string `json:"title"`
	Type  string `json:"type"`
}

type RequiredField struct {
	Name string
	Type string
}

func main() {
	// Read JSON file
	jsonFile, err := os.Open("api-struct/api-struct.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Unmarshal JSON
	var openAPI OpenAPI
	json.Unmarshal(byteValue, &openAPI)

	mediaSchema := openAPI.Components.Schemas["Media"]

	requiredFields := make([]RequiredField, len(mediaSchema.Required))

	for i, req := range mediaSchema.Required {
		requiredFields[i] = RequiredField{Name: req, Type: mediaSchema.Properties[req].Type}
	}

	// Open and parse the template file
	tmpl, err := template.ParseFiles("templates/provider.tmpl")
	if err != nil {
		panic(err)
	}

	// Execute template
	err = tmpl.Execute(os.Stdout, struct{ RequiredFields []RequiredField }{RequiredFields: requiredFields})
	if err != nil {
		panic(err)
	}
}
