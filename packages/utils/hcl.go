package utils

import (
    "encoding/json"
    "fmt"
    "io/ioutil"

    "github.com/hashicorp/hcl/v2/gohcl"
    "github.com/hashicorp/hcl/v2/hclparse"
)


func HCLtoJSON(configFilePath string) (string, error) {
    hclContent, err := ioutil.ReadFile(configFilePath)
    if err != nil {
        return "", fmt.Errorf("Failed to read HCL file: %v", err)
    }

    parser := hclparse.NewParser()
    file, diags := parser.ParseHCL(hclContent, configFilePath)
    if diags.HasErrors() {
        return "", fmt.Errorf("Failed to parse HCL: %s", diags.Error())
    }

    var config Config
    diags = gohcl.DecodeBody(file.Body, nil, &config)
    if diags.HasErrors() {
        return "", fmt.Errorf("Failed to decode HCL: %s", diags.Error())
    }

    jsonOutput, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return "", fmt.Errorf("Failed to marshal to JSON: %v", err)
    }

    return string(jsonOutput), nil
}
