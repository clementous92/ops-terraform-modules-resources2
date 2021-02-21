locals {
  cert = var.chain_cert_file_path != "" ? format("%s\n%s\n", file(var.cert_file_path), file(var.chain_cert_file_path)) : file(var.cert_file_path)
}

resource "google_compute_ssl_certificate" "cert" {
  count       = var.region == null ? 1 : 0
  name        = random_id.certificate.hex
  project     = var.project_id
  private_key = file(var.private_key_file_path)
  certificate = local.cert

  lifecycle {
    ignore_changes = [name]
  }
}

resource "google_compute_region_ssl_certificate" "cert" {
  count       = var.region != null ? 1 : 0
  region      = var.region
  name        = random_id.certificate.hex
  project     = var.project_id
  private_key = file(var.private_key_file_path)
  certificate = local.cert

  lifecycle {
    ignore_changes = [name]
  }
}

resource "random_id" "certificate" {
  byte_length = 4
  prefix      = format("%s-", var.name)

  # For security, do not expose raw certificate values in the output
  keepers = {
    private_key = filebase64sha256(var.private_key_file_path)
    certificate = base64sha256(local.cert)
  }
}
