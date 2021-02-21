package test

import (
	"fmt"
	"testing"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	iam "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/iam"
	"github.com/stretchr/testify/assert"
)

func TestServiceAccountCreation(t *testing.T) {
	t.Parallel()

	// Copy the module folder to a temp folder, this helps with making it
	// clean and reproducible
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/iam/service_account"

	terraformTempDir := test_structure.CopyTerraformFolderToTemp(
		t,
		rootFolder,
		terraformFolderRelativeToRoot,
	)

	saRandomName := fmt.Sprintf("terratest-sa-%x", random.UniqueId())
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: terraformTempDir,
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"project_id":   projectID,
			"account_id":   saRandomName,
			"display_name": "Service Account for terratest",
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// terraformOutputs is a map of string string which is used to validate the actual resources are created
	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateSaCreated(t, terraformOptions, terraformOutputs)
}

func validateSaCreated(t *testing.T, terraformOptions *terraform.Options, terraformOutputs map[string]interface{}) {
	projectID := terraformOutputs["project"].(string)
	expectedSaEmail := terraformOutputs["email"].(string)

	actualSAEmail := iam.GetServiceAccount(t, projectID, expectedSaEmail)

	assert.Equal(t, expectedSaEmail, actualSAEmail.Email)
}
