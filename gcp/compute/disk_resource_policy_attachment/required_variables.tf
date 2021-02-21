variable "disk_name" {
  description = "The name of the disk to snapshot"
  type        = string
}

variable "policy_name" {
  description = "The name of the policy to attach"
  type        = string
}

variable "project_id" {
  description = "The name of the GCP project"
  type        = string
}

variable "zone" {
  description = "Zone where the snapshot will be managed"
  type        = string
}
