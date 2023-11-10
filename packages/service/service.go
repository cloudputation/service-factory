package service

import (
    "bytes"
    "fmt"
    "net/http"
)


func SendJSONPayload(jsonData string) error {
    url := "http://10.100.200.242:48840/create"
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
