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

func TestGCPComputeDiskResourcePolicyAttachment(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/disk_resource_policy_attachment"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	diskName := fmt.Sprintf("terratest-disk-%s", strings.ToLower(random.UniqueId()))
	policyName := fmt.Sprintf("terratest-snapshot-policy-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"
	zone := "us-central1-a"

	compute.CreatePersistentDisk(t, projectID, zone, diskName, 8)

	// create snapshot policy
	compute.CreateSnapshotDailySchedule(t, projectID, region, policyName)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"disk_name":   diskName,
			"policy_name": policyName,
			"project_id":  projectID,
			"zone":        zone,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	defer compute.DeletePersistentDisk(t, projectID, zone, diskName)
	defer compute.DeleteSnapshotDailySchedule(t, projectID, region, policyName)
	defer compute.DetachSnapshotScheduleFromDisk(t, projectID, region, zone, diskName, policyName)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)
	terraformOutputs["disk_name"] = diskName
	terraformOutputs["zone"] = zone
	terraformOutputs["region"] = region

	validatePolicyAttachedToDisk(t, terraformOptions, terraformOutputs, region)
}

func validatePolicyAttachedToDisk(t *testing.T, opts *terraform.Options, out map[string]interface{}, region string) {
	output := out["output"].(map[string]interface{})

	projectID := output["project"].(string)
	zone := output["zone"].(string)
	diskName := output["disk"].(string)
	policyName := output["name"].(string)
	policyURL := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/regions/%s/resourcePolicies/%s", projectID, region, policyName)
	disk := compute.GetPersistentDisk(t, projectID, zone, diskName)
	assert.Equal(t, policyURL, disk.ResourcePolicies[0])
}
