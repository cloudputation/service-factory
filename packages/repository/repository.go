package repository

import (
    "os"
    "fmt"
    "sync"
    "os/exec"
    "path/filepath"
)


var gitCloneMutex sync.Mutex

func DownloadRepoToDatastore(repoDatastore, repoDatastoreParent, templateURL, templateName string) error {
  // Check if repo already exists in the datastore
  gitCloneMutex.Lock()
  if _, err := os.Stat(repoDatastore); os.IsNotExist(err) {
    cmd := exec.Command("mkdir", "-p", repoDatastoreParent)
    err = cmd.Run()
  	if err != nil {
      return fmt.Errorf("Failed to create namespace repository directory in datastore: %v", err)
  	}
    // If repo dont exist, clone it
    err = GitClone(templateURL, templateName, repoDatastore)
    if err != nil {
        gitCloneMutex.Unlock()
        return fmt.Errorf("Failed to clone repository to datastore: %v", err)
    }
  }
  gitCloneMutex.Unlock()


  return nil
}

func GitClone(repoProvider string, template string, cookieCutterDestination string) error {
	cookieCutterGitUrl := fmt.Sprintf("https://%s/%s.git", repoProvider, template)

	fmt.Println("Cloning from", cookieCutterGitUrl)

	cmd := exec.Command("git", "clone", cookieCutterGitUrl, cookieCutterDestination)
  err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to clone cookie cutter %s: %v", cookieCutterGitUrl, err)
	}

  gitDirectory := filepath.Join(cookieCutterDestination, ".git")
  cmd = exec.Command("rm", "-fr", gitDirectory)
  err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to delete .git directory in %s: %v", cookieCutterDestination, err)
	}


	return nil
}
