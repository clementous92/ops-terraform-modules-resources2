
variable "enabled" {
  default = true
}

variable "port" {
  default = "443"
}

variable "port_name" {
  default = "https"
}

variable "zone" {
  description = "Zone where these resources are being deployed"
  default     = "us-west1-a"
}
