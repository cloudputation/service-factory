locals {
  config = yamldecode(file("config.yaml"))
}

remote_state {
  backend = "consul"
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite"
  }
  config = {
    access_token = "${local.config.consul_access_token}"
    path = "${local.config.consul_factory_data_dir}/terraform/${local.config.service_name}"
    address = "${local.config.consul_address}:8500"
  }
}

generate "providers" {
  path      = "providers.tf"
  if_exists = "overwrite"
  contents  = <<-EOC
    variable "gitlab_api_url" {
      description = "Gitlab API base URL"
      type        = string
      default     = "https://gitlab.com/api/v4/"
    }

    terraform {
      required_providers {
        gitlab = {
          source = "gitlabhq/gitlab"
          version = "16.6.0"
        }
        consul = {
          source  = "hashicorp/consul"
          version = "2.20.0"
        }
      }
    }

    provider "gitlab" {
      token     = "${local.config.gitlab_api_token}"
      base_url  = var.gitlab_api_url
    }

    provider "consul" {
      address    = "${local.config.consul_address}:8500"
      token      = "${local.config.consul_access_token}"
    }
    EOC
}

terraform {
  source = "terraform_module"
}
