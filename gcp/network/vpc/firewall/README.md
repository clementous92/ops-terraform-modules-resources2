# Usage

Basic usage of this module is as follows:

```hcl
module "net-firewall" {
  source                  = "github.com/rocketlawyer/ops-terraform-modules-resources/gcp/network/vpc/firewall"
  project_id              = "my-project"
  network                 = "my-vpc"
  custom_rules = {
    ingress-sample = {
      description          = "Dummy sample ingress rule, tag-based."
      direction            = "INGRESS"
      action               = "allow"
      ranges               = ["192.168.0.0"]
      sources              = ["spam-tag"]
      targets              = ["foo-tag", "egg-tag"]
      use_service_accounts = false
      rules = [
        {
          protocol = "tcp"
          ports    = [22]
        }
      ]
      extra_attributes = {
          priority = 100
      }
    }
  }
}
```
