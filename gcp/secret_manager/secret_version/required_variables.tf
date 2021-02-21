variable "secret_id" {
  description = "secret id -  This must be unique within the project."
}

variable "secret_data" {
  description = "The secret data. Must be no larger than 64KiB. Note: This property is sensitive and will not be displayed in the plan."
}

variable "project_id" {
  description = "GCP project id."
}

