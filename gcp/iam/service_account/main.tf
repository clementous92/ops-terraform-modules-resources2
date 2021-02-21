resource "google_service_account" "service_account" {
  project    = var.project_id
  account_id = var.account_id

  display_name = var.display_name
  description  = var.description
}
