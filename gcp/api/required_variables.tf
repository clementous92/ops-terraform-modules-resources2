variable "apis" {
  type        = set(string)
  description = "list of api initials without googleapi.com suffix"
}

variable "project_id" {
  description = "Project ID"
}
