package test

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	compute "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	gcompute "google.golang.org/api/compute/v1"
)

func TestGCPComputeAutoScalingSingleInstance(t *testing.T) {
	setup(t)
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/autoscaler"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	groupName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	zone := "us-central1-a"
	computeBaseURL := "https://www.googleapis.com/compute/v1"
	subnetworkLink := fmt.Sprintf("%s/projects/%s/regions/%s/subnetworks/default", computeBaseURL, projectID, "us-central1")

	tpl := createInstanceTemplate(t, projectID, groupName, subnetworkLink)
	mig := createManagedInstanceGroup(t, projectID, zone, groupName, tpl.SelfLink)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id": projectID,
			"name":       fmt.Sprintf("%s-autoscaler", groupName),
			"target":     mig.SelfLink,
			"zone":       "us-central1-a",
			"autoscaling_policy": []map[string]interface{}{
				{
					"cooldown_period": 60,
					"max_replicas":    3,
					"min_replicas":    1,
					"cpu_utilization": []map[string]interface{}{
						{
							"target": "0.2",
						},
					},
					"load_balancing_utilization": []interface{}{},
					"metric":                     []interface{}{},
				},
			},
		},
	}
	defer teardown(t)
	defer compute.DeleteInstanceTemplate(t, projectID, groupName)
	defer compute.DeleteInstanceGroupManager(t, projectID, zone, groupName)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateAutoScalePolicy(t, terraformOptions, terraformOutputs)
}

func validateAutoScalePolicy(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	zoneURL, err := url.Parse(out["zone"].(string))
	if err != nil {
		log.Errorf("error getting instance group name from target URL %s: %v", zoneURL, err)
		t.FailNow()
	}

	zone := path.Base(zoneURL.Path)
	actualAs := compute.GetAutoScaler(t, projectID, zone, out["name"].(string))
	expectedAs := out["autoscaling_policy"].([]interface{})[0].(map[string]interface{})

	// Assert autoscaling policy
	assert.Equal(t, int64(expectedAs["cooldown_period"].(float64)), actualAs.AutoscalingPolicy.CoolDownPeriodSec)
	assert.Equal(
		t,
		expectedAs["cpu_utilization"].([]interface{})[0].(map[string]interface{})["target"].(float64),
		actualAs.AutoscalingPolicy.CpuUtilization.UtilizationTarget,
	)
	assert.Equal(t, int64(expectedAs["min_replicas"].(float64)), actualAs.AutoscalingPolicy.MinNumReplicas)
	assert.Equal(t, int64(expectedAs["max_replicas"].(float64)), actualAs.AutoscalingPolicy.MaxNumReplicas)

	// Assert group in range
	targetURL, err := url.Parse(out["target"].(string))
	if err != nil {
		log.Errorf("error getting instance group name from target URL %s: %v", targetURL, err)
		t.FailNow()
	}

	fmt.Printf("Waiting (30 seconds) for instance to come up from the autoscale...")
	time.Sleep(30 * time.Second)
	fmt.Printf("Done!\n")
	groupName := path.Base(targetURL.Path)
	instancesList := compute.InstanceGroupListInstances(t, projectID, zone, groupName)
	numInstances := int64(len(instancesList.Items))
	numInstancesInRange := (numInstances >= actualAs.AutoscalingPolicy.MinNumReplicas) && (numInstances <= actualAs.AutoscalingPolicy.MaxNumReplicas)
	assert.True(t, numInstancesInRange)
}

func createInstanceTemplate(t *testing.T, projectID, name, subnetworkLink string) *gcompute.InstanceTemplate {
	template := &gcompute.InstanceTemplate{
		Name: name,
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
	tpl := compute.GetInstanceTemplate(t, projectID, name)
	return tpl
}

func createManagedInstanceGroup(t *testing.T, projectID, zone, groupName, template string) *gcompute.InstanceGroupManager {
	manager := &gcompute.InstanceGroupManager{
		Name:             groupName,
		BaseInstanceName: "terratest",
		Versions: []*gcompute.InstanceGroupManagerVersion{
			{
				Name:             groupName,
				InstanceTemplate: template,
			},
		},
		TargetSize:      0,
		ForceSendFields: []string{"TargetSize"},
	}
	compute.CreateInstanceGroupManager(t, projectID, zone, manager)
	mig := compute.GetInstanceGroupManager(t, projectID, zone, groupName)
	return mig
}
