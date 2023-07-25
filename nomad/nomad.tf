terraform {
  required_providers {
    nomad = {
      source = "hashicorp/nomad"
      version = "1.4.20"
    }
  }
}

provider "nomad" {
    address = "${var.nomad_server_address}"
}

variable "nomad_server_address" {
  type = "string"
}

variable "jobspec_path" {
  type = "string"
}


resource "nomad_job" "app" {
  purge_on_destroy = true
  jobspec = file("${var.jobspec_path}")
}
