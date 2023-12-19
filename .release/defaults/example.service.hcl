service {
  name    = "test-service"
  group   = "test-apis"
  port    = "1337"
  tags    = ["SF-Managed"]

  template {
    template_url  = "gitlab.com"
    template      = "cloudputation/cookie-cutter-api"
  }

  repo {
    provider          = "gitlab"
    namespace_id      = REPOSITORY_ID
    runner_id         = RUNNER_ID
    registry_token    = GITLAB_REGISTRY_TOKEN
    repository_owner  = "cloudputation"
  }

  network {
    // The server instance of your scheduler system
    // For example, your Nomad server
    authoritative_server  = "127.0.0.1"
    // The target host to deploy the service.
    // This can be a hostname or any constraint value
    target_host           = "backend1"
  }
}
