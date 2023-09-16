package service

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "path/filepath"

    "github.com/cloudputation/service-factory/packages/utils"
    "github.com/cloudputation/service-factory/packages/config"
)


func ParseService(serviceFilePath string) {
    jsonOutput, err := utils.HCLtoJSON(serviceFilePath)
    if err != nil {
        log.Fatalf("Error during conversion: %v", err)
    }

    err = SendJSONPayload(jsonOutput)
    if err != nil {
        log.Fatalf("Failed to send JSON payload: %v", err)
    }

    fmt.Println("Payload sent successfully!")
}

func SendJSONPayload(jsonData string) error {
    url := "http://10.100.200.243:48840/create"
    req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
    }

    return nil
}
