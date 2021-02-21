locals {
  replica_num = length(var.replicas)
}

resource "google_secret_manager_secret" "secret" {
  provider  = google-beta
  project   = var.project_id
  secret_id = var.secret_id
  labels    = var.labels

  # replication block is required by the google_secret_manager_secret resource
  # The following section would set replciation to user_managed if replicas list is not empty
  # if it is empty it will set the replication to automatic
  replication {
    automatic = local.replica_num > 0 ? null : true

    dynamic "user_managed" {
      for_each = local.replica_num > 0 ? [{}] : []
      content {
        dynamic "replicas" {
          for_each = local.replica_num > 0 ? toset(var.replicas) : []
          content {
            location = replicas.value
          }
        }
      }
    }
  }

}
