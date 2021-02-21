resource "google_secret_manager_secret_version" "secret_version" {
  provider = google-beta

  secret      = var.secret_id
  secret_data = var.secret_data
}
