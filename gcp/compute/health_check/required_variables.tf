variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "healthcheck_type" {
  description = "health check type, can be one of: [http, https, https2, ssl, tcp]"
  validation {
    condition     = contains(["http", "https", "http2", "ssl", "tcp"], var.healthcheck_type)
    error_message = "The health check type must be one of the following types: [http, https, https2, ssl, tcp]."
  }
}
