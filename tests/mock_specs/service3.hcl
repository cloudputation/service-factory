service {
  name    = "test-service13"
  group   = "test-apis"
  port    = "9700"
  tags    = ["SF-Managed"]

  template {
    template_url  = "gitlab.com"
    template      = "cloudputation/cookie-cutter-api"
  }

  repository {
    provider = "gitlab"
    config {
      namespace_id      = "69638879"
      runner_id         = "20665767"
      registry_token    = "glpat-kuuxn3oXxwBsCMkxssDY"
      repository_owner  = "franksrobins"
    }
  }

  network {
    authoritative_server  = "10.100.200.241"
    target_host           = "comm1"
  }
}
