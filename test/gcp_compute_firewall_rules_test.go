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
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	gcompute "google.golang.org/api/compute/v1"
)

func TestGCPComputeFirewallRulesCustomAllow(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/vpc/firewall"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	network := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", projectID, "default")
	rule1 := fmt.Sprintf("terratest-fw-rule-1-%s", strings.ToLower(random.UniqueId()))
	rule2 := fmt.Sprintf("terratest-fw-rule-2-%s", strings.ToLower(random.UniqueId()))

	customRules := map[string]interface{}{
		rule1: map[string]interface{}{
			"description":          "allow ssh for tags",
			"direction":            "INGRESS",
			"action":               "allow",
			"sources":              []string{},
			"ranges":               []string{"0.0.0.0/0"},
			"targets":              []string{"terratest-ssh-tag"},
			"use_service_accounts": false,
			"rules": []map[string]interface{}{
				{
					"protocol": "tcp",
					"ports":    []string{"22"},
				},
			},
			"extra_attributes": map[string]interface{}{
				"disabled": false,
				"priority": int64(1000),
			},
		},
		rule2: map[string]interface{}{
			"description":          "allow connection to postgres",
			"direction":            "INGRESS",
			"action":               "allow",
			"sources":              []string{},
			"ranges":               []string{"0.0.0.0/0"},
			"targets":              []string{"terratest-postgres-tag"},
			"use_service_accounts": false,
			"rules": []map[string]interface{}{
				{
					"protocol": "tcp",
					"ports":    []string{"5432"},
				},
			},
			"extra_attributes": map[string]interface{}{
				"disabled": false,
				"priority": int64(990),
			},
		},
	}
	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":   projectID,
			"network":      network,
			"custom_rules": customRules,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	validateFirewallRuleCreated(t, terraformOptions, projectID, network, customRules)
}

func validateFirewallRuleCreated(t *testing.T, opts *terraform.Options, projectID, network string, customRules map[string]interface{}) {
	fwRules := compute.ListFirewallRules(t, projectID)

	var nilSlice []string
	rulesFound := 0

	for name, rule := range customRules {
		rule := rule.(map[string]interface{})

		r := &gcompute.Firewall{
			Name:    name,
			Network: network,
			Allowed: []*gcompute.FirewallAllowed{
				{
					IPProtocol: rule["rules"].([]map[string]interface{})[0]["protocol"].(string),
					Ports:      rule["rules"].([]map[string]interface{})[0]["ports"].([]string),
				},
			},
			Description:           rule["description"].(string),
			DestinationRanges:     nilSlice,
			Direction:             rule["direction"].(string),
			Disabled:              rule["extra_attributes"].(map[string]interface{})["disabled"].(bool),
			Priority:              rule["extra_attributes"].(map[string]interface{})["priority"].(int64),
			SourceRanges:          rule["ranges"].([]string),
			SourceServiceAccounts: nilSlice,
			SourceTags:            nilSlice,
			TargetServiceAccounts: nilSlice,
			TargetTags:            rule["targets"].([]string),
			ForceSendFields:       nilSlice,
			NullFields:            nilSlice,
			LogConfig: &gcompute.FirewallLogConfig{
				Enable:          false,
				ForceSendFields: nilSlice,
				NullFields:      nilSlice,
			},
		}

		// t.Logf("rule: %#v", r)
		for _, actualRule := range fwRules {
			if r.Name != actualRule.Name {
				continue
			}

			// align attributes
			r.CreationTimestamp = actualRule.CreationTimestamp
			r.Id = actualRule.Id
			r.SelfLink = actualRule.SelfLink
			r.ServerResponse = actualRule.ServerResponse
			r.Kind = actualRule.Kind

			if reflect.DeepEqual(r, actualRule) {
				t.Logf("Rule %s found and match the actual firewall rules", r.Name)
				rulesFound++
			} else {
				log.Error("expected rule and actual rule does not match")
				t.Logf("expectedRule:\n\n %#v\n\n", r)
				t.Logf("actualRule:\n\n %#v\n\n", actualRule)
				t.FailNow()
			}
		}
	}
	assert.Equal(t, 2, rulesFound)
}
