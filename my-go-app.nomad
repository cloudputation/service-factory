job "test-go-app" {
  datacenters = ["dc1"]

  type = "service"

  group "group1" {
    count = 1

    network {
      port "go_app" {
        static = 48840
      }
    }

    task "task1" {
      driver = "docker"

      config {
        image = "localhost:5000/my-go-app:latest"
        ports = ["go_app"]

      }

      resources {
        cpu    = 500
        memory = 256
      }

      service {
        name = "${NOMAD_JOB_NAME}"
        port = "go_app"
        tags = ["test-app"]

        check {
          type     = "http"
          path     = "/health"
          interval = "10s"
          timeout  = "2s"
        }
      }

      env {
        NOMAD_JOB_NAME = "${NOMAD_JOB_NAME}"
      }
    }
  }
}
