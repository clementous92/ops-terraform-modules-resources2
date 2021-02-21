package test

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	compute "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	"github.com/stretchr/testify/assert"
)

func TestGCPComputeInstanceGlobalTemplate(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/instance_template"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	diskSize := 10

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":              name,
			"network":           "default",
			"boot_disk_size_gb": diskSize,
			"project_id":        projectID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateVMInstanceFromTemplate(t, terraformOptions, terraformOutputs)
}

func TestGCPComputeInstanceGlobalTemplateCustom(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/instance_template"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	name := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"
	diskSize := 20
	subnetName := fmt.Sprintf("subnet-%s", strings.ToLower(random.UniqueId()))
	networkName := fmt.Sprintf("terratest-vpc-%s", strings.ToLower(random.UniqueId()))
	cidrName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	min := 1
	max := 254
	cidr := fmt.Sprintf("%d.50.10.0/24", rand.Intn(max-min)+min)
	compute.CreateNetwork(t, projectID, networkName)
	subnetwork := compute.CreateSubnetwork(t, projectID, region, networkName, subnetName, cidr, cidrName)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":               name,
			"subnetwork":         subnetwork.SelfLink,
			"boot_disk_size_gb":  diskSize,
			"project_id":         projectID,
			"enable_shielded_vm": true,
			"region":             "us-central1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	defer compute.DeleteNetwork(t, projectID, networkName)
	defer compute.DeleteSubnetwork(t, projectID, region, subnetName)
	defer compute.DeleteVMInstance(t, projectID, "us-central1-a", fmt.Sprintf("%s-instance", name))

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateVMInstanceFromTemplateCustom(t, terraformOptions, terraformOutputs, subnetwork.SelfLink)
}

func validateVMInstanceFromTemplate(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	templateName := out["self_link"].(string)
	name := fmt.Sprintf("%s-instance", out["name"])
	zone := "us-central1-a"

	compute.CreateVMInstanceFromTemplate(t, projectID, name, zone, templateName)
	instance := compute.GetVMInstance(t, projectID, zone, name)
	assert.Equal(t, name, instance.Name)
}

func validateVMInstanceFromTemplateCustom(t *testing.T, opts *terraform.Options, out map[string]interface{}, subnetwork string) {
	projectID := out["project"].(string)
	templateName := out["self_link"].(string)
	name := fmt.Sprintf("%s-instance", out["name"])
	zone := "us-central1-a"

	compute.CreateVMInstanceFromTemplate(t, projectID, name, zone, templateName)
	instance := compute.GetVMInstance(t, projectID, zone, name)
	assert.Equal(t, name, instance.Name)
	assert.True(t, instance.ShieldedInstanceConfig.EnableSecureBoot)
	assert.True(t, instance.ShieldedInstanceConfig.EnableVtpm)
	assert.True(t, instance.ShieldedInstanceConfig.EnableIntegrityMonitoring)
	assert.Equal(t, subnetwork, instance.NetworkInterfaces[0].Subnetwork)
	disk := compute.GetPersistentDisk(t, projectID, zone, name)
	assert.Equal(
		t,
		fmt.Sprintf("https://www.googleapis.com/compute/v1/%s", out["disk"].([]interface{})[0].(map[string]interface{})["source_image"].(string)),
		disk.SourceImage,
	)
}
