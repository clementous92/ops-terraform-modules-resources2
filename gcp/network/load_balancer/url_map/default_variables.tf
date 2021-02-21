
variable "description" {
  description = "https proxy description"
  type        = string
  default     = null
}

variable "default_service" {
  description = "The backend service or backend bucket to use when none of the given paths match."
  type        = string
  default     = null
}

variable "default_url_redirect" {
  description = "When none of the specified hostRules match, the request is redirected to a URL specified by defaultUrlRedirect. If defaultUrlRedirect is specified, defaultService or defaultRouteAction must not be set. Structure is documented below."
  type = set(object({
    https_redirect         = bool
    redirect_response_code = string
    strip_query            = bool
  }))
  default = []
}
