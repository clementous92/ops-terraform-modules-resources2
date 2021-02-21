locals {
  metadata                 = length(keys(var.metadata)) > 0 ? var.metadata : {}
  shielded_instance_config = var.enable_shielded_vm ? [true] : []
  boot_disk = [{
    auto_delete         = null
    boot                = true
    device_name         = null
    disk_encryption_key = []
    disk_name           = null
    disk_size_gb        = var.boot_disk_size_gb
    disk_type           = var.boot_disk_type
    interface           = null
    labels = {
      "managed-by" = "terraform"
    }
    mode         = null
    source       = null
    source_image = var.boot_source_image != "" ? data.google_compute_image.image[0].self_link : data.google_compute_image.image_family[0].self_link
    type         = null
  }]
  disks = concat(local.boot_disk, tolist(var.disks))
}

resource "google_compute_instance_template" "instance_template" {
  provider = google-beta

  name        = var.name
  name_prefix = var.name_prefix // conflicts with name

  description          = var.description
  project              = var.project_id
  machine_type         = var.machine_type
  region               = var.region
  tags                 = var.tags
  labels               = var.labels
  instance_description = var.instance_description
  can_ip_forward       = var.can_ip_forward
  enable_display       = var.enable_display
  metadata             = var.metadata
  min_cpu_platform     = var.min_cpu_platform

  dynamic "disk" {
    for_each = local.disks
    content {
      auto_delete  = disk.value["auto_delete"]
      boot         = disk.value["boot"]
      device_name  = disk.value["device_name"]
      disk_name    = disk.value["disk_name"]
      disk_size_gb = disk.value["disk_size_gb"]
      disk_type    = disk.value["disk_type"]
      interface    = disk.value["interface"]
      labels       = disk.value["labels"]
      mode         = disk.value["mode"]
      source       = disk.value["source"]
      source_image = disk.value["source_image"]
      type         = disk.value["type"]

      dynamic "disk_encryption_key" {
        for_each = disk.value.disk_encryption_key
        content {
          kms_key_self_link = disk_encryption_key.value["kms_key_self_link"]
        }
      }

    }
  }

  service_account {
    email  = var.service_account
    scopes = var.service_account_scopes
  }

  dynamic "guest_accelerator" {
    for_each = var.guest_accelerator
    content {
      count = guest_accelerator.value["count"]
      type  = guest_accelerator.value["type"]
    }
  }

  network_interface {
    # The name or self_link of the network to attach this interface to.
    # Use network attribute for Auto subnetted networks
    # subnetwork for custom subnetted networks.
    network    = var.network
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

  //  scheduling must have automatic_restart be false when preemptible is true.
  scheduling {
    preemptible       = var.preemptible
    automatic_restart = ! var.preemptible
  }


  dynamic "shielded_instance_config" {
    for_each = local.shielded_instance_config

    content {
      enable_secure_boot          = lookup(var.shielded_instance_config, "enable_secure_boot", null)
      enable_vtpm                 = lookup(var.shielded_instance_config, "enable_vtpm", null)
      enable_integrity_monitoring = lookup(var.shielded_instance_config, "enable_integrity_monitoring", null)
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}
