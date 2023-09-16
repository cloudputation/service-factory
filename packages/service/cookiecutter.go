package service

import (
  "fmt"
  "io/ioutil"
  "strings"
  "path/filepath"
)


func GetCookieCutterFiles(service_dir string) map[string]string {
    basePaths := map[string]string{
        "service.py.tpl":                   "service.py",
        "Dockerfile.tpl":                   "Dockerfile",
        ".gitlab-ci.yml.tpl":               ".gitlab-ci.yml",
        "API_VERSION":                      "API_VERSION",
        "requirements.txt":                 "requirements.txt",
        "nomad/nomad.tf.tpl":               "nomad/nomad.tf",
        "nomad/jobspec/service.nomad.tpl":  "nomad/jobspec/service.nomad",
    }

    cookieCutterFiles := make(map[string]string)

    for key, value := range basePaths {
        cookieCutterFiles[fmt.Sprintf("cookie-cutter-api/%s", key)] = fmt.Sprintf("%s/%s", service_dir, value)
    }

    return cookieCutterFiles
}
// func GetCookieCutterFiles(serviceDir string, templateDir string) map[string]string {
// 	// Create a map to store the filenames
// 	cookieCutterFiles := make(map[string]string)
//
// 	// Read the directory contents
// 	files, err := ioutil.ReadDir(templateDir)
// 	if err != nil {
// 		fmt.Println("Error reading directory:", err)
// 		return cookieCutterFiles
// 	}
//
// 	// Loop over the directory contents to find .cc files and render them
// 	for _, file := range files {
// 		if !file.IsDir() && strings.HasSuffix(file.Name(), ".cc") {
// 			key := filepath.Join(templateDir, file.Name())
// 			value := filepath.Join(serviceDir, strings.TrimSuffix(file.Name(), ".cc"))
// 			cookieCutterFiles[key] = value
// 		}
// 	}
//
// 	return cookieCutterFiles
// }
