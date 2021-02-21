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
	"github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestGCPCloudRouterBasic(t *testing.T) {
	setup(t)
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/cloud_router"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	routerName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"
	network := "default"

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id": projectID,
			"region":     region,
			"name":       routerName,
			"network":    network,
		},
	}

	defer teardown(t)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)

	validateCloudRouterBasic(t, terraformOptions, terraformOutputs)
}

func validateCloudRouterBasic(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	output := cast.ToStringMapString(out["output"])
	projectID := output["project"]
	region := output["region"]

	router := compute.GetRouter(t, projectID, region, output["name"])
	assert.Equal(t, output["name"], router.Name)
}

func TestGCPCloudRouterAdvance(t *testing.T) {
	setup(t)
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/cloud_router"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	routerName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"
	network := "default"

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id": projectID,
			"region":     region,
			"name":       routerName,
			"network":    network,
			"bgp": map[string]interface{}{
				"asn":               65415,
				"advertised_groups": []string{"ALL_SUBNETS"},
				"advertised_ip_ranges": []map[string]string{
					{
						"range":       "10.0.0.0/24",
						"description": "subnet-01",
					},
				},
			},
		},
	}

	defer teardown(t)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)

	validateCloudRouterAdvance(t, terraformOptions, terraformOutputs)
}

func validateCloudRouterAdvance(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	output := cast.ToStringMapString(out["output"])
	projectID := output["project"]
	region := output["region"]

	router := compute.GetRouter(t, projectID, region, output["name"])
	assert.Equal(t, output["name"], router.Name)
	assert.Equal(t, router.Bgp.AdvertiseMode, "CUSTOM")
	assert.Equal(t, router.Bgp.AdvertisedGroups[0], "ALL_SUBNETS")
	assert.Equal(t, router.Bgp.AdvertisedIpRanges[0].Range, "10.0.0.0/24")
}
