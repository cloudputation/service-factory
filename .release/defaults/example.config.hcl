log_dir = "sf.log"
data_dir = "sf-data"

server {
  port    = "8840"
  address = "127.0.0.1"
}

consul {
  consul_host   = "127.0.0.1"
  consul_token  = CONSUL_TOKEN
}

terraform {
  terraform_dir = "terraform"
}

repo {
  gitlab {
    access_token = GITLAB_TOKEN
  }
}
