variable "description" {
  description = "An optional description of this resource."
  type        = string
  default     = null
}

variable "autoscaling_policy" {
  description = "nested mode: NestingList, min items: 1, max items: 1"
  type = set(object(
    {
      cooldown_period = number
      cpu_utilization = list(object(
        {
          target = number
        }
      ))
      load_balancing_utilization = list(object(
        {
          target = number
        }
      ))
      max_replicas = number
      metric = list(object(
        {
          name   = string
          target = number
          type   = string
        }
      ))
      min_replicas = number
      mode         = string
    }
  ))
}

