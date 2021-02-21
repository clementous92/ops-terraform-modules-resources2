package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	iam "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/iam"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestOrgFolderCreation(t *testing.T) {
	t.Parallel()

	// Copy the module folder to a temp folder, this helps with making it
	// clean and reproducible
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/iam/org/folder"

	terraformTempDir := test_structure.CopyTerraformFolderToTemp(
		t,
		rootFolder,
		terraformFolderRelativeToRoot,
	)

	// Get Org id from domain name
	orgDomain := os.Getenv("ORG_DOMAIN")
	if orgDomain == "" {
		log.Error("error org_domain env var is missing, add it with export ORG_DOMAIN=domain")
		t.FailNow()
	}
	parent := iam.GetOrgID(t, orgDomain)
	folderRandomName := fmt.Sprintf("terratest-%x", random.UniqueId())

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: terraformTempDir,
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"parent":      parent,
			"folder_name": folderRandomName,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// terraformOutputs is a map of string string which is used to validate the actual resources are created
	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateFolderCreated(t, terraformOptions, terraformOutputs)
}

func validateFolderCreated(t *testing.T, terraformOptions *terraform.Options, terraformOutputs map[string]interface{}) {
	expectedFolderName := terraformOutputs["display_name"].(string)
	expectedFolderID := terraformOutputs["name"].(string)

	actualFolder := iam.GetFolder(t, expectedFolderName)

	// verify the sql instance exists
	assert.Equal(t, expectedFolderID, actualFolder.Name)
}

func TestFolderInFolderCreation(t *testing.T) {
	t.Parallel()

	// Copy the module folder to a temp folder, this helps with making it
	// clean and reproducible
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/iam/org/folder"

	terraformTempDir := test_structure.CopyTerraformFolderToTemp(
		t,
		rootFolder,
		terraformFolderRelativeToRoot,
	)

	// Get Org id from domain name
	orgDomain := os.Getenv("ORG_DOMAIN")
	if orgDomain == "" {
		log.Error("error org_domain env var is missing, add it with export ORG_DOMAIN=domain")
		t.FailNow()
	}
	orgID := iam.GetOrgID(t, orgDomain)

	parentFolderName := fmt.Sprintf("terratest-%x", random.UniqueId())
	parentFolder := iam.CreateFolder(t, parentFolderName, orgID)
	folderRandomName := fmt.Sprintf("terratest-%x", random.UniqueId())

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: terraformTempDir,
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"parent":      parentFolder.Name,
			"folder_name": folderRandomName,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer iam.DeleteFolder(t, parentFolder.Name)
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// terraformOutputs is a map of string string which is used to validate the actual resources are created
	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateFolderInFolderCreated(t, terraformOptions, terraformOutputs)
}

func validateFolderInFolderCreated(t *testing.T, terraformOptions *terraform.Options, terraformOutputs map[string]interface{}) {
	expectedFolderName := terraformOutputs["display_name"].(string)
	expectedFolderID := terraformOutputs["name"].(string)
	parentID := terraformOutputs["parent"].(string)
	actualFolder := iam.GetFolder(t, expectedFolderName)

	// verify the sql instance exists
	assert.Equal(t, expectedFolderID, actualFolder.Name)
	assert.Equal(t, parentID, actualFolder.Parent)
}
