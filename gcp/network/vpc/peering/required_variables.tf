variable "network" {
  type        = string
  description = "The VPC network from which the Cloud SQL instance is accessible on a private IP"
}

variable "peering_ranges_names_list" {
  description = "List of named IP address range(s) of PEERING type reserved for this service provider. Note that invoking this method with a different range when connection is already established will not reallocate already provisioned service producer subnetworks."
  type        = list(string)
}

variable "service" {
  description = "Name of the peering connection"
}
