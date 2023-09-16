package utils

import (
    "os/exec"

    "github.com/cloudputation/service-factory/packages/service"
)


func RunTerraform(
    terraformDir string,
    serviceSpecs service.ServiceSpecs,
    terraformCmd string,
) error {
    cmd := exec.Command(
        "terraform",
        "-chdir="+terraformDir,
        terraformCmd, "-auto-approve",
        "-var", "repo_name="+serviceSpecs.Scheduler.Nomad.ServiceName,
        "-var", "gitlab_token="+serviceSpecs.Repo.Token,
        "-var", "runner_id="+serviceSpecs.Repo.RunnerID,
    )
    return cmd.Run()
}
