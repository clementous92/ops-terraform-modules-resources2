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
	gcs "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/gcs"
	helpers "github.com/rocketlawyer/ops-terraform-terratest/modules/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/googleapi"
	storage "google.golang.org/api/storage/v1"
)

func TestGCPStorageBucketDefaultValues(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/storage/gcs_bucket"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":       name,
			"project_id": projectID,
			"region":     region,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateStorageBucketCreatedWithDefault(t, terraformOptions, terraformOutputs)
}

func validateStorageBucketCreatedWithDefault(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	name := helpers.KeyExists(t, "name", out).(string)
	bucket := gcs.GetBucket(t, name)
	assert.Equal(t, name, bucket.Name)
}

func TestGCPStorageBucketLifeCycleRules(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/storage/gcs_bucket"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	region := "us-central1"

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":       name,
			"project_id": projectID,
			"region":     region,
			"versioning": true,
			"lifecycle_rules": []map[string]interface{}{
				{
					"action": []map[string]string{
						{
							"type":          "SetStorageClass",
							"storage_class": "NEARLINE",
						},
					},
					"condition": []map[string]interface{}{
						{
							"age":                   60,
							"created_before":        "2018-08-20",
							"with_state":            "LIVE",
							"matches_storage_class": []string{"REGIONAL"},
							"num_newer_versions":    10,
						},
					},
				},
			},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateStorageBucketCreatedWithLifeCycleRule(t, terraformOptions, terraformOutputs)
}

func validateStorageBucketCreatedWithLifeCycleRule(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	var nilSlice []string
	name := out["name"].(string)
	bucket := gcs.GetBucket(t, name)
	assert.Equal(t, name, bucket.Name)

	rule := out["lifecycle_rule"].([]interface{})[0].(map[string]interface{})
	expectedLifeCycleRule := &storage.BucketLifecycleRule{
		Action: &storage.BucketLifecycleRuleAction{
			StorageClass: rule["action"].([]interface{})[0].(map[string]interface{})["storage_class"].(string),
			Type:         rule["action"].([]interface{})[0].(map[string]interface{})["type"].(string),
		},
		Condition: &storage.BucketLifecycleRuleCondition{
			Age:                 int64(rule["condition"].([]interface{})[0].(map[string]interface{})["age"].(float64)),
			CreatedBefore:       rule["condition"].([]interface{})[0].(map[string]interface{})["created_before"].(string),
			MatchesStorageClass: []string{rule["condition"].([]interface{})[0].(map[string]interface{})["matches_storage_class"].([]interface{})[0].(string)},
			NumNewerVersions:    int64(rule["condition"].([]interface{})[0].(map[string]interface{})["num_newer_versions"].(float64)),
		},
		ForceSendFields: nilSlice,
		NullFields:      nilSlice,
	}
	withState := rule["condition"].([]interface{})[0].(map[string]interface{})["with_state"].(string)

	switch withState {
	case "LIVE":
		expectedLifeCycleRule.Condition.IsLive = googleapi.Bool(true)
	case "ARCHIVED":
		expectedLifeCycleRule.Condition.IsLive = googleapi.Bool(false)
	case "ANY", "":
		// This is unnecessary, but set explicitly to nil for readability.
		expectedLifeCycleRule.Condition.IsLive = nil
	default:
		log.Errorf("unexpected value %q for condition.with_state", withState)
		t.FailNow()
	}
	actualRule := bucket.Lifecycle.Rule[0]
	identicalRules := reflect.DeepEqual(expectedLifeCycleRule, actualRule)
	if !identicalRules {
		log.Errorf("actual rule: %#v", actualRule)
		log.Errorf("expected rule: %#v", expectedLifeCycleRule)
		t.FailNow()
	}
	// check lifecycle rules identical
	assert.True(t, identicalRules)

	// check that versioning is enabled
	assert.True(t, bucket.Versioning.Enabled)
}
