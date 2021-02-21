resource "google_compute_managed_ssl_certificate" "cert" {
  provider = google-beta

  name    = var.name
  project = var.project_id

  managed {
    domains = var.cert_domains
  }
}
