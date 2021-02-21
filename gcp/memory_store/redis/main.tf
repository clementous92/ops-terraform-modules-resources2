locals {
  tier = var.ha ? "STANDARD_HA" : "BASIC"
}

resource "google_redis_instance" "redis" {
  project                 = var.project_id
  region                  = var.region
  location_id             = var.location_id
  alternative_location_id = var.alternative_location_id

  name          = var.name
  display_name  = var.display_name
  redis_version = var.redis_version

  tier           = local.tier
  memory_size_gb = var.memory_size_gb
  connect_mode   = var.connect_mode

  authorized_network = var.authorized_network
  reserved_ip_range  = var.reserved_ip_range

  labels = var.labels
}

