variable "address_type" {
  default = "EXTERNAL"
}

variable "internal_address" {
  default = null
}

variable "purpose" {
  description = "The purpose of the resource. which can be one the followung: [GCE_ENDPOINT] for VM instances, alias IP ranges, internal load balancers, and similar resources. [SHARED_LOADBALANCER_VIP] for an address that can be used by multiple internal load balancers This should only be set when using an Internal address. or global internal addresses it can be VPC_PEERING - for peer networks This should only be set when using an Internal address."
  default     = null
}

variable "network_tier" {
  description = "The networking tier used for configuring this address.  default PREMIUM"
  default     = "PREMIUM"
}

variable "subnetwork" {
  description = "Regional subnetwork - This field can only be used with INTERNAL type with GCE_ENDPOINT/DNS_RESOLVER purposes."
  default     = null
}
