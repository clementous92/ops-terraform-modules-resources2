resource "google_compute_url_map" "default" {
  project         = var.project_id
  name            = var.name
  default_service = var.default_service

  dynamic "default_url_redirect" {
    for_each = var.default_url_redirect
    content {
      https_redirect         = lookup(default_url_redirect.value, "https_redirect", null)
      redirect_response_code = lookup(default_url_redirect.value, "redirect_response_code", null)
      strip_query            = lookup(default_url_redirect.value, "strip_query", false)
    }
  }
}

