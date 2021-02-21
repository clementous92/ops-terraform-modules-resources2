# A service project gains access to network resources provided by its associated host project.
resource "google_compute_shared_vpc_service_project" "service" {
  service_project = var.project_id
  host_project    = var.host_project
}
