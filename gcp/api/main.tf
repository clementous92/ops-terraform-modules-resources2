# For a list of services available, visit the API library page or run gcloud services list
# https://console.cloud.google.com/apis/library

resource "google_project_service" "api" {
  for_each                   = var.apis
  project                    = var.project_id
  service                    = format("%s.googleapis.com", each.value)
  disable_dependent_services = var.disable_dependent_services
  disable_on_destroy         = var.disable_on_destroy
}
