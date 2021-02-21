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

	"github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	"github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/dns"
	"github.com/stretchr/testify/assert"
)

func TestGCPManagedSSL(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/managed_ssl_certificate"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	// create Cloud DNS record
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	zoneName := fmt.Sprintf("terratest-dns-%s", strings.ToLower(random.UniqueId()))
	m := dns.CreateManagedZone(t, projectID, zoneName, fmt.Sprintf("%s.com.", zoneName), "managed by terrtest")

	recordName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	fqdn := fmt.Sprintf("%s.%s", recordName, m.DnsName)
	certName := fmt.Sprintf("terratest-ip-%s", strings.ToLower(random.UniqueId()))
	dns.RecordSetChange(t, "add", projectID, zoneName, "A", fqdn, "10.0.0.5", 30)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":         certName,
			"project_id":   projectID,
			"cert_domains": []string{fqdn},
		},
	}

	defer dns.DeleteManagedZone(t, projectID, zoneName)
	defer dns.RecordSetChange(t, "delete", projectID, zoneName, "A", fqdn, "10.0.0.5", 30)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateManagedCertCreated(t, terraformOptions, terraformOutputs)
}

func validateManagedCertCreated(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	name := out["name"].(string)
	domains := out["managed"].([]interface{})[0].(map[string]interface{})["domains"].([]interface{})

	cert := compute.GetSSLCertificate(t, projectID, name)
	assert.Equal(t, name, cert.Name)
	reflect.DeepEqual(domains, cert.Managed.Domains)
	assert.Equal(t, "PROVISIONING", cert.Managed.Status)
}
