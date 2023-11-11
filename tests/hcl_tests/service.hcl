service {
  name    = "test-service1"
  group   = "test-apis"
  port    = "9977"
  tags    = ["SF-Managed"]

  template {
    template_url  = "gitlab.com"
    template      = "franksrobins/cookie-cutter-api"
  }

  repo {
    provider_url      = "gitlab.com"
    namespace_id      = "69638879"
    token             = "glpat-HuEekH9zXTi8DbWkyLzo"
    runner_id         = "20665767"
    registry_token    = "glpat-kuuxn3oXxwBsCMkxssDY"
    repository_owner  = "franksrobins"
  }

  network {
    authoritative_server  = "10.100.200.241"
    client_hostname       = "tower2"
  }

  // network {
  //   parameter = value <- find this option first
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
