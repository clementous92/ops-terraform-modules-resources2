package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	iam "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/iam"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestIAMProjectCreation(t *testing.T) {
	t.Parallel()
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/iam/org/project"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	billingAccount := os.Getenv("GCP_BILLING_ACCOUNT")
	if billingAccount == "" {
		t.Fatal("You must set billing account from environment variable: export GCP_BILLING_ACCOUNT=\"01A29B-000000-000000\"")
	}

	orgDomain := os.Getenv("ORG_DOMAIN")
	if orgDomain == "" {
		log.Error("error org_domain env var is missing, add it with export ORG_DOMAIN=domain")
		t.FailNow()
	}
	projectID := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	name := fmt.Sprintf("%s-01", projectID)
	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":            name,
			"project_id":      projectID,
			"billing_account": billingAccount,
			"org_domain":      orgDomain,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateProjectCreated(t, terraformOptions, terraformOutputs)
}

func validateProjectCreated(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project_id"].(string)
	proj := iam.GetProject(t, projectID)
	assert.Equal(t, projectID, proj.ProjectId)
	assert.Equal(t, "ACTIVE", proj.LifecycleState)
}

func TestIAMProjectInFolder(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/iam/org/project"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	billingAccount := os.Getenv("GCP_BILLING_ACCOUNT")
	if billingAccount == "" {
		t.Fatal("You must set billing account from environment variable: export GCP_BILLING_ACCOUNT=\"01A29B-000000-000000\"")
	}

	orgDomain := os.Getenv("ORG_DOMAIN")
	if orgDomain == "" {
		log.Error("error org_domain env var is missing, add it with export ORG_DOMAIN=domain")
		t.FailNow()
	}

	projectID := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	name := fmt.Sprintf("%s-01", projectID)

	// Get Org id from domain name
	parent := iam.GetOrgID(t, orgDomain)

	// create folder
	folderName := projectID
	folder := iam.CreateFolder(t, folderName, parent)

	// set the project within a folder
	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":            name,
			"project_id":      projectID,
			"billing_account": billingAccount,
			"org_domain":      orgDomain,
			"folder_id":       folder.Name,
		},
	}

	defer iam.DeleteFolder(t, folder.Name)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateProjectCreatedInAFolder(t, terraformOptions, terraformOutputs)
}

func validateProjectCreatedInAFolder(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project_id"].(string)
	folderID := out["folder_id"].(string)
	proj := iam.GetProject(t, projectID)
	assert.Equal(t, folderID, proj.Parent.Id)
	assert.Equal(t, "folder", proj.Parent.Type)
}
