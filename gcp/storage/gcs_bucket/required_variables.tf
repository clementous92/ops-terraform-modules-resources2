variable "name" {
  description = "The name of the bucket.  Must be globally unique."
  type        = string
}

variable "region" {
  description = "GCS bucket location"
  type        = string
}

variable "project_id" {
  description = "GCP Project ID"
}

