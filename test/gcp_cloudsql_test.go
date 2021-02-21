package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/martian/log"
	"github.com/gruntwork-io/terratest/modules/random"
	"google.golang.org/api/sqladmin/v1beta4"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	sql "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/cloud-sql"
	compute "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	serviceNetworking "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/service-networking"
	"github.com/stretchr/testify/assert"
	gcompute "google.golang.org/api/compute/v1"
)

func TestTerraformGCPCloudSQLInstanceDefaultVariablesResource(t *testing.T) {
	setup(t)
	// terratest::tag::0:: Since the state is local, copy the module to a temp directory to make sure state is not stale
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/sql/cloud_sql/instance"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	t.Logf("temp folder: %s", tempTestFolder)

	// terratest::tag::1:: Get the Project Id to use
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	// terratest::tag::2:: Give the example resource a unique name
	resourcePrefixName := "cloud-sql-instance"
	resourceName := fmt.Sprintf("terratest-%s-%s", resourcePrefixName, strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		// terratest::tag::3:: The path to where our Terraform code is located
		TerraformDir: tempTestFolder,

		// terratest::tag::4:: Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"database_name": resourceName,
			"project_id":    projectID,
			"region":        "us-central",
		},

		// terratest::tag::5:: Variables to pass to our Terraform code using TF_VAR_xxx environment variables
		EnvVars: map[string]string{},
	}

	// terratest::tag::9:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer teardown(t)
	defer terraform.Destroy(t, terraformOptions)

	// terratest::tag::6:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// terratest::tag::7:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	terraformOutputs := terraform.OutputAll(t, terraformOptions)

	validateDefaults(t, projectID, terraformOptions, terraformOutputs)
}

func validateDefaults(t *testing.T, projectID string, opts *terraform.Options, out map[string]interface{}) {
	output := out["output"].(map[string]interface{})
	instanceName := output["name"].(string)
	// actual database info
	instances := sql.ListInstances(t, projectID)

	for _, instance := range instances {
		if instance.Name == instanceName {
			t.Logf("instance %s foundm asserting connection and ip address...", instanceName)
			assert.Equal(t, instance.ConnectionName, output["connection_name"].(string))
			assert.Equal(t, instance.IpAddresses[0].IpAddress, output["ip_address"].([]interface{})[0].(map[string]interface{})["ip_address"].(string))
			t.Logf("assertions passed!")
			return
		}
	}
	log.Errorf("error: instance %s was not found within actual instances list: %v", instanceName, instances)
}

func TestGCPCloudSQLUserResource(t *testing.T) {
	setup(t)
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/sql/cloud_sql/user"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	t.Logf("temp folder: %s", tempTestFolder)

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	region := "us-central1"
	database := "POSTGRES_11"
	instanceName := fmt.Sprintf("terratest-cloud-sql-instance-%s", strings.ToLower(random.UniqueId()))

	instance := sql.CreateCloudSQLInstance(t, projectID, instanceName, region, database)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"username":      "cloud-sql-user",
			"project_id":    projectID,
			"instance_name": instance.Name,
			"password":      "s3cr3t",
		},

		EnvVars: map[string]string{},
	}

	// defer statements are executed in LIFO(Last-In, First-Out) order
	defer teardown(t)
	defer sql.DeleteInstance(t, projectID, instanceName)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := map[string]string{
		"username": terraform.Output(t, terraformOptions, "username"),
	}

	validateUserCreated(t, projectID, instance, terraformOptions, terraformOutputs)
}

func validateUserCreated(t *testing.T, projectID string, instance *sqladmin.DatabaseInstance, opts *terraform.Options, out map[string]string) {
	users := sql.ListUsers(t, projectID, instance.Name)
	found := false

	for _, user := range users.Items {
		if user.Name == out["username"] {
			found = true
		}
	}
	if !found {
		t.Fatalf("username %s was not found on instance %s", out["username"], instance.Name)
	}
}

func TestTerraformGCPCloudSQLPrivateInstance(t *testing.T) {
	setup(t)
	// terratest::tag::0:: Since the state is local, copy the module to a temp directory to make sure state is not stale
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/sql/cloud_sql/instance"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	t.Logf("temp folder: %s", tempTestFolder)

	// terratest::tag::1:: Get the Project Id to use
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	network := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", projectID, "default")
	addressName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))

	// create global static ip for peering
	rb := &gcompute.Address{
		Name:         addressName,
		Purpose:      "VPC_PEERING",
		AddressType:  "INTERNAL",
		PrefixLength: 16,
		Network:      network,
	}
	region := "us-central1"
	addr := compute.CreateComputeGlobalAddress(t, projectID, rb)
	service := "servicenetworking.googleapis.com"
	reservedPeeringRanges := []string{addr.Name}
	vpc := fmt.Sprintf("projects/%s/global/networks/%s", projectID, "default")
	serviceNetworking.CreatePrivateServiceConnection(t, vpc, service, reservedPeeringRanges)

	// terratest::tag::2:: Give the example resource a unique name
	resourcePrefixName := "cloud-sql-instance"
	resourceName := fmt.Sprintf("terratest-%s-%s", resourcePrefixName, strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		// terratest::tag::3:: The path to where our Terraform code is located
		TerraformDir: tempTestFolder,

		// terratest::tag::4:: Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"database_name": resourceName,
			"project_id":    projectID,
			"region":        region,
			"ip_configuration": map[string]interface{}{
				"authorized_networks": []string{},
				"ipv4_enabled":        false,
				"private_network":     network,
				"require_ssl":         nil,
			},
		},

		// terratest::tag::5:: Variables to pass to our Terraform code using TF_VAR_xxx environment variables
		EnvVars: map[string]string{},
	}

	// terratest::tag::9:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer teardown(t)
	defer terraform.Destroy(t, terraformOptions)
	defer serviceNetworking.DeletePrivateServiceConnection(t, projectID, "default", "cloudsql-postgres-googleapis-com")
	defer compute.DeleteComputeGlobalAddress(t, projectID, addr.Name)
	defer serviceNetworking.DeletePrivateServiceConnection(t, projectID, "default", "servicenetworking-googleapis-com")

	// terratest::tag::6:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	validatePrivateCloudSQLInstance(t, projectID, terraformOptions, resourceName)
}

func validatePrivateCloudSQLInstance(t *testing.T, projectID string, opts *terraform.Options, name string) {
	instance := sql.GetInstance(t, projectID, name)
	assert.Equal(t, "PRIVATE", instance.IpAddresses[0].Type)
}
