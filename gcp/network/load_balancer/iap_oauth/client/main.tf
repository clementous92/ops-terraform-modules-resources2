resource "google_iap_client" "project_client" {
  display_name = var.display_name
  brand        = var.brand
}
