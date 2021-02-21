variable "project_id" {
  description = "GCP Project ID"
}

variable "region" {
  description = "GCP Region"
}

variable "network" {
  type        = string
  description = "The VPC network ID"
}

variable "name" {
  description = "Router name"
}
