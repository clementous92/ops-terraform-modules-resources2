variable "description" {
  description = "An optional description of this resource"
  type        = string
  default     = ""
}

variable "timeout_sec" {
  description = "How many seconds to wait for the backend before considering it a failed request. Default is 30 seconds. Valid range is [1, 86400]."
  type        = number
  default     = null
}

variable "network" {
  description = "The URL of the network to which this backend service belongs. This field can only be specified when the load balancing scheme is set to INTERNAL"
  type        = string
  default     = null
}

variable "affinity_cookie_ttl_sec" {
  description = "Lifetime of cookies in seconds if session_affinity is GENERATED_COOKIE. If set to 0, the cookie is non-persistent and lasts only until the end of the browser session (or equivalent). The maximum allowed value for TTL is one day. When the load balancing scheme is INTERNAL, this field is not used."
  type        = number
  default     = null
}

variable "circuit_breakers" {
  description = "Settings controlling the volume of connections to a backend service. This field is applicable only when the load_balancing_scheme is set to INTERNAL_MANAGED and the protocol is set to HTTP, HTTPS, or HTTP2. Structure is documented below."
  type = set(object(
    {
      connect_timeout             = number // The timeout for new network connections to hosts. Structure is documented below.
      max_requests_per_connection = number // Maximum requests for a single backend connection. This parameter is respected by both the HTTP/1.1 and HTTP/2 implementations. If not specified, there is no limit. Setting this parameter to 1 will effectively disable keep alive.
      max_connections             = number // The maximum number of connections to the backend cluster. Defaults to 1024.
      max_pending_requests        = number // The maximum number of pending requests to the backend cluster. Defaults to 1024.
      max_requests                = number // The maximum number of parallel requests to the backend cluster. Defaults to 1024.
      max_retries                 = number // The maximum number of parallel retries to the backend cluster. Defaults to 3.
      connect_timeout = list(object({
        seconds = number // Span of time at a resolution of a second. Must be from 0 to 315,576,000,000 inclusive.
        nanos   = number // Span of time that's a fraction of a second at nanosecond resolution. Durations less than one second are represented with a 0 seconds field and a positive nanos field. Must be from 0 to 999,999,999 inclusive.
      }))
    }
  ))
  default = []
}

variable "consistent_hash" {
  description = "Consistent Hash-based load balancing can be used to provide soft session affinity based on HTTP headers, cookies or other properties."
  type = set(object(
    {
      http_cookie = list(object(
        {
          name = string
          path = string
          ttl = list(object(
            {
              nanos   = number
              seconds = number
            }
          ))
        }
      ))
      http_header_name  = string
      minimum_ring_size = number
    }
  ))
  default = []
}

variable "connection_draining_timeout_sec" {
  description = "Time for which instance will be drained (not accept new connections, but still work to finish started)."
  type        = number
  default     = null
}

variable "failover_policy" {
  description = "Policy for failovers."
  type = set(object({
    disable_connection_drain_on_failover = bool   // On failover or failback, this field indicates whether connection drain will be honored. Setting this to true has the following effect: connections to the old active pool are not drained. Connections to the new active pool use the timeout of 10 min (currently fixed). Setting to false has the following effect: both old and new connections will have a drain timeout of 10 min. This can be set to true only if the protocol is TCP. The default is false.
    drop_traffic_if_unhealthy            = bool   // This option is used only when no healthy VMs are detected in the primary and backup instance groups. When set to true, traffic is dropped. When set to false, new connections are sent across all VMs in the primary group. The default is false.
    failover_ratio                       = number // The value of the field must be in [0, 1]. If the ratio of the healthy VMs in the primary backend is at or below this number, traffic arriving at the load-balanced IP will be directed to the failover backend. In case where 'failoverRatio' is not set or all the VMs in the backup backend are unhealthy, the traffic will be directed back to the primary backend in the "force" mode, where traffic will be spread to the healthy VMs with the best effort, or to all VMs when no VM is healthy. This field is only used with l4 load balancing.
  }))
  default = []
}

variable "load_balancing_scheme" {
  description = "Indicates what kind of load balancing this regional backend service will be used for. A backend service created for one type of load balancing cannot be used with the other(s). Default value is INTERNAL. Possible values are INTERNAL and INTERNAL_MANAGED."
  type        = string
  default     = null
}

variable "locality_lb_policy" {
  description = "The load balancing algorithm [ROUND_ROBIN, LEAST_REQUEST, RING_HASH, RANDOM, ORIGINAL_DESTINATION, MAGLEV]"
  type        = string
  default     = null
}

variable "port_name" {
  description = "A named port on a backend instance group representing the port for communication to the backend VMs in that group. Required when the loadBalancingScheme is EXTERNAL, INTERNAL_MANAGED, or INTERNAL_SELF_MANAGED and the backends are instance groups."
  type        = string
  default     = null
}

variable "protocol" {
  description = "The protocol this RegionBackendService uses to communicate with backends. The default is HTTP. NOTE: HTTP2 is only valid for beta HTTP/2 load balancer types and may result in errors if used with the GA API. Possible values are HTTP, HTTPS, HTTP2, SSL, TCP, and UDP."
  type        = string
  default     = "HTTP"
}

variable "session_affinity" {
  description = "Type of session affinity to use. The default is NONE. Session affinity is not applicable if the protocol is UDP. Possible values are NONE, CLIENT_IP, CLIENT_IP_PORT_PROTO, CLIENT_IP_PROTO, GENERATED_COOKIE, HEADER_FIELD, and HTTP_COOKIE."
  type        = string
  default     = "NONE"
}

variable "log_config" {
  description = "This field denotes the logging options for the load balancer traffic served by this backend service. If logging is enabled, logs will be exported to Stackdriver."
  type = set(object(
    {
      enable      = bool
      sample_rate = number
    }
  ))
  default = []
}

variable "backends" {
  description = "List of backends, should be a map of key-value pairs for each backend, must have the 'group' key."
  type = set(object(
    {
      balancing_mode               = string,
      capacity_scaler              = number
      description                  = string
      failover                     = bool
      group                        = string
      max_connections              = number
      max_connections_per_endpoint = number
      max_connections_per_instance = number
      max_rate                     = number
      max_rate_per_endpoint        = number
      max_rate_per_instance        = number
      max_utilization              = number
    }
  ))
  default = []
}
