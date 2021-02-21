resource "google_compute_disk_resource_policy_attachment" "snapshots" {
  name    = var.policy_name
  disk    = var.disk_name
  project = var.project_id
  zone    = var.zone
}
