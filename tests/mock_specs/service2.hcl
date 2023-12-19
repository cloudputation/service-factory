service {
  name    = "test-service12"
  group   = "test-apis"
  port    = "9800"
  tags    = ["SF-Managed"]

  template {
    template_url  = "github.com"
    template      = "cloudputation/cookie-cutter-api"
  }

  repository {
    provider = "github"
    config {
      registry_token    = "glpat-kuuxn3oXxwBsCMkxssDY"
      repository_owner  = "franksrobins"
    }
  }

  network {
    authoritative_server  = "10.100.200.241"
    target_host           = "comm1"
  }
}
