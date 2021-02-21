package test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	compute "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	memoryStore "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/memory-store"
	serviceNetworking "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/service-networking"
	gcompute "google.golang.org/api/compute/v1"

	"github.com/stretchr/testify/assert"
)

func TestTerraformGCPMemoryStoreRedisInstanceBasic(t *testing.T) {
	// terratest::tag::0:: Since the state is local, copy the module to a temp directory to make sure state is not stale
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/memory_store/redis"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	t.Logf("temp folder: %s", tempTestFolder)

	// terratest::tag::1:: Get the Project Id to use
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	// terratest::tag::2:: Give the example resource a unique name
	resourcePrefixName := "memory-store"
	resourceName := fmt.Sprintf("terratest-%s-%s", resourcePrefixName, strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		// terratest::tag::3:: The path to where our Terraform code is located
		TerraformDir: tempTestFolder,

		// terratest::tag::4:: Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"project_id":     projectID,
			"region":         "us-central1",
			"name":           resourceName,
			"memory_size_gb": 1,
		},

		// terratest::tag::5:: Variables to pass to our Terraform code using TF_VAR_xxx environment variables
		EnvVars: map[string]string{},
	}

	// terratest::tag::9:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// terratest::tag::6:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// terratest::tag::7:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateBasicRedisInstance(t, terraformOptions, terraformOutputs)
}

func validateBasicRedisInstance(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	instanceID := out["id"].(string)
	host := out["host"].(string)

	instance := memoryStore.GetInstance(t, instanceID)

	assert.Equal(t, host, instance.Host)
}

func TestTerraformGCPMemoryStoreRedisInstanceFull(t *testing.T) {
	// terratest::tag::0:: Since the state is local, copy the module to a temp directory to make sure state is not stale
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/memory_store/redis"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	t.Logf("temp folder: %s", tempTestFolder)

	// terratest::tag::1:: Get the Project Id to use
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	// default network
	network := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", projectID, "default")

	// terratest::tag::2:: Give the example resource a unique name
	resourcePrefixName := "memory-store"
	resourceName := fmt.Sprintf("terratest-%s-%s", resourcePrefixName, strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		// terratest::tag::3:: The path to where our Terraform code is located
		TerraformDir: tempTestFolder,

		// terratest::tag::4:: Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"project_id":              projectID,
			"region":                  "us-central1",
			"name":                    resourceName,
			"memory_size_gb":          1,
			"ha":                      true,
			"location_id":             "us-central1-a",
			"alternative_location_id": "us-central1-f",
			"authorized_network":      network,
			"redis_version":           "REDIS_3_2",
			"display_name":            "Terratest Redis HA Instance",
			"reserved_ip_range":       "192.168.0.0/29",
			"labels": map[string]string{
				"my_key":    "my_val",
				"other_key": "other_val",
			},
		},

		// terratest::tag::5:: Variables to pass to our Terraform code using TF_VAR_xxx environment variables
		EnvVars: map[string]string{},
	}

	// terratest::tag::9:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// terratest::tag::6:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// terratest::tag::7:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateFullRedisInstance(t, terraformOptions, terraformOutputs)
}

func validateFullRedisInstance(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	labelsMap := map[string]string{}

	instanceID := out["id"].(string)
	locationID := out["current_location_id"].(string)
	altLocationID := out["alternative_location_id"].(string)
	labels := out["labels"].(map[string]interface{})
	authorizedNetworks := out["authorized_network"].(string)
	redisVersion := out["redis_version"].(string)
	displayName := out["display_name"].(string)
	reservedIPRange := out["reserved_ip_range"].(string)

	instance := memoryStore.GetInstance(t, instanceID)

	assert.Equal(t, locationID, instance.LocationId)
	assert.Equal(t, altLocationID, instance.AlternativeLocationId)
	assert.Equal(t, authorizedNetworks, instance.AuthorizedNetwork)
	assert.Equal(t, redisVersion, instance.RedisVersion)
	assert.Equal(t, displayName, instance.DisplayName)
	assert.Equal(t, reservedIPRange, instance.ReservedIpRange)

	for key, value := range labels {
		strKey := fmt.Sprintf("%v", key)
		strVal := fmt.Sprintf("%v", value)
		labelsMap[strKey] = strVal
	}
	assert.True(t, reflect.DeepEqual(labelsMap, instance.Labels))
}

func TestTerraformGCPMemoryStoreRedisInstancePrivate(t *testing.T) {
	// terratest::tag::0:: Since the state is local, copy the module to a temp directory to make sure state is not stale
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/memory_store/redis"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	t.Logf("temp folder: %s", tempTestFolder)

	// terratest::tag::1:: Get the Project Id to use
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	// default network
	network := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", projectID, "default")

	addressName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))

	// create global static ip for peering
	rb := &gcompute.Address{
		Name:         addressName,
		Purpose:      "VPC_PEERING",
		AddressType:  "INTERNAL",
		PrefixLength: 24,
		Network:      network,
	}

	addr := compute.CreateComputeGlobalAddress(t, projectID, rb)
	service := "servicenetworking.googleapis.com"
	reservedPeeringRanges := []string{addr.Name}
	vpc := fmt.Sprintf("projects/%s/global/networks/%s", projectID, "default")
	serviceNetworking.CreatePrivateServiceConnection(t, vpc, service, reservedPeeringRanges)

	// terratest::tag::2:: Give the example resource a unique name
	resourcePrefixName := "memory-store"
	resourceName := fmt.Sprintf("terratest-%s-%s", resourcePrefixName, strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		// terratest::tag::3:: The path to where our Terraform code is located
		TerraformDir: tempTestFolder,

		// terratest::tag::4:: Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"project_id":              projectID,
			"region":                  "us-central1",
			"connect_mode":            "PRIVATE_SERVICE_ACCESS",
			"name":                    resourceName,
			"memory_size_gb":          1,
			"ha":                      true,
			"location_id":             "us-central1-a",
			"alternative_location_id": "us-central1-f",
			"authorized_network":      network,
			"display_name":            "Terratest Redis HA Instance",
		},

		// terratest::tag::5:: Variables to pass to our Terraform code using TF_VAR_xxx environment variables
		EnvVars: map[string]string{},
	}

	// terratest::tag::9:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)
	defer deleteRedisPeering(t, projectID, "default")
	defer compute.DeleteComputeGlobalAddress(t, projectID, addr.Name)
	defer serviceNetworking.DeletePrivateServiceConnection(t, projectID, "default", "servicenetworking-googleapis-com")

	// terratest::tag::6:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validatePrivateRedisInstance(t, terraformOptions, terraformOutputs)
}

func validatePrivateRedisInstance(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	instanceID := out["id"].(string)
	locationID := out["current_location_id"].(string)
	altLocationID := out["alternative_location_id"].(string)
	authorizedNetworks := out["authorized_network"].(string)
	redisVersion := out["redis_version"].(string)
	displayName := out["display_name"].(string)
	reservedIPRange := out["reserved_ip_range"].(string)

	instance := memoryStore.GetInstance(t, instanceID)

	assert.Equal(t, locationID, instance.LocationId)
	assert.Equal(t, altLocationID, instance.AlternativeLocationId)
	assert.Equal(t, authorizedNetworks, instance.AuthorizedNetwork)
	assert.Equal(t, redisVersion, instance.RedisVersion)
	assert.Equal(t, displayName, instance.DisplayName)
	assert.Equal(t, reservedIPRange, instance.ReservedIpRange)
}

func deleteRedisPeering(t *testing.T, projectID, networkName string) {
	network := compute.GetNetwork(t, projectID, networkName)
	for _, peer := range network.Peerings {
		if strings.HasPrefix(peer.Name, "redis") {
			compute.RemovePeering(t, projectID, networkName, peer.Name)
		}
	}
}
