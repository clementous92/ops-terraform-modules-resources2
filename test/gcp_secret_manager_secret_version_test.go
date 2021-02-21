package test

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	secretManager "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/secret-manager"
	"github.com/stretchr/testify/assert"
)

func TestGCPSecretManagerSecretVersion(t *testing.T) {
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/secret_manager/secret_version"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	secretName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))

	// create a provider.tf in the temp folder
	CreateGoogleProvider(t, projectID, tempTestFolder)

	// create new secret
	secret := secretManager.CreateSecret(t, projectID, secretName)
	s := &secretManager.Secret{
		Project: projectID,
		Name:    secretName,
		Payload: []byte("version 1 top secret data"),
	}

	v1 := secretManager.AddSecretVersion(t, s)
	fmt.Printf("Version state: %s", v1.State)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":  projectID,
			"secret_id":   secret.Name,
			"secret_data": "v2 top secret data",
		},
	}

	defer secretManager.DeleteSecret(t, projectID, secret.Name)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateNewSecretversion(t, terraformOptions, terraformOutputs, projectID)
}

func validateNewSecretversion(t *testing.T, opts *terraform.Options, out map[string]interface{}, projectID string) {
	re := regexp.MustCompile(`secrets\/(.+)\/versions`)
	name := re.FindAllStringSubmatch(out["name"].(string), 1)[0][1]
	secret := &secretManager.Secret{
		Project: projectID,
		Name:    name,
	}

	secretData := secretManager.GetSecretData(t, secret)
	assert.Equal(t, "v2 top secret data", secretData)
}
