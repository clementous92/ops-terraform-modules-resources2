locals {
  hourly = var.hours_in_cycle > 0
  daily  = (length(var.day_of_weeks) == 0 && var.hours_in_cycle == 0)
  weekly = length(var.day_of_weeks) > 0
}

resource "google_compute_resource_policy" "policy" {
  name    = var.policy_name
  project = var.project_id

  region = var.region

  snapshot_schedule_policy {
    schedule {
      dynamic "hourly_schedule" {
        for_each = local.hourly ? [{}] : []
        content {
          hours_in_cycle = var.hours_in_cycle
          start_time     = var.start_time
        }
      }

      dynamic "daily_schedule" {
        for_each = local.daily ? [{}] : []
        content {
          days_in_cycle = var.snapshot_frequency
          start_time    = var.start_time
        }
      }

      dynamic "weekly_schedule" {
        for_each = local.weekly ? [{}] : []
        content {
          dynamic "day_of_weeks" {
            for_each = var.day_of_weeks
            content {
              day        = day_of_weeks.value.day
              start_time = day_of_weeks.value.start_time
            }
          }
        }
      }
    }

    retention_policy {
      max_retention_days    = var.max_retention_days
      on_source_disk_delete = var.on_source_disk_delete
    }
    snapshot_properties {
      guest_flush       = var.guest_flush // perform a 'guest aware' snapshot
      labels            = var.labels
      storage_locations = var.storage_locations
    }
  }
}
