variable "secret_id" {
  description = "secret id -  This must be unique within the project."
}

variable "project_id" {
  description = "GCP project id."
}

variable "members" {
  description = "members map with role - default role \"roles/secretmanager.secretAccessor\""
  type        = list(map(string))
}

# example for members lists
#  members = [
#    {
#      member               = "user:jane@example.com"
#      role                 = "roles/secretmanager.secretAccessor"
#    }
#]

