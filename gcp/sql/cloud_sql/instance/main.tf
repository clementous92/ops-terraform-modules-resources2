locals {
  ip_configuration_enabled = length(keys(var.ip_configuration)) > 0 ? true : false

  ip_configurations = {
    enabled  = var.ip_configuration
    disabled = {}
  }
}

# Zone should be set through the provider
resource "google_sql_database_instance" "database_instance" {
  provider = google-beta

  project = var.project_id
  name    = var.database_name
  region  = var.region

  deletion_protection = var.deletion_protection
  database_version    = var.database_version
  encryption_key_name = var.encryption_key_name

  # Initial root password. Required for MS SQL Server, ignored by MySQL and PostgreSQL.
  root_password = var.root_password

  master_instance_name = var.master_instance_name

  dynamic "replica_configuration" {
    for_each = var.replica_configuration
    content {
      ca_certificate            = replica_configuration.value["ca_certificate"]
      client_certificate        = replica_configuration.value["client_certificate"]
      client_key                = replica_configuration.value["client_key"]
      connect_retry_interval    = replica_configuration.value["connect_retry_interval"]
      dump_file_path            = replica_configuration.value["dump_file_path"]
      failover_target           = replica_configuration.value["failover_target"]
      master_heartbeat_period   = replica_configuration.value["master_heartbeat_period"]
      password                  = replica_configuration.value["password"]
      ssl_cipher                = replica_configuration.value["ssl_cipher"]
      username                  = replica_configuration.value["username"]
      verify_server_certificate = replica_configuration.value["verify_server_certificate"]
    }
  }

  settings {
    tier              = var.database_instance_tier
    activation_policy = var.activation_policy
    availability_type = var.availability_type
    user_labels       = var.user_labels

    dynamic "backup_configuration" {
      for_each = [var.backup_configuration]
      content {
        binary_log_enabled             = lookup(backup_configuration.value, "binary_log_enabled", false)
        enabled                        = lookup(backup_configuration.value, "enabled", false)
        start_time                     = lookup(backup_configuration.value, "start_time", null)
        point_in_time_recovery_enabled = lookup(backup_configuration.value, "point_in_time_recovery_enabled", null)
      }
    }

    dynamic "ip_configuration" {
      for_each = [local.ip_configurations[local.ip_configuration_enabled ? "enabled" : "disabled"]]
      content {
        ipv4_enabled    = lookup(ip_configuration.value, "ipv4_enabled", null)
        private_network = lookup(ip_configuration.value, "private_network", null)
        require_ssl     = lookup(ip_configuration.value, "require_ssl", null)

        dynamic "authorized_networks" {
          for_each = lookup(ip_configuration.value, "authorized_networks", [])
          content {
            expiration_time = lookup(authorized_networks.value, "expiration_time", null)
            name            = lookup(authorized_networks.value, "name", null)
            value           = lookup(authorized_networks.value, "value", null)
          }
        }
      }
    }

    disk_autoresize = var.disk_autoresize
    disk_size       = var.disk_size
    disk_type       = var.disk_type
    pricing_plan    = var.pricing_plan

    dynamic "database_flags" {
      for_each = var.database_flags
      content {
        name  = lookup(database_flags.value, "name", null)
        value = lookup(database_flags.value, "value", null)
      }
    }

    maintenance_window {
      day          = var.maintenance_window_day
      hour         = var.maintenance_window_hour
      update_track = var.maintenance_window_update_track
    }
  }

  lifecycle {
    ignore_changes = [
      settings[0].disk_size,
      # https://github.com/hashicorp/terraform-provider-google/issues/7990
      # This setting is not actually used but terraform wants to set it
      settings.0.replication_type
    ]
  }
}
