package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	compute "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestGCPComputeHTTPHealthCheck(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/health_check"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":         projectID,
			"healthcheck_name":   name,
			"healthcheck_type":   "http",
			"request_path":       "/app",
			"response":           "O.K",
			"port":               8080,
			"port_specification": "USE_FIXED_PORT",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateHealthCheckCreated(t, terraformOptions, terraformOutputs)
}

func TestGCPComputeTCPHealthCheck(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/health_check"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	name := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":         projectID,
			"healthcheck_name":   name,
			"healthcheck_type":   "tcp",
			"request":            `{"test": "123"}`,
			"response":           `{"resp": {"all works!"}}`,
			"port_name":          "listener",
			"port_specification": "USE_NAMED_PORT",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateHealthCheckCreated(t, terraformOptions, terraformOutputs)
}

func validateHealthCheckCreated(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	projectID := out["project"].(string)
	name := out["name"].(string)
	actualHc := compute.GetHealthCheck(t, projectID, name)

	var (
		expectedPort, actualPort                                                                                              int64
		expectedHost, expectedPortName, expectedPortSpecification, expectedProxyHeader, expectedRequestPath, expectedResponse string
		actualHost, actualPortName, actualPortSpecification, actualProxyHeader, actualRequestPath, actualResponse             string
	)

	switch opts.Vars["healthcheck_type"].(string) {
	case "http":
		httpHealthCheck := out["http_health_check"].([]interface{})[0].(map[string]interface{})
		expectedPort = cast.ToInt64(httpHealthCheck["port"])
		expectedHost = cast.ToString(httpHealthCheck["host"])
		expectedPortName = cast.ToString(httpHealthCheck["port_name"])
		expectedPortSpecification = cast.ToString(httpHealthCheck["port_specification"])
		expectedProxyHeader = cast.ToString(httpHealthCheck["proxy_header"])
		expectedRequestPath = cast.ToString(httpHealthCheck["request_path"])
		expectedResponse = cast.ToString(httpHealthCheck["response"])

		actualPort = actualHc.HttpHealthCheck.Port
		actualHost = actualHc.HttpHealthCheck.Host
		actualPortName = actualHc.HttpHealthCheck.PortName
		actualPortSpecification = actualHc.HttpHealthCheck.PortSpecification
		actualProxyHeader = actualHc.HttpHealthCheck.ProxyHeader
		actualRequestPath = actualHc.HttpHealthCheck.RequestPath
		actualResponse = actualHc.HttpHealthCheck.Response

	case "tcp":
		httpHealthCheck := out["tcp_health_check"].([]interface{})[0].(map[string]interface{})
		expectedPort = cast.ToInt64(httpHealthCheck["port"])
		expectedPortName = cast.ToString(httpHealthCheck["port_name"])
		expectedPortSpecification = cast.ToString(httpHealthCheck["port_specification"])
		expectedProxyHeader = cast.ToString(httpHealthCheck["proxy_header"])
		expectedRequestPath = cast.ToString(httpHealthCheck["request_path"])
		expectedResponse = cast.ToString(httpHealthCheck["response"])

		actualPort = actualHc.TcpHealthCheck.Port
		actualPortName = actualHc.TcpHealthCheck.PortName
		actualPortSpecification = actualHc.TcpHealthCheck.PortSpecification
		actualProxyHeader = actualHc.TcpHealthCheck.ProxyHeader
		actualResponse = actualHc.TcpHealthCheck.Response
	}

	assert.Equal(t, expectedPort, actualPort)
	assert.Equal(t, expectedHost, actualHost)
	assert.Equal(t, expectedPortName, actualPortName)
	assert.Equal(t, expectedPortSpecification, actualPortSpecification)
	assert.Equal(t, expectedProxyHeader, actualProxyHeader)
	assert.Equal(t, expectedRequestPath, actualRequestPath)
	assert.Equal(t, expectedResponse, actualResponse)
}
