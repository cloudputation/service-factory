service {
  name    = "test-service"
  group   = "test-apis"
  port    = "1337"
  tags    = ["SF-Managed"]

  template {
    template_url  = "gitlab.com"
    template      = "franksrobins/cookie-cutter-api"
  }

  repo {
    provider          = "gitlab"
    namespace_id      = REPOSITORY_ID
    runner_id         = RUNNER_ID
    registry_token    = GITLAB_REGISTRY_TOKEN
    repository_owner  = "franksrobins"
  }

  network {
    authoritative_server  = "127.0.0.1"
    client_hostname       = "backend1"
  }
}
