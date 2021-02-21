resource "google_storage_bucket" "bucket" {
  name          = var.name
  location      = var.region
  project       = var.project_id
  storage_class = var.storage_class
  force_destroy = var.force_destroy
  labels        = var.labels

  dynamic "lifecycle_rule" {
    for_each = var.lifecycle_rules
    content {
      dynamic "action" {
        for_each = lookup(lifecycle_rule.value, "action", [])
        content {
          storage_class = lookup(action.value, "storage_class", null)
          type          = action.value.type
        }
      }

      dynamic "condition" {
        for_each = lookup(lifecycle_rule.value, "condition", [])
        content {
          age                   = lookup(condition.value, "age", null)
          created_before        = lookup(condition.value, "created_before", null)
          is_live               = lookup(condition.value, "is_live", null)
          matches_storage_class = lookup(condition.value, "matches_storage_class", null)
          num_newer_versions    = lookup(condition.value, "num_newer_versions", null)
          with_state            = lookup(condition.value, "with_state", null)
        }
      }
    }
  }

  versioning {
    enabled = var.versioning
  }
}

