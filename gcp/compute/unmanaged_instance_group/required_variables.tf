
variable "project_id" {
  description = "The name of the GCP project"
  type        = string
}

variable "group_name" {
  type = string
}

variable "instances" {
  type = list(string)
}
