variable "backup_pool" {
  description = "URL to the backup target pool. Must also set failover_ratio."
  type        = string
  default     = null
}

variable "failover_ratio" {
  description = "Ratio (0 to 1) of failed nodes before using the backup pool (which must also be set)."
  type        = number
  default     = null
}

variable "description" {
  description = "Textual description field."
  type        = string
  default     = ""
}

variable "instances" {
  description = "List of instances in the pool. They can be given as URLs, or in the form of \"zone/name\". Note that the instances need not exist at the time of target pool creation, so there is no need to use the Terraform interpolators to create a dependency on the instances from the target pool."
  type        = list(string)
  default     = []
}

variable "session_affinity" {
  description = "How to distribute load. Options are \"NONE\" (no affinity). \"CLIENT_IP\" (hash of the source/dest addresses / ports), and \"CLIENT_IP_PROTO\" also includes the protocol (default \"NONE\")."
  type        = string
  default     = "NONE"
}
