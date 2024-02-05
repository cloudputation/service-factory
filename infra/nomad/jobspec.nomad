job "service-factory" {
  datacenters = ["dc1"]

  constraint {
    attribute = "${attr.unique.hostname}"
    value     = "nia1"
  }

  type = "service"

  group "automation" {
    count = 1

    network {
      port "service_factory_port" {
        static = 4884
      }
    }

    task "service-factory-docker-deployment" {
      driver = "docker"

      config {
        image = "registry.gitlab.com/cloudputation/service-factory:latest"
        ports = ["service_factory_port"]

        auth {
          username = "FranksRobins"
          password = "glpat-kuuxn3oXxwBsCMkxssDY"
        }

      }

      resources {
        cpu    = 200
        memory = 256
      }

      service {
        name = "${NOMAD_JOB_NAME}"
        port = "service_factory_port"
        tags = [
          "mainframe-service"
        ]
        check {
          name     = "${NOMAD_JOB_NAME} 4884 alive"
          type     = "http"
          path     = "/health"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}
