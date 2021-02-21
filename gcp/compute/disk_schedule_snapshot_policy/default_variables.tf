variable "on_source_disk_delete" {
  description = "Specifies the behavior to apply to scheduled snapshots when the source disk is deleted. Valid options are KEEP_AUTO_SNAPSHOTS and APPLY_RETENTION_POLICY"
  default     = "APPLY_RETENTION_POLICY"
}

variable "max_retention_days" {
  description = "Number of days to keep snapshots"
  default     = 14
}

variable "snapshot_frequency" {
  description = "Number of days between snapshots"
  default     = 1
}

variable "start_time" {
  description = "When snapshots are taken.  This must be in UTC HH:MM format."
  default     = "22:00"
}

variable "day_of_weeks" {
  description = "May contain up to seven (one for each day of the week) snapshot times. Possible values are: * MONDAY * TUESDAY * WEDNESDAY * THURSDAY * FRIDAY * SATURDAY * SUNDAY"
  type = list(object({
    day        = string
    start_time = string
  }))
  default = []
}

variable "hours_in_cycle" {
  description = "The policy will execute every nth hour starting at the specified time.  It must be in an hourly format \"HH:MM\""
  default     = 0
}

variable "labels" {
  description = "Snapshot labels "
  type        = map
  default = {
    "managed-by" = "terraform"
  }
}

variable "storage_locations" {
  description = "Cloud Storage bucket location to store the auto snapshot (regional or multi-regional)"
  type        = list(string)
  default     = []
}

variable "guest_flush" {
  description = "Whether to perform a 'guest aware' snapshot."
  default     = true
}
