resource "google_compute_address" "address" {
  project     = var.project_id
  name        = var.name
  description = var.description

  // An address may only be specified for INTERNAL address types
  address      = var.internal_address
  address_type = var.address_type

  purpose    = var.purpose
  region     = var.region
  subnetwork = var.subnetwork
}
