# disable-legacy-endpoints = True
# enable-oslogin           = True
# ssh-keys                 = "ssh public keys in string"
# vmdnssetting             = string // can be either ZonalOnly, ZonalPreferred,or GlobalOnly
variable "metadata" {
  type    = map(string)
  default = {}
}

variable "external_ip" {
  description = "A external IP address to attach to the instance. The default will allocate an ephemeral IP"
  type        = string
  default     = true
}

variable "static_ip" {
  description = "A external IP address to attach to the instance. The default will allocate an ephemeral IP"
  type        = string
  default     = null
}

variable "allow_stopping_for_update" {
  description = "Allow changes to the VM (resize etc..) after initial creation"
  default     = true
}

variable "description" {
  description = "Description of this VM instance"
  default     = ""
}

variable "disk_size" {
  description = "Size of boot disk (in GB)"
  default     = 100
}

variable "disks" {
  description = "A list of disks (self links) to attach to this compute instance"
  default     = []
}

variable "disk_type" {
  description = "The GCE disk type. May be set to pd-standard or pd-ssd."
  default     = "pd-ssd"
}

variable "machine_type" {
  description = "GCP VM type to provision"
  default     = "n1-standard-1"
}

variable "boot_source_image" {
  description = "Source disk image. If neither source_image nor source_image_family is specified, defaults to the latest public Debian image."
  default     = ""
}

variable "source_image_family" {
  description = "Source image family. If neither source_image nor source_image_family is specified, defaults to the latest public Debian image."
  default     = ""
}

variable "source_image_project" {
  description = "Project where the source image comes from."
  default     = ""
}

variable "tags" {
  description = "Network tags for managing firewall rules"
  type        = list(string)
  default     = []
}

variable "labels" {
  description = "(optional)"
  type        = map(string)
  default = {
    "managed-by" = "terraform"
  }
}

variable "service_account" {
  description = "Service account to use. An empty string means we'll use the default one"
  default     = ""
}

variable "service_account_scopes" {
  description = "Service account scopes"
  default     = ["cloud-platform", "userinfo-email", "compute-ro", "storage-ro", "monitoring-write", "logging-write"]
}

variable "virtual_ips" {
  description = "A list of virtual IPs to assign to this compute instance"
  default     = []
}

variable "enable_shielded_vm" {
  default     = false
  description = "Whether to enable the Shielded VM configuration on the instance. Note that the instance image must support Shielded VMs. See https://cloud.google.com/compute/docs/images"
}

variable "shielded_instance_config" {
  description = "Enable Shielded VM on this instance. Shielded VM provides verifiable integrity to prevent against malware and rootkits. Defaults to disabled."
  type = object({
    enable_secure_boot          = bool
    enable_vtpm                 = bool
    enable_integrity_monitoring = bool
  })

  default = {
    enable_secure_boot          = true
    enable_vtpm                 = true
    enable_integrity_monitoring = true
  }
}

variable "preemptible" {
  type        = bool
  description = "Allow the instance to be preempted"
  default     = false
}

variable "metadata_startup_script" {
  description = "Commands or scripts that should only run during creation"
  default     = ""
}
