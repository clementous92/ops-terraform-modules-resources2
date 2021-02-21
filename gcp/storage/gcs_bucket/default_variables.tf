variable "storage_class" {
  description = "The Storage Class of the new bucket. Supported values include: STANDARD, MULTI_REGIONAL, REGIONAL, NEARLINE, COLDLINE."
  default     = "REGIONAL"
}

variable "versioning" {
  description = "Version bucket objects?"
  default     = false
}

variable "lifecycle_rules" {
  type    = list
  default = []
}

variable "labels" {
  description = "The bucket's Lifecycle Rules configuration. See README for examples"
  type        = map

  default = {
    "managed-by" = "terraform"
  }
}

variable "force_destroy" {
  description = "When deleting a bucket, this boolean option will delete all contained objects."
  default     = "false"
}
