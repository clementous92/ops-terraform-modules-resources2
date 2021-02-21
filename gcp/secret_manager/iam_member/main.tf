# the following transformation simply adds a default secret accesor role if no role is given
locals {
  members = {
    for _, member in var.members :
    lookup(member, "member", "") => lookup(member, "role", "roles/secretmanager.secretAccessor")
  }
}


resource "google_secret_manager_secret_iam_member" "membership" {
  for_each = local.members
  provider = google-beta

  secret_id = var.secret_id
  project   = var.project_id

  member = each.key
  role   = each.value
}
