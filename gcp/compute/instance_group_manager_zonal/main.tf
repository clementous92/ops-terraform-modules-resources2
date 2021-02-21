resource "google_compute_instance_group_manager" "manager" {
  provider = google-beta

  project = var.project_id
  zone    = var.zone
  name    = var.name

  base_instance_name = var.base_instance_name
  description        = var.description
  target_pools       = var.target_pools
  target_size        = var.target_size
  wait_for_instances = var.wait_for_instances

  dynamic "stateful_disk" {
    for_each = var.stateful_disk
    content {
      device_name = stateful_disk.value["device_name"]
      delete_rule = stateful_disk.value["delete_rule"]
    }
  }

  dynamic "auto_healing_policies" {
    for_each = var.auto_healing_policies
    content {
      health_check      = auto_healing_policies.value["health_check"]
      initial_delay_sec = auto_healing_policies.value["initial_delay_sec"]
    }
  }

  dynamic "named_port" {
    for_each = var.named_port
    content {
      name = named_port.value["name"]
      port = named_port.value["port"]
    }
  }

  dynamic "update_policy" {
    for_each = var.update_policy
    content {
      max_surge_fixed         = update_policy.value["max_surge_fixed"]
      max_surge_percent       = update_policy.value["max_surge_percent"]
      max_unavailable_fixed   = update_policy.value["max_unavailable_fixed"]
      max_unavailable_percent = update_policy.value["max_unavailable_percent"]
      min_ready_sec           = update_policy.value["min_ready_sec"]
      minimal_action          = update_policy.value["minimal_action"]
      type                    = update_policy.value["type"]
    }
  }

  version {
    instance_template = var.default_version.instance_template
    name              = var.default_version.name
  }
  dynamic version {
    for_each = var.versions == null ? {} : var.versions
    iterator = version
    content {
      name              = version.key
      instance_template = version.value.instance_template
      target_size {
        fixed = (
          version.value.target_type == "fixed" ? version.value.target_size : null
        )
        percent = (
          version.value.target_type == "percent" ? version.value.target_size : null
        )
      }
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}
