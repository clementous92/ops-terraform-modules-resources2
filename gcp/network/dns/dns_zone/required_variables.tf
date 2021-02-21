variable "name" {
  type        = string
  description = "User assigned name for this resource. Must be unique within the project."
}

variable "dns_name" {
  type        = string
  description = "The DNS name of this managed zone, for instance example.com."
}

variable "project_id" {
  description = "The name of the GCP project"
  type        = string
}
