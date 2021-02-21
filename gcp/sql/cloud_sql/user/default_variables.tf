variable "host" {
  type        = string
  description = "The host the user can connect from -- MUST BE null FOR POSTGRES AND MSSQL, use % for any host on MySQL"
  default     = null
}
