variable "name" {
  description = "The name of the vm instance"
  type        = string
}

variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "subnetwork" {
  description = "Regional subnetwork within the shared VPC network"
  type        = string
}

variable "zone" {
  description = "Zone to launch VM instance in"
  type        = string
}

