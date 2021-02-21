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

func TestGCPComputeDiskDefaultValues(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/disk"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	zone := "us-central1-a"
	diskSize := 8

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":       name,
			"size":       diskSize,
			"project_id": projectID,
			"zone":       zone,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validatePersistentDiskCreatedWithDefault(t, terraformOptions, terraformOutputs)
}

func validatePersistentDiskCreatedWithDefault(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	zone := out["zone"].(string)
	name := out["name"].(string)
	diskSize := int64(out["size"].(float64))
	diskType := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/diskTypes/%s", projectID, zone, out["type"].(string))

	disk := compute.GetPersistentDisk(t, projectID, zone, name)

	assert.Equal(t, name, disk.Name)
	assert.Equal(t, diskSize, disk.SizeGb)
	assert.Equal(t, diskType, disk.Type)
}
