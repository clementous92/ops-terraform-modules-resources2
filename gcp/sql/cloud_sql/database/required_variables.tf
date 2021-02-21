variable "name" {
  description = "The name of the database in the Cloud SQL instance. This does not include the project ID or instance name."
  type        = string
}

variable "instance" {
  description = "The name of the Cloud SQL instance. This does not include the project ID."
  type        = string
}

variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

