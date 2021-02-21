resource "google_dns_managed_zone" "public" {
  count       = var.type == "public" ? 1 : 0
  project     = var.project_id
  name        = var.name
  dns_name    = format("%s.", var.dns_name)
  description = var.description
  visibility  = "public"
}

resource "google_dns_managed_zone" "private" {
  count       = var.type == "private" ? 1 : 0
  project     = var.project_id
  name        = var.name
  dns_name    = format("%s.", var.dns_name)
  description = var.description
  visibility  = "private"

  private_visibility_config {
    dynamic "networks" {
      for_each = var.private_visibility_config_networks
      content {
        network_url = networks.value
      }
    }
  }
}
