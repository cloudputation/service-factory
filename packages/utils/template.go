package utils

import (
    "os"
    "log"
    "fmt"
    "text/template"
    "path/filepath"

    "github.com/cloudputation/service-factory/packages/config"
    "github.com/cloudputation/service-factory/packages/service"
)


func RenderTemplate(templateFilename, outputFilename string, svc service.ServiceSpecs) error {
    // Check if template file exists
    if _, err := os.Stat(templateFilename); os.IsNotExist(err) {
        log.Printf("Template file does not exist: %s", templateFilename)
        return fmt.Errorf("Template file does not exist: %s", templateFilename)
    }

    // Define custom function map
    funcMap := template.FuncMap{
        "sub": func(a, b int) int {
            return a - b
        },
    }

    // Ensure the directory exists
    outputDir := filepath.Dir(outputFilename)
    err := os.MkdirAll(outputDir, 0755) // The 0755 parameter sets the directory's permissions
    if err != nil {
        log.Printf("Error creating directories: %v", err)
        return err
    }

    // Parse the template with the custom function map
    tmpl, err := template.New(filepath.Base(templateFilename)).Funcs(funcMap).ParseFiles(templateFilename)
    if err != nil {
        log.Printf("Error parsing template: %v", err)
        return err
    }

    // Create output file
    outputFile, err := os.Create(outputFilename)
    if err != nil {
        log.Printf("Error creating output file: %v", err)
        return err
    }
    defer outputFile.Close()

    // Execute the template
    err = tmpl.Execute(outputFile, svc)
    if err != nil {
        log.Printf("Error executing template: %v", err)
        return err
    }

    err = os.Remove(templateFilename)
    if err != nil {
      fmt.Println("Error:", err)
    }

    gitDir := filepath.Join(config.AppConfig.DataDir, "/services/", svc.Service.Name, ".git")

    err = os.RemoveAll(gitDir)

  	if err != nil {
  		fmt.Println("Error:", err)
  	}

    return nil
}
