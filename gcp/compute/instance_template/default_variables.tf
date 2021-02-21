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

variable "can_ip_forward" {
  description = "(optional)"
  type        = bool
  default     = null
}

variable "virtual_ips" {
  description = "A list of virtual IPs to assign to this compute instance"
  default     = []
}

variable "description" {
  description = "(optional)"
  type        = string
  default     = null
}

variable "enable_display" {
  description = "(optional)"
  type        = bool
  default     = null
}

variable "instance_description" {
  description = "(optional)"
  type        = string
  default     = null
}

variable "labels" {
  description = "(optional)"
  type        = map(string)
  default = {
    "managed-by" = "terraform"
  }
}

variable "machine_type" {
  description = "GCP VM type to provision"
  default     = "n1-standard-1"
}

variable "min_cpu_platform" {
  description = "Specifies a minimum CPU platform. Applicable values are the friendly names of CPU platforms, such as Intel Haswell or Intel Skylake."
  type        = string
  default     = null
}

variable "name_prefix" {
  description = "Creates a unique name beginning with the specified prefix. Conflicts with name."
  type        = string
  default     = null
}

variable "name" {
  description = "The name of the instance template. If you leave this blank, Terraform will auto-generate a unique name."
  type        = string
  default     = null
}


variable "region" {
  description = "An instance template is a global resource that is not bound to a zone or a region. However, you can still specify some regional resources in an instance template, which restricts the template to the region where that resource resides. For example, a custom subnetwork resource is tied to a specific region. Defaults to the region of the Provider if no value is given."
  type        = string
  default     = null
}

variable "tags" {
  description = "Tags to attach to the instance."
  type        = set(string)
  default     = null
}

variable "guest_accelerator" {
  description = "nested mode: NestingList, min items: 0, max items: 0"
  type = set(object(
    {
      count = number
      type  = string
    }
  ))
  default = []
}

variable "enable_shielded_vm" {
  description = "Whether to enable the Shielded VM configuration on the instance. Note that the instance image must support Shielded VMs. See https://cloud.google.com/compute/docs/images"
  default     = false
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

variable "service_account" {
  description = "Service account to use. An empty string means we'll use the default one"
  default     = ""
}

variable "service_account_scopes" {
  description = "Service account scopes"
  default     = ["cloud-platform", "userinfo-email", "compute-ro", "storage-ro", "monitoring-write", "logging-write"]
}

variable "disks" {
  description = "nested mode: NestingList, min items: 1, max items: 0"
  type = set(object(
    {
      auto_delete = bool
      boot        = bool
      device_name = string
      disk_encryption_key = list(object(
        {
          kms_key_self_link = string
        }
      ))
      disk_name    = string
      disk_size_gb = number
      disk_type    = string
      interface    = string
      labels       = map(string)
      mode         = string
      source       = string
      source_image = string
      type         = string
    }
  ))
  default = []
}

variable "boot_disk" {
  description = "boot disk settings"
  type = set(object(
    {
      auto_delete = bool
      boot        = bool
      device_name = string
      disk_encryption_key = list(object(
        {
          kms_key_self_link = string
        }
      ))
      disk_name    = string
      disk_size_gb = number
      disk_type    = string
      interface    = string
      labels       = map(string)
      mode         = string
      source       = string
      source_image = string
      type         = string
    }
  ))
  default = []
}

variable "boot_disk_size_gb" {
  description = "Size of boot disk (in GB)"
  default     = 100
}

variable "boot_disk_type" {
  description = "The GCE disk type. May be set to pd-standard or pd-ssd."
  default     = "pd-ssd"
}

variable "network" {
  description = "VPC network when using auto mode"
  type        = string
  default     = null
}

variable "subnetwork" {
  description = "subnetwork when using custom mode"
  type        = string
  default     = null
}
