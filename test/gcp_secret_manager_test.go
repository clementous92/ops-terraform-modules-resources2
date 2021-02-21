package test

import (
	"fmt"
	"github.com/kr/pretty"
	"os"
	"reflect"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	secretManager "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/secret-manager"
	"github.com/stretchr/testify/assert"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func TestGCPSecretManagerNewSecretAutomaticReplication(t *testing.T) {
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/secret_manager/secret"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	secretID := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"secret_id":  secretID,
			"project_id": projectID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateNewSecretWithAutomaticReplication(t, terraformOptions, terraformOutputs)
}

func validateNewSecretWithAutomaticReplication(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	secretName := out["name"].(string)
	actualSecret := secretManager.GetSecret(t, projectID, secretName)

	assert.Equal(t, secretName, actualSecret.Name)
	r := actualSecret.GetReplication().Replication
	assert.Equal(t, &secretmanagerpb.Replication_Automatic_{Automatic: &secretmanagerpb.Replication_Automatic{}}, r)
}

func TestGCPSecretManagerNewSecretUserManagedReplication(t *testing.T) {
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/secret_manager/secret"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	secretID := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"secret_id":  secretID,
			"project_id": projectID,
			"replicas":   []string{"us-east1", "us-west1", "us-central1"},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateNewSecretWithUserManagedReplication(t, terraformOptions, terraformOutputs)
}

func validateNewSecretWithUserManagedReplication(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	secretName := out["name"].(string)
	actualSecret := secretManager.GetSecret(t, projectID, secretName)
	assert.Equal(t, secretName, actualSecret.Name)
	actualManagedReplicas := actualSecret.GetReplication().Replication

	// get the value from the outputs of terraform
	replication := out["replication"].([]interface{})[0].(map[string]interface{})
	userManaged := replication["user_managed"].([]interface{})[0].(map[string]interface{})
	replicas := userManaged["replicas"].([]interface{})

	// build the user managed replicas object
	replicaLocations := []*secretmanagerpb.Replication_UserManaged_Replica{}

	// for every replica location in output add to replicaLocations list
	for _, replica := range replicas {
		replicaLocations = append(
			replicaLocations,
			&secretmanagerpb.Replication_UserManaged_Replica{Location: replica.(map[string]interface{})["location"].(string)},
		)
	}

	userManagedReplicas := &secretmanagerpb.Replication_UserManaged{
		Replicas: replicaLocations,
	}
	expectedManagedReplicas := &secretmanagerpb.Replication_UserManaged_{UserManaged: userManagedReplicas}

	if reflect.DeepEqual(expectedManagedReplicas, actualManagedReplicas) {
		t.Log("user managed replicas locations found and match the given secret")
	} else {
		log.Errorf("error! actual user managed replicas and expected managed replicas does not match:\n\n%#v\n\n", pretty.Diff(actualManagedReplicas, expectedManagedReplicas))

		log.Infof("expected:\n\n %#v\n\n", pretty.Sprint(expectedManagedReplicas))
		log.Infof("actual:\n\n %#v\n\n", pretty.Sprint(actualManagedReplicas))
		t.FailNow()
	}
}
