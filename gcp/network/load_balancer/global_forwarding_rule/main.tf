resource "google_compute_global_forwarding_rule" "rule" {
  provider = google-beta

  project               = var.project_id
  name                  = var.name
  description           = var.description
  ip_address            = var.ip_address
  ip_protocol           = var.ip_protocol
  load_balancing_scheme = var.load_balancing_scheme
  network               = var.network
  port_range            = var.port_range
  target                = var.target
  labels                = var.labels
}
