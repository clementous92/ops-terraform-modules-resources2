variable "description" {
  description = "An optional textual description of the instance group manager."
  type        = string
  default     = null
}

variable "named_port" {
  description = "The named port configuration. See the section below for details on configuration"
  type = set(object(
    {
      name = string
      port = number
    }
  ))
  default = []
}

variable "target_size" {
  description = "The full URL of all target pools to which new instances in the group are added. Updating the target pools attribute does not affect existing instances."
  type        = number
  default     = null
}

variable "target_pools" {
  description = "The full URL of all target pools to which new instances in the group are added. Updating the target pools attribute does not affect existing instances."
  type        = set(string)
  default     = null
}

variable "wait_for_instances" {
  description = "Whether to wait for all instances to be created/updated before returning. Note that if this is set to true and the operation does not succeed, Terraform will continue trying until it times out."
  type        = bool
  default     = null
}

variable "auto_healing_policies" {
  description = "The autohealing policies for this managed instance group. You can specify only one value."
  type = set(object(
    {
      health_check      = string
      initial_delay_sec = number
    }
  ))
  default = []
}

variable "stateful_disk" {
  description = "Disks created on the instances that will be preserved on instance delete, update, etc - delete_rule options are \"NEVER\" and \"ON_PERMANENT_INSTANCE_DELETION\""
  type = set(object(
    {
      device_name = string
      delete_rule = string
    }
  ))
  default = []
}

variable "update_policy" {
  description = "The update policy for this managed instance group"
  type = set(object(
    {
      minimal_action = string # minimal_action - (Required) - Minimal action to be taken on an instance. You can specify either RESTART to restart existing instances or REPLACE to delete and create new instances from the target template. If you specify a RESTART, the Updater will attempt to perform that action only. However, if the Updater determines that the minimal action you specify is not enough to perform the update, it might perform a more disruptive action.
      type           = string # type - (Required) - The type of update process. You can specify either PROACTIVE so that the instance group manager proactively executes actions in order to bring instances to their target versions or OPPORTUNISTIC so that no action is proactively executed but the update will be performed as part of other actions (for example, resizes or recreateInstances calls).

      instance_redistribution_type = string # The distribution policy (Optional), i.e. which zone(s) should instances be create in. Default is all zones in given region.
      max_surge_fixed              = number # max_surge_fixed - (Optional), The maximum number of instances that can be created above the specified targetSize during the update process. Conflicts with max_surge_percent. If neither is set, defaults to 1
      max_surge_percent            = number # max_surge_percent - (Optional), The maximum number of instances(calculated as percentage) that can be created above the specified targetSize during the update process. Conflicts with max_surge_fixed.
      max_unavailable_fixed        = number # max_unavailable_fixed - (Optional), The maximum number of instances that can be unavailable during the update process. Conflicts with max_unavailable_percent. If neither is set, defaults to 1
      max_unavailable_percent      = number # max_unavailable_percent - (Optional), The maximum number of instances(calculated as percentage) that can be unavailable during the update process. Conflicts with max_unavailable_fixed.
      min_ready_sec                = number # min_ready_sec - (Optional), Minimum number of seconds to wait for after a newly created instance becomes available. This value must be from range [0, 3600]
    }
  ))
  default = []
}

variable "distribution_policy_zones" {
  description = "The distribution policy, i.e. which zone(s) should instances be create in. Default is all zones in given region."
  type        = list(string)
  default     = []
}
