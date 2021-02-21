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
	"github.com/kr/pretty"
	compute "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	gcompute "google.golang.org/api/compute/v1"
)

func TestGCPComputeZonalInstanceGroupManager(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/instance_group_manager_zonal"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	groupName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	zone := "us-central1-a"
	computeBaseURL := "https://www.googleapis.com/compute/v1"
	subnetworkLink := fmt.Sprintf("%s/projects/%s/regions/%s/subnetworks/default", computeBaseURL, projectID, "us-central1")
	template := &gcompute.InstanceTemplate{
		Name: groupName,
		Properties: &gcompute.InstanceProperties{
			MachineType: "f1-micro",
			Disks: []*gcompute.AttachedDisk{
				{
					InitializeParams: &gcompute.AttachedDiskInitializeParams{
						DiskSizeGb:  10,
						SourceImage: "projects/debian-cloud/global/images/debian-10-buster-v20200618",
					},
					Boot: true,
				},
			},
			NetworkInterfaces: []*gcompute.NetworkInterface{
				{
					Subnetwork: subnetworkLink,
					AccessConfigs: []*gcompute.AccessConfig{
						{
							Name:        "External NAT",
							Type:        "ONE_TO_ONE_NAT",
							NetworkTier: "PREMIUM",
						},
					},
				},
			},
		},
	}
	compute.CreateInstanceTemplate(t, projectID, template)
	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":         projectID,
			"name":               groupName,
			"base_instance_name": "terratest",
			"zone":               zone,
			"default_version": map[string]interface{}{
				"name":              groupName,
				"instance_template": fmt.Sprintf("projects/%s/global/instanceTemplates/%s", projectID, groupName),
			},
			"target_size": 1,
		},
	}

	defer compute.DeleteInstanceTemplate(t, projectID, groupName)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateInstanceGroupCreated(t, terraformOptions, terraformOutputs)
}

func validateInstanceGroupCreated(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	zone := out["zone"].(string)
	targetSize := int64(out["target_size"].(float64))
	mig := out["name"].(string)
	template := out["version"].([]interface{})[0].(map[string]interface{})["instance_template"].(string)

	actualMig := compute.GetInstanceGroupManager(t, projectID, zone, mig)

	assert.Equal(t, targetSize, actualMig.TargetSize)
	assert.Equal(t, template, actualMig.InstanceTemplate)
	assert.Equal(t, fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s", projectID, out["zone"].(string)), actualMig.Zone)
}

func TestGCPComputeZonalInstanceGroupManagerWithAutoHealingHealthChecks(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/instance_group_manager_zonal"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	groupName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))

	// Create health check
	hc := &gcompute.HealthCheck{
		Name:               groupName,
		CheckIntervalSec:   30,
		HealthyThreshold:   1,
		TimeoutSec:         10,
		UnhealthyThreshold: 3,
		Type:               "TCP",
		TcpHealthCheck: &gcompute.TCPHealthCheck{
			Port: 22,
		},
	}

	compute.CreateHealthCheck(t, projectID, hc)
	healthCheck := compute.GetHealthCheck(t, projectID, hc.Name)
	zone := "us-central1-a"
	computeBaseURL := "https://www.googleapis.com/compute/v1"
	subnetworkLink := fmt.Sprintf("%s/projects/%s/regions/%s/subnetworks/default", computeBaseURL, projectID, "us-central1")
	template := &gcompute.InstanceTemplate{
		Name: groupName,
		Properties: &gcompute.InstanceProperties{
			MachineType: "f1-micro",
			Disks: []*gcompute.AttachedDisk{
				{
					InitializeParams: &gcompute.AttachedDiskInitializeParams{
						DiskSizeGb:  10,
						SourceImage: "projects/debian-cloud/global/images/debian-10-buster-v20200618",
					},
					Boot: true,
				},
			},
			NetworkInterfaces: []*gcompute.NetworkInterface{
				{
					Subnetwork: subnetworkLink,
					AccessConfigs: []*gcompute.AccessConfig{
						{
							Name:        "External NAT",
							Type:        "ONE_TO_ONE_NAT",
							NetworkTier: "PREMIUM",
						},
					},
				},
			},
		},
	}
	compute.CreateInstanceTemplate(t, projectID, template)
	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":         projectID,
			"name":               groupName,
			"base_instance_name": "terratest",
			"zone":               zone,
			"auto_healing_policies": []map[string]interface{}{
				{
					"health_check":      healthCheck.SelfLink,
					"initial_delay_sec": 60,
				},
			},
			"default_version": map[string]interface{}{
				"name":              groupName,
				"instance_template": fmt.Sprintf("projects/%s/global/instanceTemplates/%s", projectID, groupName),
			},
			"target_size": 1,
		},
	}

	defer compute.DeleteInstanceTemplate(t, projectID, groupName)
	defer compute.DeleteHealthCheck(t, projectID, groupName)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateInstanceGroupWithAutoHealingCreated(t, terraformOptions, terraformOutputs, healthCheck.SelfLink)
}

func validateInstanceGroupWithAutoHealingCreated(t *testing.T, opts *terraform.Options, out map[string]interface{}, healthCheckSelfLink string) {
	projectID := out["project"].(string)
	zone := out["zone"].(string)
	targetSize := int64(out["target_size"].(float64))
	mig := out["name"].(string)
	template := out["version"].([]interface{})[0].(map[string]interface{})["instance_template"].(string)
	actualMig := compute.GetInstanceGroupManager(t, projectID, zone, mig)

	assert.Equal(t, healthCheckSelfLink, actualMig.AutoHealingPolicies[0].HealthCheck)
	assert.Equal(t, targetSize, actualMig.TargetSize)
	assert.Equal(t, template, actualMig.InstanceTemplate)
	assert.Equal(t, fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s", projectID, out["zone"].(string)), actualMig.Zone)
}

func TestGCPComputeRegionalInstanceGroupManagerAllZones(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/instance_group_manager_regional"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	groupName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"
	computeBaseURL := "https://www.googleapis.com/compute/v1"
	subnetworkLink := fmt.Sprintf("%s/projects/%s/regions/%s/subnetworks/default", computeBaseURL, projectID, "us-central1")
	template := &gcompute.InstanceTemplate{
		Name: groupName,
		Properties: &gcompute.InstanceProperties{
			MachineType: "f1-micro",
			Disks: []*gcompute.AttachedDisk{
				{
					InitializeParams: &gcompute.AttachedDiskInitializeParams{
						DiskSizeGb:  10,
						SourceImage: "projects/debian-cloud/global/images/debian-10-buster-v20200618",
					},
					Boot: true,
				},
			},
			NetworkInterfaces: []*gcompute.NetworkInterface{
				{
					Subnetwork: subnetworkLink,
					AccessConfigs: []*gcompute.AccessConfig{
						{
							Name:        "External NAT",
							Type:        "ONE_TO_ONE_NAT",
							NetworkTier: "PREMIUM",
						},
					},
				},
			},
		},
	}
	compute.CreateInstanceTemplate(t, projectID, template)
	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":         projectID,
			"name":               groupName,
			"base_instance_name": "terratest",
			"region":             region,
			"default_version": map[string]interface{}{
				"name":              groupName,
				"instance_template": fmt.Sprintf("projects/%s/global/instanceTemplates/%s", projectID, groupName),
			},
			"target_size": 1,
		},
	}

	defer compute.DeleteInstanceTemplate(t, projectID, groupName)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateInstanceGroupCreatedWithZones(t, terraformOptions, terraformOutputs)
}

func validateInstanceGroupCreatedWithZones(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	region := out["region"].(string)
	distributionPolicyZones := out["distribution_policy_zones"].([]interface{})
	targetSize := int64(out["target_size"].(float64))
	mig := out["name"].(string)
	template := out["version"].([]interface{})[0].(map[string]interface{})["instance_template"].(string)

	actualMig := compute.GetRegionalInstanceGroupManager(t, projectID, region, mig)

	assert.Equal(t, targetSize, actualMig.TargetSize)
	assert.Equal(t, template, actualMig.InstanceTemplate)

	actualZones := []string{}
	for _, zoneInfo := range actualMig.DistributionPolicy.Zones {
		actualZones = append(actualZones, zoneInfo.Zone)
	}

	expectedZones := cast.ToStringSlice(distributionPolicyZones)
	if reflect.DeepEqual(actualZones, expectedZones) {
		log.Errorf("error! actual zones and expected zones do not match:\n\n%#v\n\n", pretty.Diff(actualZones, expectedZones))

		t.Logf("expected zones:\n\n %#v\n\n", expectedZones)
		t.Logf("actual zones:\n\n %#v\n\n", actualZones)
		t.FailNow()
	}
}

func TestGCPComputeRegionalInstanceGroupManagerFewZones(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/instance_group_manager_regional"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	groupName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"
	computeBaseURL := "https://www.googleapis.com/compute/v1"
	subnetworkLink := fmt.Sprintf("%s/projects/%s/regions/%s/subnetworks/default", computeBaseURL, projectID, "us-central1")
	template := &gcompute.InstanceTemplate{
		Name: groupName,
		Properties: &gcompute.InstanceProperties{
			MachineType: "f1-micro",
			Disks: []*gcompute.AttachedDisk{
				{
					InitializeParams: &gcompute.AttachedDiskInitializeParams{
						DiskSizeGb:  10,
						SourceImage: "projects/debian-cloud/global/images/debian-10-buster-v20200618",
					},
					Boot: true,
				},
			},
			NetworkInterfaces: []*gcompute.NetworkInterface{
				{
					Subnetwork: subnetworkLink,
					AccessConfigs: []*gcompute.AccessConfig{
						{
							Name:        "External NAT",
							Type:        "ONE_TO_ONE_NAT",
							NetworkTier: "PREMIUM",
						},
					},
				},
			},
		},
	}
	compute.CreateInstanceTemplate(t, projectID, template)
	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":         projectID,
			"name":               groupName,
			"base_instance_name": "terratest",
			"region":             region,
			"distribution_policy_zones": []string{
				"us-central1-a",
				"us-central1-c",
			},
			"default_version": map[string]interface{}{
				"name":              groupName,
				"instance_template": fmt.Sprintf("projects/%s/global/instanceTemplates/%s", projectID, groupName),
			},
			"target_size": 1,
		},
	}

	defer compute.DeleteInstanceTemplate(t, projectID, groupName)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateInstanceGroupCreatedWithZones(t, terraformOptions, terraformOutputs)
}
