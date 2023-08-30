package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

type Config struct {
	Repo  Repo   `hcl:"repo,block"`
	Nomad Nomad  `hcl:"nomad,block"`
}

type Repo struct {
	TemplateURL           string `hcl:"template_url"`
	NamespaceID           string `hcl:"namespace_id"`
	Token                 string `hcl:"token"`
	RunnerID              string `hcl:"runner_id"`
	ServiceRegistryToken  string `hcl:"service_registry_token"`
	ServiceRepositoryOwner string `hcl:"service_repository_owner"`
}

type Nomad struct {
	ServerAddress string `hcl:"server_address"`
}

func main() {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile("config.hcl")
	if diags.HasErrors() {
		log.Fatalf("Failed to parse HCL: %s", diags.Error())
	}

	var config Config
	diags = gohcl.DecodeBody(file.Body, nil, &config)
	if diags.HasErrors() {
		log.Fatalf("Failed to decode HCL: %s", diags.Error())
	}

	jsonOutput, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal to JSON: %v", err)
	}

	fmt.Println(string(jsonOutput))
}
