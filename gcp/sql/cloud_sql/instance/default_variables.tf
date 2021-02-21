variable "availability_type" {
  description = "This specifies whether a PostgreSQL instance should be set up for high availability (REGIONAL) or single zone (ZONAL)."
  default     = "ZONAL"
}

variable "database_instance" {
  type        = string
  default     = ""
  description = "override database instance name (defaults to the name of the created instance)"
}

variable "deletion_protection" {
  description = "Whether or not to allow Terraform to destroy the instance"
  default     = true
}

variable "disk_size" {
  default     = 10
  description = "database disk size"
}

variable "disk_type" {
  default     = "PD_SSD"
  description = "database disk type"
}

variable "database_instance_tier" {
  default     = "db-f1-micro"
  description = "database instance tier"
}

variable "database_version" {
  default     = "POSTGRES_11"
  description = "database engine version"
}

variable "activation_policy" {
  description = "The activation policy for the master instance.Can be either `ALWAYS`, `NEVER` or `ON_DEMAND`."
  type        = string
  default     = "ALWAYS"
}

variable "disk_autoresize" {
  description = "Configuration to increase storage size."
  type        = bool
  default     = true
}

variable "maintenance_window_day" {
  description = "The day of week (1-7) for the master instance maintenance."
  type        = number
  default     = 1
}

variable "maintenance_window_hour" {
  description = "The hour of day (0-23) maintenance window for the master instance maintenance."
  type        = number
  default     = 23
}

variable "maintenance_window_update_track" {
  description = "The update track of maintenance window for the master instance maintenance.Can be either `canary` or `stable`."
  type        = string
  default     = "canary"
}

variable "database_flags" {
  description = "The database flags for the master instance. See [more details](https://cloud.google.com/sql/docs/mysql/flags)"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "backup_configuration" {
  description = "The backup_configuration settings subblock for the database setings"
  type = object({
    enabled                        = bool
    binary_log_enabled             = bool
    point_in_time_recovery_enabled = string
    start_time                     = string
  })
  default = {
    enabled                        = false
    binary_log_enabled             = false
    point_in_time_recovery_enabled = null
    start_time                     = null
  }
}

variable "ip_configuration" {
  description = "The ip configuration for the master instances."
  type = object({
    authorized_networks = list(map(string))
    ipv4_enabled        = bool
    private_network     = string
    require_ssl         = bool
  })
  default = {
    authorized_networks = []
    ipv4_enabled        = true
    private_network     = null
    require_ssl         = null
  }
}

variable "pricing_plan" {
  description = "The pricing plan for the master instance."
  type        = string
  default     = "PER_USE"
}


variable "root_password" {
  description = "Initial root password. Required for MS SQL Server, ignored by MySQL and PostgreSQL"
  default     = null
}

variable "encryption_key_name" {
  description = "The full path to the encryption key used for the CMEK disk encryption"
  type        = string
  default     = null
}

variable "region" {
  description = "GCP Region"
  default     = null
}

variable "master_instance_name" {
  description = "The name of the instance that will act as the master in the replication setup. Note, this requires the master to have binary_log_enabled set, as well as existing backups."
  default     = null
}

variable "replica_configuration" {
  description = "The configuration for replication."
  type = set(object(
    {
      ca_certificate            = string // PEM representation of the trusted CA's x509 certificate.
      client_certificate        = string // PEM representation of the slave's x509 certificate.
      client_key                = string // PEM representation of the slave's private key. The corresponding public key in encoded in the client_certificate.
      connect_retry_interval    = number // The number of seconds between connect retries. (default: 60)
      dump_file_path            = string // Path to a SQL file in GCS from which slave instances are created. Format is gs://bucket/filename.
      failover_target           = bool   // Specifies if the replica is the failover target. If the field is set to true the replica will be designated as a failover replica. If the master instance fails, the replica instance will be promoted as the new master instance.
      master_heartbeat_period   = number // Time in ms between replication heartbeats.
      password                  = string // Password for the replication connection.
      ssl_cipher                = string // Permissible ciphers for use in SSL encryption.
      username                  = string // Username for replication connection.
      verify_server_certificate = bool   // True if the master's common name value is checked during the SSL handshake.
    }
  ))
  default = []
}

variable "user_labels" {
  description = "A set of key/value user label pairs to assign to the instance."
  type        = map(string)
  default = {
    "managed-by" = "terraform"
  }
}
