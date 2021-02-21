variable "description" {
  description = "An optional description of this resource"
  type        = string
  default     = ""
}

variable "ip_address" {
  description = "The IP address that this forwarding rule is serving on behalf of. when load_balancing_sheme is EXTERNAL address must be global, Internal must be regional - if empty ephemral global/regional address will be assigned"
  type        = string
  default     = null
}

variable "ip_protocol" {
  description = "The IP protocol to which this rule applies. When the load balancing scheme is INTERNAL, only TCP and UDP are valid. Possible values are TCP, UDP, ESP, AH, SCTP, and ICMP"
  type        = string
  default     = "TCP"
}

variable "backend_service" {
  description = "A BackendService to receive the matched traffic. This is used only for INTERNAL load balancing."
  type        = string
  default     = null
}

variable "load_balancing_scheme" {
  description = "This signifies what the ForwardingRule will be used for and can be EXTERNAL, INTERNAL, or INTERNAL_MANAGED. EXTERNAL is used for Classic Cloud VPN gateways, protocol forwarding to VMs from an external IP address, and HTTP(S), SSL Proxy, TCP Proxy, and Network TCP/UDP load balancers. INTERNAL is used for protocol forwarding to VMs from an internal IP address, and internal TCP/UDP load balancers. INTERNAL_MANAGED is used for internal HTTP(S) load balancers. Default value is EXTERNAL. Possible values are EXTERNAL, INTERNAL, and INTERNAL_MANAGED."
  type        = string
  default     = "EXTERNAL"
}

variable "is_mirroring_collector" {
  description = "Beta Indicates whether or not this load balancer can be used as a collector for packet mirroring. To prevent mirroring loops, instances behind this load balancer will not have their traffic mirrored even if a PacketMirroring rule applies to them. This can only be set to true for load balancers that have their loadBalancingScheme set to INTERNAL."
  type        = bool
  default     = false
}

/*
This field is used along with the target field for TargetHttpProxy, TargetHttpsProxy, TargetSslProxy, TargetTcpProxy, TargetVpnGateway, TargetPool, TargetInstance. Applicable only when IPProtocol is TCP, UDP, or SCTP, only packets addressed to ports in the specified range will be forwarded to target. Forwarding rules with the same [IPAddress, IPProtocol] pair must have disjoint port ranges. Some types of forwarding target have constraints on the acceptable ports:

TargetHttpProxy: 80, 8080
TargetHttpsProxy: 443
TargetTcpProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995, 1883, 5222
TargetSslProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995, 1883, 5222
TargetVpnGateway: 500, 4500
*/
variable "port_range" {
  description = "This field is used along with the target field for TargetHttpProxy, TargetHttpsProxy, TargetSslProxy, TargetTcpProxy, TargetVpnGateway, TargetPool, TargetInstance. Applicable only when IPProtocol is TCP, UDP, or SCTP"
  type        = string
  default     = ""
}

variable "ports" {
  description = "This field is used along with the backend_service field for internal load balancing. When the load balancing scheme is INTERNAL, a single port or a comma separated list of ports can be configured"
  type        = list(string)
  default     = null
}

variable "network" {
  description = "For internal load balancing, this field identifies the network that the load balanced IP should belong to for this Forwarding Rule. If this field is not specified, the default network will be used. This field is only used for INTERNAL load balancing."
  type        = string
  default     = null
}

variable "subnet" {
  description = "The subnetwork that the load balanced IP should belong to for this Forwarding Rule. This field is only used for INTERNAL load balancing. If the network specified is in auto subnet mode, this field is optional. However, if the network is in custom subnet mode, a subnetwork must be specified."
  type        = string
  default     = null
}

variable "target" {
  description = "The URL of the target resource to receive the matched traffic. The target must live in the same region as the forwarding rule. The forwarded traffic must be of a type appropriate to the target object."
  type        = string
  default     = null
}

variable "global_access" {
  description = " If true, clients can access ILB from all regions. Otherwise only allows from the local region the ILB is located at."
  type        = bool
  default     = false
}

variable "labels" {
  description = "Labels to apply to this forwarding rule. A list of key->value pairs."
  type        = map(string)
  default = {
    "managed-by" = "terraform"
  }
}

variable "all_ports" {
  description = "For internal TCP/UDP load balancing (i.e. load balancing scheme is INTERNAL and protocol is TCP/UDP), set this to true to allow packets addressed to any ports to be forwarded to the backends configured with this forwarding rule. Used with backend service. Cannot be set if port or portRange are set."
  type        = bool
  default     = false
}

variable "network_tier" {
  description = " The networking tier used for configuring this address. If this field is not specified, it is assumed to be PREMIUM. Possible values are PREMIUM and STANDARD."
  type        = string
  default     = "PREMIUM"
}

variable "service_label" {
  description = "An optional prefix to the service name for this Forwarding Rule. If specified, will be the first label of the fully qualified service name. This field is only used for INTERNAL load balancing."
  type        = string
  default     = null
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = null
}
