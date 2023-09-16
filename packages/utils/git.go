package utils

import (
    "os/exec"

    "github.com/cloudputation/service-factory/packages/service"
)


func GitClone(serviceSpecs service.ServiceSpecs) error {
	repoProvider := serviceSpecs.Repo.Provider
	template := serviceSpecs.Repo.TemplateName
	cookieCutterGitUrl := fmt.Sprintf("https://%s/%s.git", repoProvider, template)

	cmd := exec.Command("git", "clone", cookieCutterGitUrl, template)
	return cmd.Run()
}
