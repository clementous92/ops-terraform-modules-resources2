resource "google_compute_target_pool" "pool" {
  project = var.project_id
  region  = var.region

  name             = var.name
  description      = var.description
  backup_pool      = var.backup_pool
  failover_ratio   = var.failover_ratio
  health_checks    = var.health_checks
  instances        = var.instances
  session_affinity = var.session_affinity
}
