
variable "chain_cert_file_path" {
  description = "Full path to the SSL chain certificate (.crt) file"
  default     = ""
  # default     = "/Volumes/GoogleDrive/Shared drives/Certificates Team Drive/SSL Certificates/gd_bundle-g2-g1.crt" move this to services godaddy_cert
}

variable "region" {
  description = "if set certifcate would not be global but regional"
  default     = null
}
