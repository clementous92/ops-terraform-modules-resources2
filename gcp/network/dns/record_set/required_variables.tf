variable "project_id" {
  description = "The name of the GCP project in which this zone is provisioned"
  type        = string
}

variable "record_name" {
  type        = string
  description = "The DNS name to add within the zone (e.g www)"
}

variable "record_data" {
  type        = list(string)
  description = "list of dns record values (IP, CNAME, NS, TXT)"
}

variable "managed_zone" {
  type        = string
  description = "dns managed zone name"
}

variable "dns_name" {
  type        = string
  description = "dns name including the trailing period."
}
