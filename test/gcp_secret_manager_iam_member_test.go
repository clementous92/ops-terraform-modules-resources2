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
	iam "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/iam"
	secretManager "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/secret-manager"
	log "github.com/sirupsen/logrus"
	giam "google.golang.org/api/iam/v1"
)

// For this test to run, the runner must have the following permissions "roles/iam.serviceAccountTokenCreator"
func TestGCPSecretManagerIAMMember(t *testing.T) {
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/secret_manager/iam_member"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	// create a provider.tf in the temp folder
	CreateGoogleProvider(t, projectID, tempTestFolder)

	// create new secret
	secretName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	secret := secretManager.CreateSecret(t, projectID, secretName)

	// Add version to secret
	s := &secretManager.Secret{
		Project: projectID,
		Name:    secretName,
		Payload: []byte("version 1 top secret data"),
	}
	v1 := secretManager.AddSecretVersion(t, s)
	fmt.Printf("Version state: %s", v1.State)

	// create service account for admin and dev users
	admin := iam.CreateServiceAccount(t, "terratest-admin", "terratest admin user", projectID)
	dev := iam.CreateServiceAccount(t, "terratest-dev", "terratest dev user", projectID)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id": projectID,
			"secret_id":  secret.Name,
			"members": []map[string]string{
				{"member": fmt.Sprintf("serviceAccount:%s", admin.Email), "role": "roles/secretmanager.admin"},
				{"member": fmt.Sprintf("serviceAccount:%s", dev.Email)},
			},
		},
	}

	defer secretManager.DeleteSecret(t, projectID, secret.Name)
	defer iam.DeleteServiceAccount(t, admin.Email)
	defer iam.DeleteServiceAccount(t, dev.Email)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	validateIAMForSecret(t, terraformOptions, secretName, admin.Email, dev.Email)
}

func validateIAMForSecret(t *testing.T, opts *terraform.Options, secretName, adminEmail, devEmail string) {
	projectID := opts.Vars["project_id"].(string)

	actualPolicy := secretManager.GetSecretIAMPolicy(t, projectID, secretName)

	expectedPolicy := &giam.Policy{
		Bindings: []*giam.Binding{
			&giam.Binding{
				Role:    "roles/secretmanager.admin",
				Members: []string{fmt.Sprintf("serviceAccount:%s", adminEmail)},
			},
			&giam.Binding{
				Role:    "roles/secretmanager.secretAccessor",
				Members: []string{fmt.Sprintf("serviceAccount:%s", adminEmail)},
			},
		},
	}

	if reflect.DeepEqual(expectedPolicy, actualPolicy) {
		log.Errorf("error! actual policy and expected policy do not match:\n\n%#v\n\n", pretty.Diff(actualPolicy, expectedPolicy))

		t.Logf("expected policy:\n\n %#v\n\n", expectedPolicy)
		t.Logf("actual policy:\n\n %#v\n\n", actualPolicy)
		t.FailNow()
	}
}
