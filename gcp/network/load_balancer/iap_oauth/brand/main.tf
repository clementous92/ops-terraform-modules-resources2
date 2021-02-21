// org_internal_only
resource "google_iap_brand" "oauth_consent_screen" {
  count             = var.create_brand ? 1 : 0
  support_email     = var.support_email
  application_title = var.application_title
  project           = var.project_id
}
