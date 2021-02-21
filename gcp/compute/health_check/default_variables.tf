variable "healthcheck_name" {
  description = "healthcheck name must be 1-63 characters long, and comply with RFC1035."
  type        = string
}
variable "healthcheck_log_enabled" {
  description = "enable log for the health check? default false"
  default     = false
}

variable "description" {
  description = "An optional description of this resource"
  type        = string
  default     = null
}

variable "check_interval_sec" {
  description = "How often (in seconds) to send a health check. The default value is 5 seconds."
  type        = number
  default     = 20
}

variable "healthy_threshold" {
  description = "A so-far unhealthy instance will be marked healthy after this many\nconsecutive successes. The default value is 2."
  type        = number
  default     = 2
}

variable "timeout_sec" {
  description = "How long (in seconds) to wait before claiming failure.\nThe default value is 5 seconds.  It is invalid for timeoutSec to have\ngreater value than checkIntervalSec."
  type        = number
  default     = 20
}

variable "unhealthy_threshold" {
  description = "A so-far healthy instance will be marked unhealthy after this many\nconsecutive failures. The default value is 2."
  type        = number
  default     = 2
}

variable "host" {
  description = "The value of the host header in the HTTP health check request. If left empty (default value), the public IP on behalf of which this health check is performed will be used."
  default     = null
}
variable "request_path" {
  description = "The request path of the HTTP health check request. The default value is /."
  default     = "/"
}

variable "request" {
  description = "The application data to send once the TCP connection has been established (default value is empty). If both request and response are empty, the connection establishment alone will indicate health. The request data can only be ASCII."
  default     = ""
}

variable "response" {
  description = "The bytes to match against the beginning of the response data. If left empty (the default value), any response will indicate health. The response data can only be ASCII."
  default     = ""
}

variable "port" {
  description = "The port number for the  health check request. The defaults are 80 for HTTP,  443 for TCP, SSL and HTTP2."
  default     = null
}

variable "port_name" {
  description = "Port name as defined in InstanceGroup#NamedPort#name. If both port and port_name are defined, port takes precedence."
  default     = null

}

variable "proxy_header" {
  description = "Specifies the type of proxy header to append before sending data to the backend. Default value is NONE. Possible values are NONE and PROXY_V1."
  default     = "NONE"
}

# USE_FIXED_PORT: The port number in port is used for health checking.
# USE_NAMED_PORT: The portName is used for health checking.
# USE_SERVING_PORT: For NetworkEndpointGroup, the port specified for each network endpoint is used for health checking. For other backends, the port or named port specified in the Backend Service is used for health checking. If not specified, HTTP health check follows behavior specified in port and portName fields. Possible values are USE_FIXED_PORT, USE_NAMED_PORT, and USE_SERVING_PORT.
variable "port_specification" {
  description = "Specifies how port is selected for health checking, can be one of the following values: [USE_FIXED_PORT, USE_NAMED_PORT, USE_SERVING_PORT]"
  default     = null
}
