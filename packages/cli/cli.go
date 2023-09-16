package cli

import (
		"fmt"
		"io/ioutil"
		"log"
		"path/filepath"

    "github.com/cloudputation/service-factory/packages/config"
)


func RunConfig() {
	// Read HCL content and convert to JSON
	hclPath := filepath.Join(".", "config.hcl")
	hclContent, err := ioutil.ReadFile(hclPath)
	if err != nil {
		log.Fatalf("Failed to read HCL file: %v", err)
	}

	jsonOutput, err := config.HCLtoJSON(string(hclContent))
	if err != nil {
		log.Fatalf("Error during conversion: %v", err)
	}

	fmt.Println(jsonOutput)
}
