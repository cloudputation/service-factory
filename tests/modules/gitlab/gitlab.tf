// Define repo variables
variable "data_dir" {
  description = "Service Factory data directory"
  type        = string
}

variable "repo_name" {
  description = "The name of the GitLab repository"
  type        = string
}

variable "runner_id" {
  description = "The ID of the gitlab runner"
  type        = string
  sensitive   = true
}

variable "gitlab_token" {
  description = "GitLab personal access token"
  type        = string
  sensitive   = true
}

variable "commit_message" {
  description = "The commit message for the new file"
  type        = string
  default     = "Managed by Service Factory"
}

variable "author_email" {
  description = "The email of the author of the commit"
  type        = string
  default     = "dev@service-factory"
}

variable "author_name" {
  description = "The name of the author of the commit"
  type        = string
  default     = "fishstock dev"
}

// Get project files
locals {
  all_files = fileset("${var.data_dir}/services/${var.repo_name}/repo", "/**/*")
  filtered_files = { for file in local.all_files : file => file if file != ".gitlab-ci.yml" }
}

// Configure provider
terraform {
  required_providers {
    gitlab = {
      source = "gitlabhq/gitlab"
      version = "15.7.1"
    }
    time = {
      source = "hashicorp/time"
      version = "0.9.1"
    }
  }
}

provider "gitlab" {
  base_url = "https://gitlab.com/api/v4/"
  token = var.gitlab_token
}

// Create new repo
resource "gitlab_project" "new_project" {
  name = var.repo_name
  visibility_level = "private"
  shared_runners_enabled = false
}

data "gitlab_project" "project_info" {
  path_with_namespace = "FranksRobins/${var.repo_name}"
  depends_on = [
    gitlab_project.new_project
  ]
}

// Upload project files
resource "gitlab_repository_file" "upload_files" {
  for_each = local.filtered_files

  project        = gitlab_project.new_project.id
  file_path      = "${each.key}"
  branch         = "master"
  content        = file("${var.data_dir}/services/${var.repo_name}/repo/${each.key}")
  author_email   = var.author_email
  author_name    = var.author_name
  commit_message = var.commit_message

  depends_on = [
    gitlab_project.new_project
  ]
}

// Upload CI file
resource "gitlab_repository_file" "upload_ci_file" {
  project        = gitlab_project.new_project.id
  file_path      = ".gitlab-ci.yml"
  branch         = "master"
  content        = file("${var.data_dir}/services/${var.repo_name}/repo/.gitlab-ci.yml")
  author_email   = var.author_email
  author_name    = var.author_name
  commit_message = var.commit_message

  depends_on = [
    gitlab_repository_file.upload_files
  ]
}


// Connect repo to pipeline
resource "gitlab_project_runner_enablement" "enable_runner" {
  project   = data.gitlab_project.project_info.id
  runner_id = var.runner_id

  depends_on = [
    gitlab_repository_file.upload_files
  ]
}
