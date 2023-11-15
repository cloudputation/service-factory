package utils

import (
  "fmt"
  "os/exec"
  "path/filepath"
)


func GitClone(repoProvider string, template string, cookieCutterDestination string) error {
	cookieCutterGitUrl := fmt.Sprintf("https://%s/%s.git", repoProvider, template)

	fmt.Println("Cloning from", cookieCutterGitUrl)

	cmd := exec.Command("git", "clone", cookieCutterGitUrl, cookieCutterDestination)
  err := cmd.Run()
	if err != nil {
		return err
	}

  gitDirectory := filepath.Join(cookieCutterDestination, ".git")
  err = DeleteDir(gitDirectory)
  if err != nil {
		return err
	}

	return nil
}
