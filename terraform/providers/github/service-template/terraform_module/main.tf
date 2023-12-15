// Create new repository
resource "github_repository" "project" {
  name          = var.repo_name
  visibility    = "private"
  has_issues    = true
  has_projects  = true
  has_wiki      = false
  // Additional GitHub-specific configurations
}

// Upload project files
resource "github_repository_file" "files" {
  for_each = local.filtered_files

  repository    = github_repository.project.name
  file          = "${each.key}"
  branch        = "master"
  content       = file("${var.data_dir}/services/${var.repo_name}/repo/${each.key}")
  commit_author   = var.author_name
  commit_email    = var.author_email
  commit_message  = var.commit_message

  // Depends on repository being created
  depends_on = [
    github_repository.project
  ]
}

// Upload CI file
resource "github_repository_file" "ci_file" {
  repository      = github_repository.project.name
  file            = ".github/workflows/ci.yml"
  branch          = "main"
  content         = file("${var.data_dir}/services/${var.repo_name}/repo/.github/workflows/action.yml")
  commit_author   = var.author_name
  commit_email    = var.author_email
  commit_message  = var.commit_message

  // Depends on project files
  depends_on = [
    github_repository_file.files
  ]
}

// Register service manifest in Consul
resource "consul_keys" "manifest_key" {
  key {
    path   = "${var.consul_path}/${var.repo_name}"
    value  = jsonencode({
        "service_id": var.service_id,
        "repository_provider": "github",
        "repository_owner": var.repository_owner,
        "repository_id": github_repository.project.id
    })
  }

  // Depends on repository setup
  depends_on = [
    github_repository_file.ci_file
  ]
}
