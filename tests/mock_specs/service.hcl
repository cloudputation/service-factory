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
}
