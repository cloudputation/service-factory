package service

import (
	"os"
	"path/filepath"
	"strings"
)

func GetCookieCutterFiles(templateDir string) (map[string]string, error) {
	cookieCutterFiles := make(map[string]string)

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".ck") {
			relPath, _ := filepath.Rel(templateDir, path)
			key := filepath.Join(templateDir, relPath)
			value := filepath.Join(templateDir, strings.TrimSuffix(relPath, ".ck"))
			cookieCutterFiles[key] = value
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return cookieCutterFiles, nil
}
