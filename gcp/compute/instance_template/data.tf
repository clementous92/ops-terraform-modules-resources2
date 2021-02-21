locals {
  default_image_name           = var.enable_shielded_vm ? "ubuntu-1804-bionic-v20200317" : "debian-10-buster-v20200805"
  default_family               = var.enable_shielded_vm ? "ubuntu-1804-lts" : "debian-10"
  default_non_shielded_project = var.source_image_project == "" && ! var.enable_shielded_vm ? "debian-cloud" : var.source_image_project
  default_shielded_project     = var.source_image_project == "" && var.enable_shielded_vm ? "gce-uefi-images" : var.source_image_project
}

# Possible premutations
# If boot_source_image is empty => find latest image by family name else use image name
# {
#   "non-shieled" + "no project" => "debian-cloud"
#   "non-shielded" + "project" => "project"
#   "shielded" + "no project" => "gce-uefi-images"
#   "shielded" + "project" => "project"
# }

// command to get gce-uefi-images
// gcloud compute images list --project gce-uefi-images --no-standard-images

// returns selected full image url
data "google_compute_image" "image" {
  count   = var.boot_source_image != "" ? 1 : 0
  project = var.enable_shielded_vm ? local.default_shielded_project : local.default_non_shielded_project
  name    = var.boot_source_image != "" ? var.boot_source_image : local.default_image_name
}

// return latest full image url from selected family
data "google_compute_image" "image_family" {
  count   = var.boot_source_image == "" ? 1 : 0
  project = var.enable_shielded_vm ? local.default_shielded_project : local.default_non_shielded_project
  family  = var.source_image_family != "" ? var.source_image_family : local.default_family
}
