
variable "description" {
  description = "https proxy description"
  type        = string
  default     = null
}

variable "quic_override" {
  description = "Specifies the QUIC override policy for this resource. This determines whether the load balancer will attempt to negotiate QUIC with clients or not. Can specify one of NONE, ENABLE, or DISABLE. If NONE is specified, uses the QUIC policy with no user overrides, which is equivalent to DISABLE. Default value is NONE. Possible values are NONE, ENABLE, and DISABLE."
  type        = bool
  default     = false
}

variable "ssl_policy" {
  description = "A reference to the SslPolicy resource that will be associated with the TargetHttpsProxy resource. If not set, the TargetHttpsProxy resource will not have any SSL policy configured."
  type        = string
  default     = null
}
