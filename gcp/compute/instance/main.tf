locals {
  metadata                 = length(keys(var.metadata)) > 0 ? var.metadata : {}
  shielded_instance_config = var.enable_shielded_vm ? [true] : []
  source_image             = var.boot_source_image != "" ? data.google_compute_image.image[0].self_link : data.google_compute_image.image_family[0].self_link
}

resource "google_compute_instance" "vm_instance" {
  provider = google-beta

  name                      = var.name
  project                   = var.project_id
  description               = var.description
  machine_type              = var.machine_type
  allow_stopping_for_update = var.allow_stopping_for_update
  zone                      = var.zone
  tags                      = var.tags
  labels                    = var.labels
  metadata_startup_script   = var.metadata_startup_script

  boot_disk {
    initialize_params {
      image = local.source_image
      type  = var.disk_type
      size  = var.disk_size
    }
  }

  metadata = var.metadata

  service_account {
    email  = var.service_account
    scopes = var.service_account_scopes
  }

  network_interface {
    subnetwork = var.subnetwork

    // If var.static_ip is set use that IP, otherwise this will generate an ephemeral IP
    // if you want to disable external ip set external_ip to false
    dynamic "access_config" {
      for_each = var.external_ip ? [{}] : []
      content {
        nat_ip = var.static_ip
      }
    }
  }

  // Support passing in additional IPs to assign as virtual interfaces
  dynamic "network_interface" {
    for_each = var.virtual_ips
    content {
      subnetwork = var.subnetwork
      network_ip = network_interface.value
    }
  }

  dynamic "attached_disk" {
    for_each = var.disks
    content {
      source = attached_disk.value
    }
  }

  dynamic "shielded_instance_config" {
    for_each = local.shielded_instance_config

    content {
      enable_secure_boot          = lookup(var.shielded_instance_config, "enable_secure_boot", null)
      enable_vtpm                 = lookup(var.shielded_instance_config, "enable_vtpm", null)
      enable_integrity_monitoring = lookup(var.shielded_instance_config, "enable_integrity_monitoring", null)
    }
  }

  //  scheduling must have automatic_restart be false when preemptible is true.
  scheduling {
    preemptible       = var.preemptible
    automatic_restart = ! var.preemptible
  }
}
