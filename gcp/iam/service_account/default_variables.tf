variable "display_name" {
  description = "The display name for the service account. Can be updated without creating a new resource"
  default     = null
}

variable "description" {
  description = " A text description of the service account. Must be less than or equal to 256 UTF-8 bytes."
  default     = null
}
