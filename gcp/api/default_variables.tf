variable "disable_dependent_services" {
  description = "If true, services that are enabled and which depend on this service should also be disabled when this service is destroyed. If false or unset, an error will be generated if any enabled services depend on this service when destroying it."
  default     = false
}

variable "disable_on_destroy" {
  description = "If true, disable the service when the terraform resource is destroyed. Defaults to true. May be useful in the event that a project is long-lived but the infrastructure running in that project changes frequently."
  default     = false
}
