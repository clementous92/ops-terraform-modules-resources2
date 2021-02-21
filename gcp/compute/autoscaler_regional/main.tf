resource "google_compute_region_autoscaler" "autoscaler" {
  provider = google-beta
  project  = var.project_id

  description = var.description
  name        = var.name
  target      = var.target
  region      = var.region

  dynamic "autoscaling_policy" {
    for_each = var.autoscaling_policy
    content {
      cooldown_period = autoscaling_policy.value["cooldown_period"]
      max_replicas    = autoscaling_policy.value["max_replicas"]
      min_replicas    = autoscaling_policy.value["min_replicas"]
      mode            = autoscaling_policy.value["mode"]

      dynamic "cpu_utilization" {
        for_each = autoscaling_policy.value.cpu_utilization
        content {
          target = cpu_utilization.value["target"]
        }
      }

      dynamic "load_balancing_utilization" {
        for_each = autoscaling_policy.value.load_balancing_utilization
        content {
          target = load_balancing_utilization.value["target"]
        }
      }

      dynamic "metric" {
        for_each = autoscaling_policy.value.metric
        content {
          name   = metric.value["name"]
          target = metric.value["target"]
          type   = metric.value["type"]
        }
      }

    }
  }
}
