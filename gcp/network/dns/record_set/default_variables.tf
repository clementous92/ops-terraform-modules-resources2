variable "record_type" {
  type    = string
  default = "A"
}

variable "record_ttl" {
  type        = number
  default     = 60
  description = "The time-to-live of this record set (seconds)."
}
