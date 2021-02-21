variable "name" {
  description = "Name of the resource, must be 1-63 characters long first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash"
  type        = string
}

variable "description" {
  description = "descriptive description"
  type        = string
}

variable "project_id" {
  description = "The name of the GCP project"
  type        = string
}

variable "region" {
  description = "Region where this will be deployed"
}
