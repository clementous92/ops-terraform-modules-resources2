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
	gcompute "google.golang.org/api/compute/v1"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestGCPCloudNatBasic(t *testing.T) {
	t.Parallel()
	setup(t)
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/cloud_nat"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	routerName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"
	network := "default"

	// Create cloud router pre-requisite for cloud nat
	compute.CreateRouter(t, projectID, region, network, routerName)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":  projectID,
			"region":      region,
			"name":        name,
			"router_name": routerName,
		},
	}

	defer teardown(t)
	defer terraform.Destroy(t, terraformOptions)
	defer compute.DeleteRouter(t, projectID, region, routerName)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)

	validateCloudNatCreatedWithDefault(t, terraformOptions, terraformOutputs, routerName)
}

func validateCloudNatCreatedWithDefault(t *testing.T, opts *terraform.Options, out map[string]interface{}, routerName string) {
	output := cast.ToStringMapString(out["output"])
	projectID := output["project"]
	region := output["region"]

	router := compute.GetRouter(t, projectID, region, routerName)
	assert.Equal(t, output["name"], router.Nats[0].Name)
}

func TestGCPCloudNatManualIPs(t *testing.T) {
	t.Parallel()
	setup(t)
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/cloud_nat"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	routerName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"
	network := "default"

	// Create cloud router pre-requisite for cloud nat
	compute.CreateRouter(t, projectID, region, network, routerName)

	// create compute regional static address
	rb := &gcompute.Address{
		Name: name,
	}
	staticIP := compute.CreateComputeRegionalAddress(t, projectID, region, rb)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":                         projectID,
			"region":                             region,
			"name":                               name,
			"router_name":                        routerName,
			"nat_ip_allocate_option":             true,
			"nat_ips":                            []string{staticIP.SelfLink},
			"source_subnetwork_ip_ranges_to_nat": "LIST_OF_SUBNETWORKS",
			"subnetworks": []map[string]interface{}{
				{
					"name":                     "default",
					"source_ip_ranges_to_nat":  []string{"ALL_IP_RANGES"},
					"secondary_ip_range_names": []string{},
				},
			},
		},
	}

	defer teardown(t)
	defer compute.DeleteComputeRegionalAddress(t, projectID, region, staticIP.Name)
	defer compute.DeleteRouter(t, projectID, region, routerName)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)

	validateCloudNatCreatedManualIPs(t, terraformOptions, terraformOutputs, routerName)
}

func validateCloudNatCreatedManualIPs(t *testing.T, opts *terraform.Options, out map[string]interface{}, routerName string) {
	output := cast.ToStringMapString(out["output"])
	projectID := output["project"]
	region := output["region"]
	computeBaseURL := "https://www.googleapis.com/compute/v1"
	subnetworkLink := fmt.Sprintf("%s/projects/%s/regions/%s/subnetworks/default", computeBaseURL, projectID, region)
	router := compute.GetRouter(t, projectID, region, routerName)
	assert.Equal(t, output["nat_ip_allocate_option"], "MANUAL_ONLY")
	assert.Equal(t, router.Nats[0].Subnetworks[0].Name, subnetworkLink)
	assert.Equal(t, router.Nats[0].Subnetworks[0].SourceIpRangesToNat[0], "ALL_IP_RANGES")
	assert.Equal(t, output["name"], router.Nats[0].Name)
}
