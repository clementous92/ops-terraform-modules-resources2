variable "custom_rules" {
  description = "List of custom rule definitions (refer to variables file for syntax)."
  type = map(object({
    description          = string
    direction            = string
    action               = string # (allow|deny)
    ranges               = list(string)
    sources              = list(string)
    targets              = list(string)
    use_service_accounts = bool
    rules = list(object({
      protocol = string
      ports    = list(string)
    }))
    extra_attributes = map(string) # (disabled and priority set here)
  }))
  default = {}
}
