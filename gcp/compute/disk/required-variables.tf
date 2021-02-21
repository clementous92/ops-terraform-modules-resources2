
variable "project_id" {
  description = "The name of the GCP project"
  type        = string
}

variable "zone" {
  description = "A reference to the zone where the disk resides"
  type        = string
}

variable "name" {
  type = string
}

variable "size" {
  description = "Size of the persistent disk, specified in GB"
  type        = number
}
