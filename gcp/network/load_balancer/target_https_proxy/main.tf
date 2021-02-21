resource "google_compute_target_https_proxy" "default" {
  project = var.project_id
  name    = var.name
  url_map = var.url_map

  ssl_certificates = var.ssl_certificates
  ssl_policy       = var.ssl_policy
  quic_override    = var.quic_override ? "ENABLE" : null
}
