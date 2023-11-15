package utils

import (
  "fmt"
  "os/exec"

  "github.com/cloudputation/service-factory/packages/config"
)


func InitTerraform(terraformDir string) error {
    cmd := exec.Command(
        "terraform",
        "-chdir="+terraformDir,
        "init",
    )
    return cmd.Run()
}

func RunTerraform(
	terraformDir string,
	terraformCmd string,
  serviceName  string,
  runnerID     string,
  repoToken    string,
) error {
	cmd := exec.Command(
		"terraform",
		"-chdir="+terraformDir,
		terraformCmd, "-auto-approve",
		"-var", "data_dir="+config.AppConfig.DataDir,
		"-var", "repo_name="+serviceName,
		"-var", "runner_id="+runnerID,
		"-var", "gitlab_token="+repoToken,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Terraform failed with the following output:\n%s\n", string(output))
		return err
	}
	fmt.Printf("Terraform succeeded with the following output:\n%s\n", string(output))
	return nil
}
