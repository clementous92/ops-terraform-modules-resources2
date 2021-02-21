locals {
  org_id = var.folder_id != null ? null : data.google_organization.org.org_id
}

resource "google_project" "project" {
  org_id          = local.org_id
  billing_account = var.billing_account

  name       = var.name
  project_id = var.project_id
  folder_id  = var.folder_id
}

