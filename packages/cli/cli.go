package cli

import (
		"encoding/json"
		"fmt"
		"io/ioutil"
		"os"
		"path/filepath"

		"github.com/hashicorp/hcl/v2/gohcl"
		"github.com/hashicorp/hcl/v2/hclparse"

		"github.com/cloudputation/service-factory/packages/config"
)


func CheckConfig() error {  // Added 'error' return type
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		// Load the configuration
		err := config.LoadConfiguration()
		if err != nil {
			return fmt.Errorf("Failed to load configuration: %v", err)
		}

		// Convert the HCL file to JSON for debugging
		hclPath := filepath.Join(".", "config.hcl")
		jsonOutput, err := AppHCLtoJSON(hclPath)
		if err != nil {
			return fmt.Errorf("Error during conversion: %v", err)
		}

		// Output the JSON and loaded configuration
		fmt.Println("Debug Mode: Converted HCL to JSON:")
		fmt.Println(jsonOutput)

	}
	return nil
}

func AppHCLtoJSON(HCLFilePath string) (string, error) {
    hclContent, err := ioutil.ReadFile(HCLFilePath)
    if err != nil {
        return "", fmt.Errorf("Failed to read HCL file: %v", err)
    }

    parser := hclparse.NewParser()
    file, diags := parser.ParseHCL(hclContent, HCLFilePath)
    if diags.HasErrors() {
        return "", fmt.Errorf("Failed to parse HCL: %s", diags.Error())
    }

    // Populate the package-level Config variable
    diags = gohcl.DecodeBody(file.Body, nil, &config.AppConfig) // Using package-level Config
    if diags.HasErrors() {
        return "", fmt.Errorf("Failed to decode HCL: %s", diags.Error())
    }

    jsonOutput, err := json.MarshalIndent(config.AppConfig, "", "  ")
    if err != nil {
        return "", fmt.Errorf("Failed to marshal to JSON: %v", err)
    }

    return string(jsonOutput), nil
}
