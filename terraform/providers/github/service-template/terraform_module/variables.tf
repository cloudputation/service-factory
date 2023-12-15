// Define repo variables
variable "consul_path" {
  description = "Service Factory consul path"
  type        = string
}

variable "data_dir" {
  description = "Service Factory data directory"
  type        = string
}

variable "repo_name" {
  description = "The name of the repository"
  type        = string
}

variable "service_id" {
  description = "The ID of the resident service"
  type        = string
}

variable "repository_owner" {
  description = "The owner of the repository"
  type        = string
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
  filtered_files = { for file in local.all_files : file => file if file != ".github/workflows/ci.yml" }
}
