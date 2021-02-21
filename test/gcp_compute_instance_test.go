package test

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	compute "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	iam "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/iam"
	"github.com/stretchr/testify/assert"
	gcompute "google.golang.org/api/compute/v1"
)

func TestTerraformVMInstanceDefault(t *testing.T) {
	t.Parallel()

	// terratest::tag::0:: Since the state is local, copy the module to a temp directory to make sure state is not stale
	// use env var SKIP_setup=true to bypass the copy and use the local folder for
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/instance"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	// terratest::tag::1:: Get the Project Id to use
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	// terratest::tag::2:: Give the example resource a unique name
	resourcePrefixName := "vm-instance"
	resourceName := fmt.Sprintf("terratest-%s-%s", resourcePrefixName, strings.ToLower(random.UniqueId()))
	subnetworkLink := fmt.Sprintf("projects/%s/regions/us-central1/subnetworks/default", projectID)
	terraformOptions := &terraform.Options{
		// terratest::tag::3:: The path to where our Terraform code is located
		TerraformDir: tempTestFolder,

		// terratest::tag::4:: Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"name":       resourceName,
			"project_id": projectID,
			"zone":       "us-central1-a",
			"subnetwork": subnetworkLink,
		},

		// terratest::tag::5:: Variables to pass to our Terraform code using TF_VAR_xxx environment variables
		EnvVars: map[string]string{},
	}

	// terratest::tag::9:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// terratest::tag::6:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateVMInstanceDefault(t, terraformOptions, terraformOutputs)
}

func TestTerraformVMInstancePreemptible(t *testing.T) {
	t.Parallel()

	// terratest::tag::0:: Since the state is local, copy the module to a temp directory to make sure state is not stale
	// use env var SKIP_setup=true to bypass the copy and use the local folder for
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/instance"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	// terratest::tag::1:: Get the Project Id to use
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	// terratest::tag::2:: Give the example resource a unique name
	resourcePrefixName := "vm-instance"
	resourceName := fmt.Sprintf("terratest-%s-%s", resourcePrefixName, strings.ToLower(random.UniqueId()))
	subnetworkLink := fmt.Sprintf("projects/%s/regions/us-central1/subnetworks/default", projectID)
	terraformOptions := &terraform.Options{
		// terratest::tag::3:: The path to where our Terraform code is located
		TerraformDir: tempTestFolder,

		// terratest::tag::4:: Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"name":        resourceName,
			"project_id":  projectID,
			"zone":        "us-central1-a",
			"subnetwork":  subnetworkLink,
			"preemptible": true,
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

	// terratest::tag::8:: At the end of the test, run `terraform destroy` to clean up any resources that were created

	validateVMInstancePreemptible(t, terraformOptions, terraformOutputs)
}

func TestTerraformVMInstanceCustom(t *testing.T) {
	t.Parallel()

	// terratest::tag::0:: Since the state is local, copy the module to a temp directory to make sure state is not stale
	// use env var SKIP_setup=true to bypass the copy and use the local folder for
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/instance"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	// terratest::tag::1:: Get the Project Id to use
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	resourcePrefixName := "vm-instance"
	resourceName := fmt.Sprintf("terratest-%s-%s", resourcePrefixName, strings.ToLower(random.UniqueId()))
	// Create service account
	sa := iam.CreateServiceAccount(t, resourceName, "VM Instance Terratest", projectID)

	// terratest::tag::2:: Give the example resource a unique name
	subnetworkLink := fmt.Sprintf("projects/%s/regions/us-central1/subnetworks/default", projectID)
	region := "us-central1"
	// Create static IP address
	staticIPAddressName := fmt.Sprintf("terratest-%s-%s", "static-ip", strings.ToLower(random.UniqueId()))
	rb := &gcompute.Address{
		Name: staticIPAddressName,
	}
	staticIP := compute.CreateComputeRegionalAddress(t, projectID, region, rb)
	terraformOptions := &terraform.Options{
		// terratest::tag::3:: The path to where our Terraform code is located
		TerraformDir: tempTestFolder,

		// terratest::tag::4:: Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"name":            resourceName,
			"project_id":      projectID,
			"zone":            fmt.Sprintf("%s-a", region),
			"subnetwork":      subnetworkLink,
			"service_account": sa.Email,
			"tags":            []string{"abc", "robert", "john", "bastion"},
			"metadata": map[string]interface{}{
				"ssh-keys":                 "ssh-rsa12345",
				"disable-legacy-endpoints": "TRUE",
				"enable-oslogin":           "TRUE",
				"vmdnssetting":             "GlobalOnly",
				"startup-script":           "echo hi! > /tmp",
			},
			"external_ip":        true,
			"static_ip":          staticIP.Address,
			"enable_shielded_vm": true,
			"shielded_instance_config": map[string]interface{}{
				"enable_secure_boot":          true,
				"enable_vtpm":                 true,
				"enable_integrity_monitoring": true,
			},
		},

		// terratest::tag::5:: Variables to pass to our Terraform code using TF_VAR_xxx environment variables
		EnvVars: map[string]string{},
	}

	defer terraform.Destroy(t, terraformOptions)
	defer compute.DeleteComputeRegionalAddress(t, projectID, region, staticIPAddressName)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	terraformOutputs["allocatedStaticIP"] = staticIP

	validateVMInstanceCustom(t, terraformOptions, terraformOutputs)
}

func TestTerraformVMInstancePrivate(t *testing.T) {
	t.Parallel()

	// terratest::tag::0:: Since the state is local, copy the module to a temp directory to make sure state is not stale
	// use env var SKIP_setup=true to bypass the copy and use the local folder for
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/instance"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	// terratest::tag::1:: Get the Project Id to use
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	resourcePrefixName := "vm-instance"
	resourceName := fmt.Sprintf("terratest-%s-%s", resourcePrefixName, strings.ToLower(random.UniqueId()))

	// terratest::tag::2:: Give the example resource a unique name
	subnetworkLink := fmt.Sprintf("projects/%s/regions/us-central1/subnetworks/default", projectID)
	terraformOptions := &terraform.Options{
		// terratest::tag::3:: The path to where our Terraform code is located
		TerraformDir: tempTestFolder,

		// terratest::tag::4:: Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"name":               resourceName,
			"project_id":         projectID,
			"zone":               "us-central1-a",
			"subnetwork":         subnetworkLink,
			"external_ip":        false,
			"enable_shielded_vm": true,
			"shielded_instance_config": map[string]interface{}{
				"enable_secure_boot":          true,
				"enable_vtpm":                 true,
				"enable_integrity_monitoring": true,
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

	validateVMInstancePrivate(t, terraformOptions, terraformOutputs)
}

func validateVMInstanceDefault(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	zone := out["zone"].(string)
	instanceName := out["name"].(string)

	// validate default os
	expectedSourceImage := "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-10-buster-v20200805"

	vmInstance := compute.GetVMInstance(t, projectID, zone, instanceName)
	bootDisk := compute.GetPersistentDisk(t, projectID, zone, vmInstance.Name)
	assert.Equal(t, bootDisk.SourceImage, expectedSourceImage)
	assert.Equal(t, vmInstance.Name, instanceName)
}

func validateVMInstanceCustom(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	zone := out["zone"].(string)
	instanceName := out["name"].(string)
	sa := out["service_account"].([]interface{})[0].(map[string]interface{})["email"].(string)
	tags := out["tags"].([]interface{})
	metadata := out["metadata"].(interface{})
	allocatedStaticIP := out["network_interface"].([]interface{})[0].(map[string]interface{})["access_config"].([]interface{})[0].(map[string]interface{})["nat_ip"].(string)
	expectedSourceImage := "https://www.googleapis.com/compute/v1/projects/gce-uefi-images/global/images/ubuntu-1804-bionic-v20200317"
	vmInstance := compute.GetVMInstance(t, projectID, zone, instanceName)
	bootDisk := compute.GetPersistentDisk(t, projectID, zone, vmInstance.Name)

	// assertions
	assert.Equal(t, bootDisk.SourceImage, expectedSourceImage)
	reflect.DeepEqual(tags, vmInstance.Tags.Items)
	reflect.DeepEqual(sa, vmInstance.ServiceAccounts)
	reflect.DeepEqual(metadata, vmInstance.Metadata.Items)
	assert.Equal(t, vmInstance.Name, instanceName)
	assert.Equal(t, allocatedStaticIP, vmInstance.NetworkInterfaces[0].AccessConfigs[0].NatIP)

	// shielded config
	assert.True(t, vmInstance.ShieldedInstanceConfig.EnableIntegrityMonitoring, "enable integrity monitoring should be set to: true")
	assert.True(t, vmInstance.ShieldedInstanceConfig.EnableSecureBoot, "enable secure boot should be set to: true")
	assert.True(t, vmInstance.ShieldedInstanceConfig.EnableVtpm, "enable Vtpm should be set to: true")
}

func validateVMInstancePrivate(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	zone := out["zone"].(string)
	instanceName := out["name"].(string)
	vmInstance := compute.GetVMInstance(t, projectID, zone, instanceName)

	assert.Equal(t, vmInstance.Name, instanceName)

	// assert no public IP created
	assert.Equal(t, 0, len(vmInstance.NetworkInterfaces[0].AccessConfigs))
}

func validateVMInstancePreemptible(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	zone := out["zone"].(string)
	instanceName := out["name"].(string)

	vmInstance := compute.GetVMInstance(t, projectID, zone, instanceName)
	assert.Equal(t, vmInstance.Name, instanceName)
	assert.False(t, *vmInstance.Scheduling.AutomaticRestart)
	assert.True(t, *&vmInstance.Scheduling.Preemptible)
}
