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
	"github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/dns"
	"github.com/stretchr/testify/assert"
	gdns "google.golang.org/api/dns/v1"
)

type simpleRecordSet struct {
	projectID   string
	recordName  string
	recordType  string
	recordTTL   int64
	recordValue []string
	managedZone string
}

func TestGCPAddDNSRecordDefaults(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/dns/record_set"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	// create Cloud DNS record
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	zoneName := fmt.Sprintf("terratest-dns-%s", strings.ToLower(random.UniqueId()))
	dns.CreateManagedZone(t, projectID, zoneName, fmt.Sprintf("%s.com.", zoneName), "managed by terrtest")
	recordName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":   projectID,
			"record_name":  recordName,
			"record_data":  []string{"10.0.0.5"},
			"managed_zone": zoneName,
		},
	}

	defer dns.DeleteManagedZone(t, projectID, zoneName)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	rs := &simpleRecordSet{
		projectID:   terraformOutputs["project"].(string),
		recordName:  terraformOutputs["name"].(string),
		recordType:  terraformOutputs["type"].(string),
		recordTTL:   int64(terraformOutputs["ttl"].(float64)),
		recordValue: []string{terraformOutputs["rrdatas"].([]interface{})[0].(string)},
		managedZone: terraformOutputs["managed_zone"].(string),
	}
	// // We have a [string]struct that defines out test cases.
	// the struct holds test in and out (expectations).
	testCases := []struct {
		name      string
		functions func(*testing.T, *terraform.Options, *simpleRecordSet)
	}{
		{
			name:      "Validate dns record created",
			functions: validateDNSRecordCreated,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.functions(t, terraformOptions, rs)
		})
	}
}

func validateDNSRecordCreated(t *testing.T, opts *terraform.Options, rs *simpleRecordSet) {
	var (
		found           bool
		forceSendFields []string
		nullFields      []string
	)

	recordSets := dns.ListRecordSets(t, rs.projectID, rs.managedZone)

	recordSet := &gdns.ResourceRecordSet{
		Kind:             "dns#resourceRecordSet",
		Name:             rs.recordName,
		Rrdatas:          rs.recordValue,
		SignatureRrdatas: []string{},
		Ttl:              rs.recordTTL,
		Type:             rs.recordType,
		ForceSendFields:  forceSendFields,
		NullFields:       nullFields,
	}

	for _, r := range recordSets {
		if reflect.DeepEqual(r, recordSet) {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestGCPAddTxtRecord(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/dns/record_set"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	// create Cloud DNS record
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	zoneName := fmt.Sprintf("terratest-dns-%s", strings.ToLower(random.UniqueId()))
	dns.CreateManagedZone(t, projectID, zoneName, fmt.Sprintf("%s.com.", zoneName), "managed by terrtest")
	recordName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":   projectID,
			"record_name":  recordName,
			"record_type":  "TXT",
			"record_ttl":   "60",
			"record_data":  []string{"terratest placeholder text"},
			"managed_zone": zoneName,
		},
	}

	defer dns.DeleteManagedZone(t, projectID, zoneName)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})
	rs := &simpleRecordSet{
		projectID:   terraformOutputs["project"].(string),
		recordName:  terraformOutputs["name"].(string),
		recordType:  terraformOutputs["type"].(string),
		recordTTL:   int64(terraformOutputs["ttl"].(float64)),
		recordValue: []string{terraformOutputs["rrdatas"].([]interface{})[0].(string)},
		managedZone: terraformOutputs["managed_zone"].(string),
	}
	validateDNSRecordCreated(t, terraformOptions, rs)
}

func TestDNSPublicZone(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/dns/dns_zone"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	zoneName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	dnsName := fmt.Sprintf("%s.com", zoneName)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id": projectID,
			"type":       "public",
			"name":       zoneName,
			"dns_name":   dnsName,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})
	nameServers := []string{}
	for _, ns := range terraformOutputs["name_servers"].([]interface{}) {
		nameServers = append(nameServers, ns.(string))
	}

	rs := &simpleRecordSet{
		projectID:   projectID,
		recordName:  terraformOutputs["dns_name"].(string),
		recordType:  "NS",
		recordTTL:   21600,
		recordValue: nameServers,
		managedZone: terraformOutputs["name"].(string),
	}

	validateDNSRecordCreated(t, terraformOptions, rs)
}

func TestDNSPrivateZone(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/dns/dns_zone"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	// create Cloud DNS record
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	zoneName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	dnsName := fmt.Sprintf("%s.com", zoneName)
	network := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", projectID, "default")

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":                         projectID,
			"name":                               zoneName,
			"dns_name":                           dnsName,
			"type":                               "private",
			"private_visibility_config_networks": []string{network},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	nameServers := []string{}
	for _, ns := range terraformOutputs["name_servers"].([]interface{}) {
		nameServers = append(nameServers, ns.(string))
	}

	rs := &simpleRecordSet{
		projectID:   projectID,
		recordName:  terraformOutputs["dns_name"].(string),
		recordType:  "NS",
		recordTTL:   21600,
		recordValue: nameServers,
		managedZone: terraformOutputs["name"].(string),
	}

	validateDNSRecordCreated(t, terraformOptions, rs)
}
