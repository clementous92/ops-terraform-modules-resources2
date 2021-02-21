variable "cert_domains" {
  description = "List of domains that this certificate should respond for"
  type        = list(string)
}

variable "project_id" {
  description = "The name of the GCP project"
  type        = string
}

variable "name" {
  description = "Name to assign this SSL certificate.  Should correspond to the domain name (but no dots allowed)"
  type        = string
}
