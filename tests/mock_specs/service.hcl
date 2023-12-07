service {
  name    = "test-service11"
  group   = "test-apis"
  port    = "9900"
  tags    = ["SF-Managed"]

  template {
    template_url  = "gitlab.com"
    template      = "franksrobins/cookie-cutter-api"
  }

  repo {
    provider          = "gitlab"
    namespace_id      = "69638879"
    runner_id         = "20665767"
    registry_token    = "glpat-kuuxn3oXxwBsCMkxssDY"
    repository_owner  = "franksrobins"
  }

  network {
    authoritative_server  = "10.100.200.241"
    client_hostname       = "comm1"
  }

  // network {
  //   parameter = value <- find this option first repo_service_name = ?
  //   consul {
  //     consul_host = "10.100.200.241"
  //   }
  // }
  //
  // scheduler {
  //   parameter = value <- find this option first
  //   nomad {
  //     nomad_host  = "10.100.200.241"
  //   }
  // }
}
