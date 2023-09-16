package utils

import (
    "os/exec"

    "github.com/cloudputation/service-factory/packages/service"
)


func CreateServiceDir(serviceSpecs service.ServiceSpecs) error {
    cmd := exec.Command(
      "mkdir",
      "-p",
      serviceSpecs.Scheduler.Nomad.ServiceName+"/nomad/jobspec"
    )
    return cmd.Run()
}
var Template = service.ServiceSpecs.Repo.TemplateName

func CleanServiceDir(serviceSpecs service.ServiceSpecs) error {
    cmd := exec.Command(
      "rm",
      "-fr",
      serviceSpecs.Scheduler.Nomad.ServiceName,
      serviceSpecs.Repo.TemplateName
    )
    return cmd.Run()
}
