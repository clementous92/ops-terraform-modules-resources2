variable "name" {
  description = "The name of the autoscaler resource."
  type        = string
}

variable "target" {
  description = "The URL of the managed instance group that this autoscaler will scale."
  type        = string
}

variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "region"
  type        = string
}
