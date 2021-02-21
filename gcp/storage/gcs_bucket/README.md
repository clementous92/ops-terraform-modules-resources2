# Google Storage Bucket

This terraform module provisions a Google Cloud Storage bucket with ACLs.

## Usage Example

```hcl
module "my_bucket" {
  source             = "path/to/this/module-repo"

  # Required Parameters:
  bucket_name        = "${var.name}"
  project            = "${var.project}"
  region             = "${var.region}"

  # Optional Parameters:
  storage_class      = "REGIONAL"
  default_acl        = "projectPrivate"
  force_destroy      = "true"
  versioning_enabled = true

  labels = {
    "managed-by" = "terraform"
  }

  lifecycle_rules = [{
    action = [{
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }]

    condition = [{
      age                   = 60
      created_before        = "2018-08-20"
      with_state            = "ANY" # [LIVE, ARCHIVED, or ANY ]
      matches_storage_class = ["REGIONAL"]
      num_newer_versions    = 10
    }]
  }]
```
