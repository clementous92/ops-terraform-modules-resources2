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
)

func TestGCPComputeUnmanagedInstanceGroup(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/unmanaged_instance_group"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	groupName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	zone := "us-central1-a"
	subnetworkLink := fmt.Sprintf("projects/%s/regions/us-central1/subnetworks/default", projectID)
	a := fmt.Sprintf("terratest-vm-instance-%s", strings.ToLower(random.UniqueId()))
	b := fmt.Sprintf("terratest-vm-instance-%s", strings.ToLower(random.UniqueId()))
	compute.CreateVMInstance(t, projectID, zone, a, "f1-micro", subnetworkLink, "debian", 10, []string{})
	compute.CreateVMInstance(t, projectID, zone, b, "n1-standard-1", subnetworkLink, "debian", 10, []string{})

	instances := []string{
		fmt.Sprintf("projects/%s/zones/%s/instances/%s", projectID, zone, a),
		fmt.Sprintf("projects/%s/zones/%s/instances/%s", projectID, zone, b),
	}
	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"group_name": groupName,
			"instances":  instances,
			"project_id": projectID,
			"zone":       zone,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	defer compute.DeleteVMInstance(t, projectID, zone, a)
	defer compute.DeleteVMInstance(t, projectID, zone, b)

	terraform.InitAndApply(t, terraformOptions)
	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateUnmanagedInstanceGroupCreated(t, terraformOptions, terraformOutputs, instances)
}

func validateUnmanagedInstanceGroupCreated(t *testing.T, opts *terraform.Options, out map[string]interface{}, instances []string) {
	projectID := out["project"].(string)
	zone := out["zone"].(string)
	actualInstancesList := instances
	expectedInstancesList := []string{}

	groupName := out["name"].(string)
	instanceGroupList := compute.InstanceGroupListInstances(t, projectID, zone, groupName)

	for _, i := range instanceGroupList.Items {
		expectedInstancesList = append(expectedInstancesList, i.Instance)
	}

	reflect.DeepEqual(expectedInstancesList, actualInstancesList)
}
