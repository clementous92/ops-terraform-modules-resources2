variable "name" {
  description = "Name of the resource"
  type        = string
}

variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "health_checks" {
  description = "The set of URLs to HealthCheck resources for health checking this RegionBackendService. Currently at most one health check can be specified."
  type        = list(string)
}
