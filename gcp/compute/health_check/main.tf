locals {
  http_health_check  = var.healthcheck_type == "http" ? [{}] : []
  https_health_check = var.healthcheck_type == "https" ? [{}] : []
  http2_health_check = var.healthcheck_type == "http2" ? [{}] : []
  ssl_health_check   = var.healthcheck_type == "ssl" ? [{}] : []
  tcp_health_check   = var.healthcheck_type == "tcp" ? [{}] : []
  enable_log         = var.healthcheck_log_enabled ? [{}] : []
}

resource "google_compute_health_check" "healthcheck" {
  provider            = google-beta
  project             = var.project_id
  name                = var.healthcheck_name
  description         = var.description
  check_interval_sec  = var.check_interval_sec
  timeout_sec         = var.timeout_sec
  healthy_threshold   = var.healthy_threshold
  unhealthy_threshold = var.unhealthy_threshold

  dynamic "http2_health_check" {
    for_each = local.http2_health_check
    content {
      host               = var.host
      request_path       = var.request_path
      response           = var.response
      port               = var.port
      port_name          = var.port_name
      proxy_header       = var.proxy_header
      port_specification = var.port_specification
    }
  }

  dynamic "http_health_check" {
    for_each = local.http_health_check
    content {
      host               = var.host
      request_path       = var.request_path
      response           = var.response
      port               = var.port
      port_name          = var.port_name
      proxy_header       = var.proxy_header
      port_specification = var.port_specification
    }
  }

  dynamic "https_health_check" {
    for_each = local.https_health_check
    content {
      host               = var.host
      request_path       = var.request_path
      response           = var.response
      port               = var.port
      port_name          = var.port_name
      proxy_header       = var.proxy_header
      port_specification = var.port_specification
    }
  }

  dynamic "tcp_health_check" {
    for_each = local.tcp_health_check
    content {
      request            = var.request
      response           = var.response
      port               = var.port
      port_name          = var.port_name
      proxy_header       = var.proxy_header
      port_specification = var.port_specification
    }
  }

  dynamic "ssl_health_check" {
    for_each = local.ssl_health_check
    content {
      request            = var.request
      response           = var.response
      port               = var.port
      port_name          = var.port_name
      proxy_header       = var.proxy_header
      port_specification = var.port_specification
    }
  }

  dynamic "log_config" {
    for_each = local.enable_log
    content {
      enable = var.healthcheck_log_enabled
    }
  }
}
