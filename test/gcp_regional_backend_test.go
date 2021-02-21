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
	// "github.com/stretchr/testify/assert"
	gcompute "google.golang.org/api/compute/v1"
)

func TestGCPRegionalBackendDefaults(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/load_balancer/regional_backend_service"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"
	compute.CreateHealthCheck(t, projectID, &gcompute.HealthCheck{
		Name: name,
		Type: "HTTP",
		HttpHealthCheck: &gcompute.HTTPHealthCheck{
			Host: "/",
			Port: 80,
		},
	})
	hcName := compute.GetHealthCheck(t, projectID, name)
	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":          name,
			"project_id":    projectID,
			"region":        region,
			"health_checks": []string{hcName.SelfLink},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateRegionalBackendWithDefaults(t, terraformOptions, terraformOutputs)
}

func validateRegionalBackendWithDefaults(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	return
	// projectID := out["project"].(string)
	// zone := out["zone"].(string)
	// name := out["name"].(string)
	// diskSize := int64(out["size"].(float64))
	// diskType := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/diskTypes/%s", projectID, zone, out["type"].(string))

	// disk := compute.GetPersistentDisk(t, projectID, zone, name)

	// assert.Equal(t, name, disk.Name)
	// assert.Equal(t, diskSize, disk.SizeGb)
	// assert.Equal(t, diskType, disk.Type)
}
