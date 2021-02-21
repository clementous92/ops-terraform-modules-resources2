resource "google_compute_disk" "disk" {
  name    = var.name
  project = var.project_id

  zone = var.zone
  size = var.size
  type = var.type
}
