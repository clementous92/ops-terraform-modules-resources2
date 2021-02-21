variable "support_email" {
  type        = string
  description = "Support email displayed on the OAuth consent screen. Can be either a user or group email. When a user email is specified, the caller must be the user with the associated email address. When a group email is specified, the caller can be either a user or a service account which is an owner of the specified group in Cloud Identity."
}

variable "application_title" {
  type        = string
  description = "Application name displayed on OAuth consent screen."
}

variable "project_id" {
  type        = string
  description = "The ID of the project in which the resource belongs."
}



