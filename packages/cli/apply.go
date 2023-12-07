package cli

import (
    "fmt"
    "bytes"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/hashicorp/hcl/v2/gohcl"
    "github.com/hashicorp/hcl/v2/hclparse"

    "github.com/cloudputation/service-factory/packages/config"
    l "github.com/cloudputation/service-factory/packages/logger"
    "github.com/cloudputation/service-factory/packages/service"
)


type ServicesWrapper struct {
    Services []service.ServiceSpecs `json:"services"`
}

var serviceSpecs *service.ServiceSpecs


func ApplyServiceSpecs(serviceFiles []string) error {
    var servicesWrapper ServicesWrapper
    var errors []string

    for _, serviceFile := range serviceFiles {
        jsonOutput, err := ServiceHCLtoJSON(serviceFile)
        if err != nil {
            errorMsg := fmt.Sprintf("Error during conversion of %s: %v", serviceFile, err)
            errors = append(errors, errorMsg)
            continue
        }

        var serviceSpec service.ServiceSpecs
        err = json.Unmarshal([]byte(jsonOutput), &serviceSpec)
        if err != nil {
            errorMsg := fmt.Sprintf("Error unmarshalling JSON for %s: %v", serviceFile, err)
            errors = append(errors, errorMsg)
            continue
        }

        servicesWrapper.Services = append(servicesWrapper.Services, serviceSpec)
    }

    wrapperJSON, err := json.Marshal(servicesWrapper)
    if err != nil {
        errorMsg := fmt.Sprintf("Error marshalling wrapper JSON: %v", err)
        errors = append(errors, errorMsg)
    } else {
        err = SendJSONPayload(string(wrapperJSON))
        if err != nil {
            errorMsg := fmt.Sprintf("Failed to send JSON payload: %v", err)
            errors = append(errors, errorMsg)
        } else {
            l.Info("Payloads sent successfully!")
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("encountered errors: %v", errors)
    }
    return nil
}


func ServiceHCLtoJSON(HCLFilePath string) (string, error) {
    hclContent, err := ioutil.ReadFile(HCLFilePath)
    if err != nil {
        return "", fmt.Errorf("Failed to read HCL file: %w", err)
    }

    parser := hclparse.NewParser()
    file, diags := parser.ParseHCL(hclContent, HCLFilePath)
    if diags.HasErrors() {
        return "", fmt.Errorf("Failed to parse HCL: %s", diags.Error())
    }

    var localServiceSpecs service.ServiceSpecs
    diags = gohcl.DecodeBody(file.Body, nil, &localServiceSpecs)
    if diags.HasErrors() {
        return "", fmt.Errorf("Failed to decode HCL: %s", diags.Error())
    }

    jsonOutput, err := json.MarshalIndent(localServiceSpecs, "", "  ")
    if err != nil {
        return "", fmt.Errorf("Failed to marshal to JSON: %v", err)
    }

    return string(jsonOutput), nil
}

func SendJSONPayload(jsonData string) error {
    serverAddress := config.AppConfig.Server.ServerAddress
    serverPort := config.AppConfig.Server.ServerPort
    serverEndpoint := fmt.Sprintf("http://%s:%s/apply", serverAddress, serverPort)

    req, err := http.NewRequest("POST", serverEndpoint, bytes.NewBuffer([]byte(jsonData)))
    if err != nil {
        return fmt.Errorf("Failed to create new HTTP request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("Failed to execute HTTP request: %v", err)
    }
    defer resp.Body.Close()

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("Failed to read response body: %v", err)
    }

    response := string(bodyBytes)
    l.Info("Response from server: %s", response)

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
    }

    return nil
}
