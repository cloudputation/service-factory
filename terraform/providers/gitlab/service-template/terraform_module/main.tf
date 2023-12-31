// Create new repo
resource "gitlab_project" "project" {
  name = var.repo_name
  visibility_level = "private"
  shared_runners_enabled = false
}

data "gitlab_project" "project_info" {
  path_with_namespace = "${var.repository_owner}/${var.repo_name}"

  depends_on = [
    gitlab_project.project
  ]
}

// Upload project files
resource "gitlab_repository_file" "files" {
  for_each = local.filtered_files

  project        = gitlab_project.project.id
  file_path      = "${each.key}"
  branch         = "master"
  content        = file("${var.data_dir}/services/${var.repo_name}/repo/${each.key}")
  encoding       = "text"
  author_email   = var.author_email
  author_name    = var.author_name
  commit_message = var.commit_message

  depends_on = [
    gitlab_project.project
  ]
}

// Upload CI file
resource "gitlab_repository_file" "ci_file" {
  project        = gitlab_project.project.id
  file_path      = ".gitlab-ci.yml"
  branch         = "master"
  content        = file("${var.data_dir}/services/${var.repo_name}/repo/.gitlab-ci.yml")
  encoding       = "text"
  author_email   = var.author_email
  author_name    = var.author_name
  commit_message = var.commit_message

  depends_on = [
    gitlab_repository_file.files
  ]
}

// Connect repo to pipeline
resource "gitlab_project_runner_enablement" "ci_runner" {
  project   = data.gitlab_project.project_info.id
  runner_id = var.runner_id

  depends_on = [
    gitlab_repository_file.files
  ]
}

// Register service manifest
resource "consul_keys" "manifest_key" {
  key {
    path   = "${var.consul_path}/${var.repo_name}"
    value  = jsonencode({
        "service_id": var.service_id,
        "repository_provider": "gitlab",
        "repository_owner": var.repository_owner,
        "repository_id": data.gitlab_project.project_info.id,
        "runner_id": var.runner_id
    })
  }

  depends_on = [
    gitlab_project_runner_enablement.ci_runner
  ]
}
