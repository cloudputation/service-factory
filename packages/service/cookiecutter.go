package service

import (
		"os"
		"fmt"
		"path/filepath"
		"strings"
		"text/template"

		"github.com/cloudputation/service-factory/packages/config"
		l "github.com/cloudputation/service-factory/packages/logger"
)


type RequestBody struct {
	ServiceSpecs
}

type ServiceSpecs struct {
	Service Service `hcl:"service,block"`
}

type Service struct {
  Name	string		`hcl:"name"`
  Group	string		`hcl:"group"`
  Port	string		`hcl:"port"`
  Tags	[]string	`hcl:"tags"`
  Template  			`hcl:"template,block"`
  Repo      			`hcl:"repo,block"`
  Network      		`hcl:"network,block"`
}

type Template struct {
	TemplateURL  string `hcl:"template_url"`
	TemplateName string `hcl:"template"`
}

type Repo struct {
	Provider     		string `hcl:"provider"`
	NamespaceID     string `hcl:"namespace_id"`
	RunnerID        string `hcl:"runner_id"`
	RegistryToken   string `hcl:"registry_token"`
	RepositoryOwner string `hcl:"repository_owner"`
}

type Network struct {
	AuthoritativeServer	string `hcl:"authoritative_server"`
	ClientHostname			string `hcl:"client_hostname"`
}


func GetCookieCutterFiles(templateDir string) (map[string]string, error) {
	cookieCutterFiles := make(map[string]string)

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Failed to analyze template directory %s: %v", templateDir, err)
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".ck") {
			relPath, _ := filepath.Rel(templateDir, path)
			key := filepath.Join(templateDir, relPath)
			value := filepath.Join(templateDir, strings.TrimSuffix(relPath, ".ck"))
			cookieCutterFiles[key] = value
			l.Info("Template file: %s, Rendered file: %s\n", key, value)
		}


		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Could not render repository template files: %v", err)
	}


	return cookieCutterFiles, nil
}

func RenderTemplate(templateFilename, outputFilename string, svc ServiceSpecs) error {
  // Check if template file exists
  if _, err := os.Stat(templateFilename); os.IsNotExist(err) {
      l.Error("Template file: %s does not exist: %s", templateFilename, err)
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
  err := os.MkdirAll(outputDir, 0755)
  if err != nil {
      return fmt.Errorf("Failed to create directories: %v", err)
  }

  // Parse the template with the custom function map
  tmpl, err := template.New(filepath.Base(templateFilename)).Funcs(funcMap).ParseFiles(templateFilename)
  if err != nil {
      return fmt.Errorf("Failed to parse template: %v", err)
  }

  // Create output file
  outputFile, err := os.Create(outputFilename)
  if err != nil {
      return fmt.Errorf("Failed to create output file: %v", err)
  }
  defer outputFile.Close()

  // Execute the template
  err = tmpl.Execute(outputFile, svc)
  if err != nil {
      return fmt.Errorf("Failed to execute template: %v", err)
  }

  err = os.Remove(templateFilename)
  if err != nil {
    return fmt.Errorf("Failed to remove template source file %s: %v", templateFilename, err)
  }

  gitDir := filepath.Join(config.AppConfig.DataDir, "/services/", svc.Service.Name, ".git")
  err = os.RemoveAll(gitDir)
	if err != nil {
		return fmt.Errorf("Failed to remove .git directory in %s: %v", gitDir, err)
	}

  return nil
}
