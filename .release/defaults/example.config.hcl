log_dir = "log"
data_dir = "sf-data"

server {
  port    = "8840"
  // Address of the server. It could be a local or remote address.
  // This is relevant when using the cli to communicate with the agent
  address = "127.0.0.1"
}

consul {
  consul_host   = "127.0.0.1"
  // Consul Gossip token
  consul_token  = CONSUL_TOKEN
}

terraform {
  terraform_dir = "terraform"
}

repo {
  gitlab {
    // Your Gitlab PAT
    access_token = GITLAB_TOKEN
  }
}
