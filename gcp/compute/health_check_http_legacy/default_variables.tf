
variable "check_interval_sec" {
  description = "How often (in seconds) to send a health check. The default value is 5 seconds."
  type        = number
  default     = 5
}

variable "description" {
  description = "An optional description of this resource. Provide this property when you create the resource."
  type        = string
  default     = ""
}

variable "healthy_threshold" {
  description = "A so-far unhealthy instance will be marked healthy after this many consecutive successes. The default value is 2."
  type        = number
  default     = 2
}

variable "host" {
  description = "The value of the host header in the HTTP health check request. If left empty (default value), the public IP on behalf of which this health check is performed will be used."
  type        = string
  default     = ""
}

variable "port" {
  description = "The TCP port number for the HTTP health check request. The default value is 80."
  type        = string
  default     = 80
}

variable "request_path" {
  description = "The request path of the HTTP health check request. The default value is /."
  type        = string
  default     = "/"
}

variable "timeout_sec" {
  description = "How long (in seconds) to wait before claiming failure. The default value is 5 seconds. It is invalid for timeoutSec to have greater value than checkIntervalSec."
  type        = number
  default     = 5
}

variable "unhealthy_threshold" {
  description = " A so-far healthy instance will be marked unhealthy after this many consecutive failures. The default value is 2."
  type        = number
  default     = 2
}
