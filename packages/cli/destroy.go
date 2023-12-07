package cli

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/cloudputation/service-factory/packages/config"
    l "github.com/cloudputation/service-factory/packages/logger"
)

type ServicesToDestroy struct {
    ServiceNames []string `json:"service-names"`
}

func DestroyService(serviceNames []string) error {
    services := ServicesToDestroy{ServiceNames: serviceNames}
    jsonData, err := json.Marshal(services)
    if err != nil {
        return fmt.Errorf("Error marshalling JSON:", err)
    }

    apiEndpoint := fmt.Sprintf(
        "http://%s:%s/destroy",
        config.AppConfig.Server.ServerAddress,
        config.AppConfig.Server.ServerPort,
    )

    resp, err := http.Post(apiEndpoint, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("Failed to make the POST request:", err)
    }
    defer resp.Body.Close()

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("Failed to read response body:", err)
    }

    response := string(bodyBytes)
    l.Info("Server Response: %s", response)


    return nil
}
