package utils

import (
    "os"
    "text/template"

    "github.com/cloudputation/service-factory/packages/service"
)


func RenderTemplate(templateFilename, outputFilename string, svc service.ServiceSpecs) error {
	tmpl, err := template.ParseFiles(templateFilename)
	if err != nil {
		return err
	}
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	err = tmpl.Execute(outputFile, svc)
	if err != nil {
		return err
	}
	return nil
}
