package cli

import (
  "fmt"
  "io/ioutil"
  "net/http"

  "github.com/cloudputation/service-factory/packages/config"
)


func DeleteService(serviceName string) error {
    apiEndpoint := fmt.Sprintf(
        "http://%s:%s/delete?serviceName=%s",
        config.AppConfig.Server.ServerAddress,
        config.AppConfig.Server.ServerPort,
        serviceName,
    )

    resp, err := http.Post(apiEndpoint, "application/json", nil)
    if err != nil {
        fmt.Println("Failed to make the POST request:", err)
        return err
    }
    defer resp.Body.Close()

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Failed to read response body:", err)
        return err
    }

    fmt.Println("Response:", string(bodyBytes))
    return nil
}
