resource "google_compute_forwarding_rule" "rule" {
  provider = google-beta

  project = var.project_id
  region  = var.region

  name                   = var.name
  description            = var.description
  ip_address             = var.ip_address
  ip_protocol            = var.ip_protocol
  is_mirroring_collector = var.is_mirroring_collector
  backend_service        = var.backend_service
  load_balancing_scheme  = var.load_balancing_scheme
  network                = var.network
  port_range             = var.port_range
  ports                  = var.ports
  subnetwork             = var.subnet
  target                 = var.target
  allow_global_access    = var.global_access
  labels                 = var.labels
  all_ports              = var.all_ports
  network_tier           = var.network_tier
  service_label          = var.service_label
}
