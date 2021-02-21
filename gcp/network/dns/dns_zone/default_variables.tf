# This variable is never used, but left here for backward compatibility
variable "sharedvpc_project" {
  type    = string
  default = "not_used"
}

variable "description" {
  default = "RocketLawyer Terraform Managed Zone"
}

variable "type" {
  description = "The zone's visibility type: public zones are exposed to the Internet, while private zones are visible only to Virtual Private Cloud resources. Default value: public Possible values are: *private *public"
  default     = "public"
}

variable "private_visibility_config_networks" {
  description = "List of VPC self links that can see this zone."
  default     = []
  type        = list(string)
}

