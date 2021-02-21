resource "google_compute_region_backend_service" "backend_service" {
  provider = google-beta
  project  = var.project_id
  region   = var.region
  network  = var.network

  name        = var.name
  description = var.description

  health_checks                   = var.health_checks
  affinity_cookie_ttl_sec         = var.affinity_cookie_ttl_sec
  connection_draining_timeout_sec = var.connection_draining_timeout_sec
  locality_lb_policy              = var.locality_lb_policy
  load_balancing_scheme           = var.load_balancing_scheme
  port_name                       = var.port_name
  protocol                        = var.protocol
  session_affinity                = var.session_affinity
  timeout_sec                     = var.timeout_sec

  dynamic "backend" {
    for_each = var.backends
    content {
      group                        = lookup(backend.value, "group", null)
      description                  = lookup(backend.value, "description", null)
      balancing_mode               = lookup(backend.value, "balancing_mode", null)
      capacity_scaler              = lookup(backend.value, "capacity_scaler", null)
      failover                     = lookup(backend.value, "failover", null)
      max_connections              = lookup(backend.value, "max_connections", null)
      max_connections_per_endpoint = lookup(backend.value, "max_connections_per_endpoint", null)
      max_connections_per_instance = lookup(backend.value, "max_connections_per_instance", null)
      max_rate                     = lookup(backend.value, "max_rate", null)
      max_rate_per_endpoint        = lookup(backend.value, "max_rate_per_endpoint", null)
      max_rate_per_instance        = lookup(backend.value, "max_rate_per_instance", null)
      max_utilization              = lookup(backend.value, "max_utilization", null)
    }
  }

  dynamic "circuit_breakers" {
    for_each = var.circuit_breakers
    content {
      max_requests_per_connection = lookup(circuit_breakers.value, "max_requests_per_connection", null)
      max_connections             = lookup(circuit_breakers.value, "max_connections", null)
      max_pending_requests        = lookup(circuit_breakers.value, "max_pending_requests", null)
      max_requests                = lookup(circuit_breakers.value, "max_requests", null)
      max_retries                 = lookup(circuit_breakers.value, "max_retries", null)
      dynamic "connect_timeout" {
        for_each = circuit_breakers.value.connect_timeout
        content {
          seconds = lookup(connect_timeout.value, "seconds", null)
          nanos   = lookup(connect_timeout.value, "nanos", null)
        }
      }
    }
  }

  dynamic "consistent_hash" {
    for_each = var.consistent_hash
    content {
      http_header_name  = lookup(consistent_hash.value, "http_header_name", null)
      minimum_ring_size = lookup(consistent_hash.value, "minimum_ring_size", null)

      dynamic "http_cookie" {
        for_each = consistent_hash.value.http_cookie
        content {
          name = lookup(http_cookie.value, "name", null)
          path = lookup(http_cookie.value, "path", null)

          dynamic "ttl" {
            for_each = http_cookie.value.ttl
            content {
              nanos   = lookup(ttl.value, "nanos", null)
              seconds = lookup(ttl.value, "seconds", null)
            }
          }
        }
      }
    }
  }

  dynamic "failover_policy" {
    for_each = var.failover_policy
    content {
      disable_connection_drain_on_failover = lookup(failover_policy.value, "disable_connection_drain_on_failover", null)
      drop_traffic_if_unhealthy            = lookup(failover_policy.value, "drop_traffic_if_unhealthy", null)
      failover_ratio                       = lookup(failover_policy.value, "failover_ratio", null)
    }
  }

  dynamic "log_config" {
    for_each = var.log_config
    content {
      enable      = lookup(log_config.value, "enable", false)
      sample_rate = lookup(log_config.value, "sample_rate", null)
    }
  }
}

