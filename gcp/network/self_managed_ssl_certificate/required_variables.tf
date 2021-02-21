variable "name" {
  description = "Resource name (not the domain name)"
  type        = string
}

variable "project_id" {
  description = "The name of the GCP project"
  type        = string
}

variable "private_key_file_path" {
  description = "Full path to the SSL key (.key) file"
  type        = string
}

variable "cert_file_path" {
  description = "Full path to the SSL certificate (.crt) file"
  type        = string
}
