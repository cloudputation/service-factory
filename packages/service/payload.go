package service

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"

  "github.com/hashicorp/hcl/v2/gohcl"
  "github.com/hashicorp/hcl/v2/hclparse"
)


func ParsePayload(serviceFilePath string) {
    jsonOutput, err := ServiceHCLtoJSON(serviceFilePath)
    if err != nil {
        log.Fatalf("Error during conversion: %v", err)
    }

    err = SendJSONPayload(jsonOutput)
    if err != nil {
        log.Fatalf("Failed to send JSON payload: %v", err)
    }

    fmt.Println("Payload sent successfully!")
}

func ServiceHCLtoJSON(HCLFilePath string) (string, error) {
    hclContent, err := ioutil.ReadFile(HCLFilePath)
    if err != nil {
        return "", fmt.Errorf("Failed to read HCL file: %v", err)
    }

    parser := hclparse.NewParser()
    file, diags := parser.ParseHCL(hclContent, HCLFilePath)
    if diags.HasErrors() {
        return "", fmt.Errorf("Failed to parse HCL: %s", diags.Error())
    }

    var serviceSpecs ServiceSpecs
    diags = gohcl.DecodeBody(file.Body, nil, &serviceSpecs)
    if diags.HasErrors() {
        return "", fmt.Errorf("Failed to decode HCL: %s", diags.Error())
    }

    jsonOutput, err := json.MarshalIndent(serviceSpecs, "", "  ")
    if err != nil {
        return "", fmt.Errorf("Failed to marshal to JSON: %v", err)
    }

    return string(jsonOutput), nil
}
