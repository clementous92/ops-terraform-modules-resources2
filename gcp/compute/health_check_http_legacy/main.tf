resource "google_compute_http_health_check" "legacy_healthcheck" {
  provider = google-beta

  project             = var.project_id
  name                = var.name
  description         = var.description
  check_interval_sec  = var.check_interval_sec
  healthy_threshold   = var.healthy_threshold
  host                = var.host
  port                = var.port
  request_path        = var.request_path
  timeout_sec         = var.timeout_sec
  unhealthy_threshold = var.unhealthy_threshold
}
