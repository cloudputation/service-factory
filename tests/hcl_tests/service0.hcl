service {
  name    = "my-first-service"
  group   = "tests-apis"
  port    = "9999"
  tags    = ["SF-Managed"]

  template {
    provider_url  = "gitlab"
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

  scheduler {
    nomad {
      type            = "service"
      target_client   = "tower2"
    }
  }
}
