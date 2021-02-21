variable "project_id" {
  description = "GCP project ID."
  type        = string
}

variable "name" {
  description = "Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with RFC1035"
  type        = string
}

variable "ssl_certificates" {
  description = "A list of SslCertificate resources that are used to authenticate connections between users and the load balancer. At least one SSL certificate must be specified."
  type        = list(string)
}

variable "url_map" {
  description = "self_link of UrlMap resource that defines the mapping from URL to the BackendService."
  type        = string
}
