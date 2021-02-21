variable "labels" {
  description = "The labels assigned to this Secret."
  type        = map

  default = {
    "managed-by" = "terraform"
  }
}

variable "replicas" {
  description = "Replciation regions list"
  type        = list(string)
  default     = []
}
