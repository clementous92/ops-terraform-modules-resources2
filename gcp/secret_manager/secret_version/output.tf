output "output" {
  value = google_secret_manager_secret_version.secret_version
}

output "project_id" {
  value = var.project_id
}
