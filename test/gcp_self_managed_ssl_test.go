package test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"net"
	"strings"
	"testing"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	log "github.com/sirupsen/logrus"

	"github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	ssl "github.com/rocketlawyer/ops-terraform-terratest/modules/ssl"
	"github.com/stretchr/testify/assert"
)

func TestGCPRegionalSelfManagedSSL(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/self_managed_ssl_certificate"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-ip-%s", strings.ToLower(random.UniqueId()))

	createSelfSignedCert(t, tempTestFolder)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":                  name,
			"project_id":            projectID,
			"region":                "us-central1",
			"private_key_file_path": filepath.Join(tempTestFolder, "cert.key"),
			"cert_file_path":        filepath.Join(tempTestFolder, "cert.pem"),
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateSelfManagedCertCreated(t, terraformOptions, terraformOutputs)
}

func TestGCPGlobalSelfManagedSSL(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/self_managed_ssl_certificate"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-ip-%s", strings.ToLower(random.UniqueId()))

	createSelfSignedCert(t, tempTestFolder)
	fmt.Printf("cert key: %s\n\n", filepath.Join(tempTestFolder, "cert.key"))
	fmt.Printf("cert pem: %s\n\n", filepath.Join(tempTestFolder, "cert.pem"))

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"name":                  name,
			"project_id":            projectID,
			"private_key_file_path": filepath.Join(tempTestFolder, "cert.key"),
			"cert_file_path":        filepath.Join(tempTestFolder, "cert.pem"),
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateSelfManagedCertCreated(t, terraformOptions, terraformOutputs)
}

func validateSelfManagedCertCreated(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	var region []string
	if val, ok := out["region"]; ok {
		region = []string{val.(string)}
	}

	projectID := out["project"].(string)
	name := out["name"].(string)
	selflink := out["self_link"].(string)

	cert := compute.GetSSLCertificate(t, projectID, name, region...)
	assert.Equal(t, name, cert.Name)
	assert.Equal(t, selflink, cert.SelfLink)
	assert.Nil(t, nil, cert.Region)
}

func createSelfSignedCert(t *testing.T, path string) {
	san := []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback}
	ca, caErr := ssl.CertTemplate("RocketLawyer Inc CA.", "Terraform RL CA", true, nil)
	if caErr != nil {
		log.Errorf("error creating cert template: %v", caErr)
		t.FailNow()
	}

	caPK := ssl.CreatePrivateKey(t)
	_, caPem := ssl.CreateCert(t, ca, ca, &caPK.PublicKey, caPK)

	// Create and sign the cert (self signed) using the root CA)
	certPK := ssl.CreatePrivateKey(t)
	cert, err := ssl.CertTemplate("RocketLawyer Inc.", "Terraform RL", false, san)
	if err != nil {
		log.Errorf("error creating cert template: %v", err)
		t.FailNow()
	}

	_, certPem := ssl.CreateCert(t, cert, ca, &certPK.PublicKey, caPK)

	WriteFile(t, path, "ca.pem", string(caPem))
	WriteFile(t, path, "ca.key", ssl.ExportPrivateKeyToPem(caPK))
	WriteFile(t, path, "cert.pem", string(certPem))
	WriteFile(t, path, "cert.key", ssl.ExportPrivateKeyToPem(certPK))
}

func WriteFile(t *testing.T, path, filename, data string) {
	err := WriteFileE(t, path, filename, data)

	if err != nil {
		log.Error(err)
		t.FailNow()
	}
}

func WriteFileE(t *testing.T, path, filename, data string) error {
	file, cErr := os.Create(filepath.Join(path, filename))
	if cErr != nil {
		return fmt.Errorf("error creating file %s: %v", filename, cErr)
	}

	defer file.Close()

	_, wErr := io.WriteString(file, data)
	if wErr != nil {
		return fmt.Errorf("error writing to file %s: %v", filename, wErr)
	}
	return file.Sync()
}
