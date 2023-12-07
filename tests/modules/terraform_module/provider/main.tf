// Create new repo
resource "gitlab_project" "new_project" {
  name = var.repo_name
  visibility_level = "private"
  shared_runners_enabled = false
}

data "gitlab_project" "project_info" {
  path_with_namespace = "${var.repo_owner}/${var.repo_name}"
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
