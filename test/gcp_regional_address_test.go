package test

import (
	"fmt"
	"net"
	"os"
	"strings"
	"testing"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	compute "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGCPRegionalAddressDefaultExternal(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/regional_address"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-ip-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":        name,
			"project_id":  projectID,
			"description": "terratest regional external address",
			"region":      region,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateRegionalExternalAddressCreated(t, terraformOptions, terraformOutputs)
}

func validateRegionalExternalAddressCreated(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	name := out["name"].(string)
	projectID := out["project"].(string)
	region := out["region"].(string)
	ip := out["address"].(string)
	selflink := out["self_link"].(string)

	addr := compute.GetComputeRegionalAddress(t, projectID, region, name)
	assert.Equal(t, ip, addr.Address)
	assert.Equal(t, selflink, addr.SelfLink)
	assert.Equal(t, "EXTERNAL", addr.AddressType)
	assert.Equal(t, "RESERVED", addr.Status)
}

func TestGCPRegionalAddressInternalAllocatedFromSubnet(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/regional_address"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-ip-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"
	subnetworkName := "default"
	subnetworkLink := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", projectID, region, subnetworkName)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":         name,
			"project_id":   projectID,
			"description":  "terratest regional internal address",
			"region":       region,
			"subnetwork":   subnetworkLink,
			"address_type": "INTERNAL",
			"purpose":      "GCE_ENDPOINT",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateRegionalInternalAddressAllocated(t, terraformOptions, terraformOutputs, subnetworkName)
}

func validateRegionalInternalAddressAllocated(t *testing.T, opts *terraform.Options, out map[string]interface{}, subnetworkName string) {
	name := out["name"].(string)
	projectID := out["project"].(string)
	region := out["region"].(string)
	address := out["address"].(string)
	selflink := out["self_link"].(string)

	addr := compute.GetComputeRegionalAddress(t, projectID, region, name)

	// get subnet cidr range
	subnet := compute.GetSubnetwork(t, projectID, region, subnetworkName)

	// parse CIDR range
	_, ipNet, err := net.ParseCIDR(subnet.IpCidrRange)

	if err != nil {
		log.Errorf("error parsing subnet: %v", err)
		t.FailNow()
	}

	// check if the address is within that range
	addressInRange := ipNet.Contains(net.ParseIP(address))

	assert.Equal(t, address, addr.Address)
	assert.Equal(t, selflink, addr.SelfLink)
	assert.Equal(t, "GCE_ENDPOINT", addr.Purpose)
	assert.Equal(t, "INTERNAL", addr.AddressType)
	assert.Equal(t, "RESERVED", addr.Status)
	assert.True(t, addressInRange)
}

func TestGCPRegionalAddressInternalSpecified(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/regional_address"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	region := "us-central1"
	subnetworkName := "default"
	addressFromSubnet := "10.128.0.20"
	subnetworkLink := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", projectID, region, subnetworkName)
	name := fmt.Sprintf("terratest-internal-ip-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":             name,
			"project_id":       projectID,
			"description":      "terratest regional internal address",
			"region":           region,
			"subnetwork":       subnetworkLink,
			"address_type":     "INTERNAL",
			"internal_address": addressFromSubnet,
			"purpose":          "GCE_ENDPOINT",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateRegionalInternalAddressAllocated(t, terraformOptions, terraformOutputs, subnetworkName)
}
