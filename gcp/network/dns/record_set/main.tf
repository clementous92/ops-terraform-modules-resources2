resource "google_dns_record_set" "recordset" {
  project = var.project_id

  name = var.record_name == "" ? var.dns_name : format("%s.%s", var.record_name, var.dns_name)
  type = var.record_type
  ttl  = var.record_ttl

  managed_zone = var.managed_zone
  rrdatas      = var.record_data
}
