terraform {
  required_version = ">= 0.13"
  required_providers {
    google-beta = {
      source  = "hashicorp/google-beta"
      version = ">= 3.36"
    }
  }
}
