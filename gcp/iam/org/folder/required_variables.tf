variable "folder_name" {
  description = "The folder’s display name. A folder’s display name must be unique amongst its siblings, e.g. no two folders with the same parent can share the same display name. The display name must start and end with a letter or digit, may contain letters, digits, spaces, hyphens and underscores and can be no longer than 30 characters."
}

variable "parent" {
  description = "parent id to create the folder in"
}
