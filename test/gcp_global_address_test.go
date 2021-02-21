package test

import (
	"fmt"
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

func TestGCPGlobalAddress(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/global_address"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-ip-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":       name,
			"project_id": projectID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateGlobalAddressCreated(t, terraformOptions, terraformOutputs)
}

func validateGlobalAddressCreated(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	ip := out["address"].(string)
	selflink := out["self_link"].(string)
	name := out["name"].(string)
	actualAddr := compute.GetComputeGlobalAddress(t, projectID, name)

	assert.Equal(t, ip, actualAddr.Address)
	assert.Equal(t, selflink, actualAddr.SelfLink)
}

func TestGCPGlobalInternalAddressForPeering(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/global_address"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-ip-%s", strings.ToLower(random.UniqueId()))
	network := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", projectID, "default")

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":          name,
			"project_id":    projectID,
			"address_type":  "INTERNAL",
			"prefix_length": 16,
			"purpose":       "VPC_PEERING",
			"network":       network,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateGlobalAddressForPeeringCreated(t, terraformOptions, terraformOutputs)
}

func validateGlobalAddressForPeeringCreated(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	ip := out["address"].(string)
	selflink := out["self_link"].(string)
	name := out["name"].(string)
	addressType := out["address_type"].(string)
	purpose := out["purpose"].(string)
	prefixLength := int64(out["prefix_length"].(float64))
	actualAddr := compute.GetComputeGlobalAddress(t, projectID, name)

	assert.Equal(t, ip, actualAddr.Address)
	assert.Equal(t, selflink, actualAddr.SelfLink)
	assert.Equal(t, addressType, actualAddr.AddressType)
	assert.Equal(t, purpose, actualAddr.Purpose)
	assert.Equal(t, prefixLength, actualAddr.PrefixLength)
}
