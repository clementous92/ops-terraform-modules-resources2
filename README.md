## Terraform Resources

**Required Terrafrom Version: >=0.13**

This repository contains **single** resources as modules to be combined together to create services or use directly when needed.

| Category | SubCategory |Module Name | Description |
|----------|-------------|------------|-------------|
| **API**  |     -           | api        | batch enable services |
| **Compute** |  managed instance group | autoscaler | autoscaler policy that can be attached to a managed instance group |
|  |                         | health_check | one of http, https, https2, ssl, tcp checks to be attached to a managed instance group |
|  |                         | instance_group_manager_regional | selected zones or all zones in region |
|  |                         | instance_group_manager_zonal | single zone |
|  |                         | instance_template | vm instance template |
|  |                         | unmanaged_instance_group | collection of vm instances |
|  |  persistent disk | disk | persistent disk - default to pd-ssd |
|  |                  | disk_resource_policy_attachment | attach a snapshot policy to disk |
|  |                  | disk_schedule_snapshot_policy | set snapshot policy for hourly, weekly, or daily |
|  |  static address | global_address | global static address |
|  |                 | regional_address | regional static address |
|  |  vm's | instance | single vm instance. |
| **IAM** |  org | folder | org folder - can use either another folder or a project as child |
|  |  | project | org project |
|  |  | service_account | IAM service account |
| **Network** | routing | cloud_nat | |
|  |  | cloud_router | |
|  |  | dns_zone | public or private DNS zone |
|  |  | record_set | DNS record set |
|  | certificates | managed_ssl_certificate | Automate certificates with google and let's encrypt |
|  | certificates | self_managed_ssl_certificate | Self managed ssl certificates |
| **VPC** | network | firewall | Allow or denied firewall rules |
|  |  | host_project | set a project to become a host project (shared-vpc) |
|  |  | network | create a network |
|  |  | routes | setup routes |
|  |  | service_project | set service projects as part of the shared-vpc |
|  |  | subnets | create automatic or custom subnets |
| **SecretManager** |  | iam_member | manage secret access by membership |
|  |  | secret | create a secret (no content) |
|  |  | secret_version | create a secret version (content) |
| **Storage** |  databases | cloud_sql > instance |  db instance  |
|  |   | cloud_sql > user |  create user for database  |
|  |   | bucket |  GCS bucket  |
|  |   | databases |  memory_store  | redis memory store (private, full) |

### Testing

The tests are configured to create a new project by default

In order to run tests you must have the following environment variables set:

* GCP_BILLING_ACCOUNT=<billing_account_number>
* ORG_DOMAIN=<org_domain_name>

If you want to use an existing project use the following environment variable:

* GOOGLE_PROJECT=<project_id>

You can either run a single test using the `go test` like so:

```bash
go test -v -run <test-function-name>
```

or use [GoTestSum](https://github.com/gotestyourself/gotestsum) to run all tests and get a useful summarized output you can install

![Demo](https://raw.githubusercontent.com/gotestyourself/gotestsum/master/docs/demo.gif)

Run the following command:

```bash
$ cd test && gotestsum --format testname
PASS test.TestTerraformPrivateVMInstance/Validate_VM_instance_with_private_ip_address_only (0.31s)
PASS test.TestTerraformPrivateVMInstance (16.36s)
PASS test.TestGCPComputeDailySnapshoptResourcePolicyDefaultValues/Validate_policy_created (0.30s)
PASS test.TestGCPComputeDailySnapshoptResourcePolicyDefaultValues (28.42s)
PASS test.TestGCPComputeDiskDefaultValues/Validate_persistent_disk_created (0.31s)
PASS test.TestGCPComputeDiskDefaultValues (28.72s)
PASS test.TestGCPComputeDiskResourcePolicyAttachment/Validate_persistent_disk_created (0.30s)
PASS test.TestGCPComputeDiskResourcePolicyAttachment (48.21s)
PASS test.TestTerraformVMInstanceCreationDefault/Validate_VM_instance_with_defaults_created (0.34s)
PASS test.TestTerraformVMInstanceCreationDefault (155.67s)
PASS test.TestTerraformVMInstanceCustom/Validate_VM_instance_with_custom_values_created (0.36s)
PASS test.TestTerraformVMInstanceCustom (169.70s)
PASS test

DONE 12 tests in 173.276s
```
