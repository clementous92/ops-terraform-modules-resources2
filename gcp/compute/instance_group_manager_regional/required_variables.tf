variable "base_instance_name" {
  description = "The base instance name to use for instances in this group. The value must be a valid RFC1035 name. Supported characters are lowercase letters, numbers, and hyphens (-)"
  type        = string
}

variable "default_version" {
  description = "Default application version template. Additional versions can be specified via the `versions` variable."
  type = object({
    instance_template = string
    name              = string
  })
}

variable "versions" {
  description = "Additional application versions, target_type is either 'fixed' or 'percent'."
  type = map(object({
    instance_template = string
    target_type       = string # fixed | percent
    target_size       = number
  }))
  default = null
}


variable "name" {
  description = "The name of the instance group manager."
  type        = string
}

variable "region" {
  description = "The region that instances in this group should be created in"
  type        = string
}

variable "project_id" {
  description = "GCP Project ID"
  type        = string
}
