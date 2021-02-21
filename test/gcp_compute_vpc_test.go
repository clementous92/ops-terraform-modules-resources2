package test

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"

	"strings"
	"testing"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/kr/pretty"
	api "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/api"
	compute "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	iam "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/iam"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	v1 "google.golang.org/api/cloudresourcemanager/v1"
	gcompute "google.golang.org/api/compute/v1"
)

func TestGCPEnableHostProject(t *testing.T) {
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/vpc/host_project"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	hostProjectID := fmt.Sprintf("terratest-host-project-%s", strings.ToLower(random.UniqueId()))

	orgDomain := os.Getenv("ORG_DOMAIN")
	if orgDomain == "" {
		log.Error("error org_domain env var is missing, add it with export ORG_DOMAIN=domain")
		t.FailNow()
	}

	parent := &v1.ResourceId{
		Id:   strings.Split(iam.GetOrgID(t, orgDomain), "/")[1],
		Type: "organization",
	}

	// Get BillingAccount
	billingAccount := os.Getenv("GCP_BILLING_ACCOUNT")
	if billingAccount == "" {
		t.Fatal("You must set billing account from environment variable: export GCP_BILLING_ACCOUNT=\"01A29B-000000-000000\"")
	}

	// create a project to be the host
	project := iam.CreateProject(t, hostProjectID, "Terratest Host Project", parent, billingAccount)

	// enable the compute engine api
	api.EnableAPI(t, project, []string{"compute.googleapis.com"})

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"host_project": hostProjectID,
		},
	}

	defer iam.DeleteProject(t, hostProjectID)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// terraformOutputs is a map of string string which is used to validate the actual resources are created
	terraformOutputs := terraform.OutputAll(t, terraformOptions)

	validateProjectIsAHostProject(t, terraformOptions, terraformOutputs)
}

func validateProjectIsAHostProject(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	hostProjectID := out["output"].(map[string]interface{})["id"].(string)
	proj := compute.GetProject(t, hostProjectID)
	assert.Equal(t, "HOST", proj.XpnProjectStatus)
}

func TestVPCStandardNetwork(t *testing.T) {
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/vpc/network"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":   projectID,
			"network_name": fmt.Sprintf("terratest-network-1-%s", strings.ToLower(random.UniqueId())),
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// terraformOutputs is a map of string string which is used to validate the actual resources are created
	terraformOutputs := terraform.OutputAll(t, terraformOptions)

	validateVPCStandardNetworkCreated(t, terraformOptions, terraformOutputs)
}

func validateVPCStandardNetworkCreated(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	output := out["output"].(map[string]interface{})
	name := output["name"].(string)
	project := output["project"].(string)
	selfLink := output["self_link"].(string)
	vpc := compute.GetNetwork(t, project, name)
	assert.Equal(t, selfLink, vpc.SelfLink)
}

func TestVPCSubnetsCreation(t *testing.T) {
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/vpc/subnets"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	// create VPC network
	networkName := fmt.Sprintf("terratest-vpc-network-1-%s", strings.ToLower(random.UniqueId()))
	compute.CreateNetwork(t, projectID, networkName)
	subnetBase := "10.10"
	cidrRangeMin := 0
	cidrRangeMax := 255
	subnetMask := "24"
	// subnets configuration
	subnets := []map[string]interface{}{
		{
			"subnet_name":   "subnet-01",
			"subnet_ip":     fmt.Sprintf("%s.%d.0/%s", subnetBase, rand.Intn(cidrRangeMax-cidrRangeMin)+cidrRangeMin, subnetMask),
			"subnet_region": "us-west1",
		}, {
			"subnet_name":           "subnet-02",
			"subnet_ip":             fmt.Sprintf("%s.%d.0/%s", subnetBase, rand.Intn(cidrRangeMax-cidrRangeMin)+cidrRangeMin, subnetMask),
			"subnet_region":         "us-west1",
			"subnet_private_access": true,
			"subnet_flow_logs":      true,
			"description":           "This subnet has a description",
		}, {
			"subnet_name":               "subnet-03",
			"subnet_ip":                 fmt.Sprintf("%s.%d.0/%s", subnetBase, rand.Intn(cidrRangeMax-cidrRangeMin)+cidrRangeMin, subnetMask),
			"subnet_region":             "us-west1",
			"subnet_private_access":     true,
			"subnet_flow_logs":          true,
			"subnet_flow_logs_interval": "INTERVAL_10_MIN",
			"subnet_flow_logs_sampling": "0.7",
			"subnet_flow_logs_metadata": "INCLUDE_ALL_METADATA",
		},
	}

	secondaryRanges := map[string]interface{}{
		"subnet-01": []map[string]string{
			{
				"range_name":    "subnet-01-secondary-01",
				"ip_cidr_range": fmt.Sprintf("192.168.%d.0/%s", rand.Intn(cidrRangeMax-cidrRangeMin)+cidrRangeMin, subnetMask),
			},
		},
		"subnet-02": []map[string]string{
			{
				"range_name":    "subnet-02-secondary-01",
				"ip_cidr_range": fmt.Sprintf("192.168.%d.0/%s", rand.Intn(cidrRangeMax-cidrRangeMin)+cidrRangeMin, subnetMask),
			}, {
				"range_name":    "subnet-02-secondary-02",
				"ip_cidr_range": fmt.Sprintf("192.168.%d.0/%s", rand.Intn(cidrRangeMax-cidrRangeMin)+cidrRangeMin, subnetMask),
			},
		},
	}

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":       projectID,
			"network_name":     networkName,
			"subnets":          subnets,
			"secondary_ranges": secondaryRanges,
		},
	}

	defer compute.DeleteNetwork(t, projectID, networkName)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// terraformOutputs is a map of string string which is used to validate the actual resources are created
	terraformOutputs := terraform.OutputAll(t, terraformOptions)

	validateVPCSubnetsCreation(t, terraformOptions, terraformOutputs, projectID)
}

func validateVPCSubnetsCreation(t *testing.T, opts *terraform.Options, out map[string]interface{}, project string) {
	subnets := out["output"].(map[string]interface{})

	var nilSlice []string
	var nilString string

	for _, subnet := range subnets {
		info := subnet.(map[string]interface{})

		// Get the actual subnet
		actualSubnet := compute.GetSubnetwork(t, project, info["region"].(string), info["name"].(string))

		// build expected subnet data
		expectedSubnet := &gcompute.Subnetwork{
			CreationTimestamp:       actualSubnet.CreationTimestamp,
			Description:             info["description"].(string),
			EnableFlowLogs:          false,
			Fingerprint:             actualSubnet.Fingerprint,
			GatewayAddress:          info["gateway_address"].(string),
			Id:                      actualSubnet.Id,
			IpCidrRange:             info["ip_cidr_range"].(string),
			Ipv6CidrRange:           nilString,
			Kind:                    "compute#subnetwork",
			Name:                    info["name"].(string),
			Network:                 info["network"].(string),
			PrivateIpGoogleAccess:   info["private_ip_google_access"].(bool),
			PrivateIpv6GoogleAccess: actualSubnet.PrivateIpv6GoogleAccess,
			Purpose:                 actualSubnet.Purpose,
			Region:                  fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/regions/%s", project, info["region"].(string)),
			Role:                    actualSubnet.Role,
			SelfLink:                info["self_link"].(string),
			State:                   actualSubnet.State,
			ServerResponse:          actualSubnet.ServerResponse,
			ForceSendFields:         nilSlice,
			NullFields:              nilSlice,
		}

		if len(info["log_config"].([]interface{})) > 0 {
			expectedSubnet.EnableFlowLogs = true
		}

		expectedSubnet.LogConfig = &gcompute.SubnetworkLogConfig{
			Enable:              actualSubnet.LogConfig.Enable,
			FilterExpr:          actualSubnet.LogConfig.FilterExpr,
			MetadataFields:      actualSubnet.LogConfig.MetadataFields,
			ForceSendFields:     nilSlice,
			NullFields:          nilSlice,
			AggregationInterval: "",
			FlowSampling:        float64(0),
		}

		if actualSubnet.LogConfig.Enable {
			expectedSubnet.LogConfig.AggregationInterval = info["log_config"].([]interface{})[0].(map[string]interface{})["aggregation_interval"].(string)
			expectedSubnet.LogConfig.FlowSampling = float64(info["log_config"].([]interface{})[0].(map[string]interface{})["flow_sampling"].(float64))
			expectedSubnet.LogConfig.Metadata = info["log_config"].([]interface{})[0].(map[string]interface{})["metadata"].(string)
		}

		if actualSubnet.SecondaryIpRanges != nil {
			expectedSubnet.SecondaryIpRanges = []*gcompute.SubnetworkSecondaryRange{}

			for _, subnet := range info["secondary_ip_range"].([]interface{}) {
				expectedSubnet.SecondaryIpRanges = append(
					expectedSubnet.SecondaryIpRanges,
					&gcompute.SubnetworkSecondaryRange{
						IpCidrRange:     subnet.(map[string]interface{})["ip_cidr_range"].(string),
						RangeName:       subnet.(map[string]interface{})["range_name"].(string),
						ForceSendFields: nilSlice,
						NullFields:      nilSlice,
					},
				)

			}
		}

		identicalRules := reflect.DeepEqual(expectedSubnet, actualSubnet)
		if !identicalRules {
			log.Errorf("%s error! actual subnet and expected does not match:\n\n%#v\n\n", info["name"].(string), pretty.Diff(actualSubnet, expectedSubnet))
			t.FailNow()
		}
		assert.True(t, reflect.DeepEqual(expectedSubnet, actualSubnet))
	}
}

func TestGCPServiceProjects(t *testing.T) {

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/vpc/service_project"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	hostProjectID := fmt.Sprintf("terratest-host-project-%s", strings.ToLower(random.UniqueId()))

	orgDomain := os.Getenv("ORG_DOMAIN")
	if orgDomain == "" {
		log.Error("error org_domain env var is missing, add it with export ORG_DOMAIN=domain")
		t.FailNow()
	}

	parent := &v1.ResourceId{
		Id:   strings.Split(iam.GetOrgID(t, orgDomain), "/")[1],
		Type: "organization",
	}

	// Get BillingAccount
	billingAccount := os.Getenv("GCP_BILLING_ACCOUNT")
	if billingAccount == "" {
		t.Fatal("You must set billing account from environment variable: export GCP_BILLING_ACCOUNT=\"01A29B-000000-000000\"")
	}

	// create a project to be the host
	hostProject := iam.CreateProject(t, hostProjectID, "Terratest Host Project", parent, billingAccount)

	// create a service project
	serviceProjectID := fmt.Sprintf("terratest-service-proj-%s", strings.ToLower(random.UniqueId()))
	serviceProject := iam.CreateProject(t, serviceProjectID, "Terratest Service 1 Project", parent, billingAccount)

	// enable the compute engine api for both projets
	api.EnableAPI(t, hostProject, []string{"compute.googleapis.com"})
	api.EnableAPI(t, serviceProject, []string{"compute.googleapis.com"})

	// enable vpc host project
	compute.EnableXPNHost(t, hostProjectID)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"host_project": hostProjectID,
			"project_id":   serviceProjectID,
		},
	}

	// Deferred function calls are pushed onto a stack. When a function returns, its deferred calls are executed in last-in-first-out order.
	defer iam.DeleteProject(t, hostProjectID)
	defer compute.DisableXPNHost(t, hostProjectID)
	defer iam.DeleteProject(t, serviceProjectID)
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// terraformOutputs is a map of string string which is used to validate the actual resources are created
	terraformOutputs := terraform.OutputAll(t, terraformOptions)

	validateHostProjectWithServiceProject(t, terraformOptions, terraformOutputs)
}

func validateHostProjectWithServiceProject(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	output := out["output"].(map[string]interface{})

	hostProjectID := output["host_project"].(string)
	serviceProjectID := output["service_project"].(string)
	resources := compute.GetXPNResources(t, hostProjectID)
	assert.Equal(t, serviceProjectID, resources.Resources[0].Id)
	assert.Equal(t, "PROJECT", resources.Resources[0].Type)
}
