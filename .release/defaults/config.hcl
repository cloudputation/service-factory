log_dir = "log"
data_dir = "sf-data"

server {
  port    = "48840"
  address = "10.100.200.210"
}

consul {
  consul_host   = "10.100.200.210"
  consul_token  = "pLggHasySyfQhkTBWuX30jJ2eM0UJ3Rs6rK8cFNTt/o="
}

repository {
  gitlab {
    access_token = "glpat-HuEekH9zXTi8DbWkyLzo"
  }
  github {
    access_token = "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"
  }
}
