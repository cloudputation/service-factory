package main

import (
	"os"
	"text/template"
)

type ServiceConfig struct {
	Meta map[string]string `hcl:"meta"`
}

const tmplStr = `
Service Meta Data:
{{- range $key, $value := .Meta }}
  {{ $key }}: {{ $value }}
{{- end }}
`

func main() {
	// Sample data
	serviceConfig := ServiceConfig{
		Meta: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	// Parse the template
	tmpl, err := template.New("service").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	// Execute the template
	if err := tmpl.Execute(os.Stdout, serviceConfig); err != nil {
		panic(err)
	}
}
