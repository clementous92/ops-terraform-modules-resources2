variable "brand" {
  type        = string
  description = "Identifier of the brand to which this client is attached to. The format is projects/{project_number}/brands/{brand_id}/identityAwareProxyClients/{client_id}."
}

variable "display_name" {
  type        = string
  description = "Human-friendly name given to the OAuth client."
}
