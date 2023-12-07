package cli

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"

  	"github.com/hashicorp/hcl/v2/gohcl"
  	"github.com/hashicorp/hcl/v2/hclparse"

    "github.com/cloudputation/service-factory/packages/config"
    l "github.com/cloudputation/service-factory/packages/logger"
)


func GetFactoryStatus() error {
  apiEndpoint := fmt.Sprintf(
      "http://%s:%s/system/status",
      config.AppConfig.Server.ServerAddress,
      config.AppConfig.Server.ServerPort,
  )

  resp, err := http.Get(apiEndpoint)
  if err != nil {
      return fmt.Errorf("Failed to make the GET request:", err)
  }
  defer resp.Body.Close()

  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      return fmt.Errorf("Failed to read response body:", err)
  }

  l.Info(string(bodyBytes))


  return nil
}

func GetServiceStatus(serviceName string) error {
  apiEndpoint := fmt.Sprintf(
      "http://%s:%s/service/status?serviceName=%s",
      config.AppConfig.Server.ServerAddress,
      config.AppConfig.Server.ServerPort,
      serviceName,
  )

  resp, err := http.Get(apiEndpoint)
  if err != nil {
      return fmt.Errorf("Failed to make the GET request:", err)
  }
  defer resp.Body.Close()

  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      return fmt.Errorf("Failed to read response body:", err)
  }

  l.Info(string(bodyBytes))


  return nil
}

func CheckConfig() error {
		err := config.LoadConfiguration()
		if err != nil {
				return fmt.Errorf("Failed to load configuration: %v", err)
		}

		hclPath := config.GetConfigPath()
		jsonOutput, err := AppHCLtoJSON(hclPath)
		if err != nil {
				return fmt.Errorf("Error during conversion: %v", err)
		}

		fmt.Println(jsonOutput)


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

    diags = gohcl.DecodeBody(file.Body, nil, &config.AppConfig)
    if diags.HasErrors() {
        return "", fmt.Errorf("Failed to decode HCL: %s", diags.Error())
    }

    jsonOutput, err := json.MarshalIndent(config.AppConfig, "", "  ")
    if err != nil {
        return "", fmt.Errorf("Failed to marshal to JSON: %v", err)
    }


    return string(jsonOutput), nil
}
