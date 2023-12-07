log_dir = "sf.log"
data_dir = "sf-data"

server {
  port    = "48840"
  address = "192.168.0.10"
}

consul {
  consul_host   = "192.168.0.10"
  consul_token  = "pLggHasySyfQhkTBWuX30jJ2eM0UJ3Rs6rK8cFNTt/o="
}

terraform {
  terraform_dir = "terraform"
}

repo {
  gitlab {
    access_token = "glpat-HuEekH9zXTi8DbWkyLzo"
  }
}
