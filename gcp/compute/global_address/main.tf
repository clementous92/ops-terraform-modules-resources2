resource "google_compute_global_address" "address" {
  project     = var.project_id
  name        = var.name
  description = var.description

  // An address may only be specified for INTERNAL address types
  address       = var.internal_address
  address_type  = var.address_type
  prefix_length = var.prefix_length

  purpose = var.purpose
  network = var.network
}
