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

	compute "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	gcompute "google.golang.org/api/compute/v1"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func TestGCPRouteBasic(t *testing.T) {
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/vpc/routes"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	// default network
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	routeName := fmt.Sprintf("terratest-route-egress-nat-%s", strings.ToLower(random.UniqueId()))

	// custom routes
	customRoutes := []map[string]string{
		{
			"name":              routeName,
			"description":       "route through IGW to access internet",
			"destination_range": "0.0.0.0/0",
			"tags":              "egress-inet",
			"next_hop_internet": "true",
		},
	}
	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":   projectID,
			"network_name": "default",
			"routes":       customRoutes,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateBasicRoute(t, terraformOptions, terraformOutputs, projectID, "default")
}

func validateBasicRoute(t *testing.T, opts *terraform.Options, routes map[string]interface{}, projectID, networkName string) {
	networkLink := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", projectID, networkName)

	routesFound := 0
	var nilSlice []string

	actualRoutesPerNetwork := compute.ListRoutes(t, projectID, networkLink)

	for name, route := range routes {
		tags := []string{}
		route := route.(map[string]interface{})

		for _, t := range route["tags"].([]interface{}) {
			tags = append(tags, t.(string))
		}

		r := &gcompute.Route{
			Description:      route["description"].(string),
			DestRange:        route["dest_range"].(string),
			Name:             name,
			Network:          route["network"].(string),
			Tags:             tags,
			NextHopGateway:   route["next_hop_gateway"].(string),
			NextHopIlb:       "",
			NextHopInstance:  "",
			NextHopIp:        "",
			NextHopNetwork:   "",
			NextHopPeering:   "",
			NextHopVpnTunnel: "",
			Priority:         int64(route["priority"].(float64)),
			SelfLink:         route["self_link"].(string),
			ForceSendFields:  nilSlice,
			NullFields:       nilSlice,
		}

		for _, actualRoute := range actualRoutesPerNetwork {
			// skip if expected route is not this route
			if r.Name != actualRoute.Name {
				continue
			}

			// align attributes
			r.CreationTimestamp = actualRoute.CreationTimestamp
			r.Id = actualRoute.Id
			r.SelfLink = actualRoute.SelfLink
			r.ServerResponse = actualRoute.ServerResponse
			r.Warnings = actualRoute.Warnings
			r.Kind = actualRoute.Kind

			if reflect.DeepEqual(r, actualRoute) {
				t.Logf("Route %s found and match the actual route", r.Name)
				routesFound++
			} else {
				log.Errorf("%s error! actual route and expected route does not match:\n\n%#v\n\n", r.Name, pretty.Diff(actualRoute, r))

				t.Logf("expected route:\n\n %#v\n\n", r)
				t.Logf("actual route:\n\n %#v\n\n", actualRoute)
				t.FailNow()
			}
		}
	}
}

func TestGCPInstanceRoute(t *testing.T) {
	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/network/vpc/routes"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	region := "us-central1"
	zone := "us-central1-a"

	// create network and subnet
	routeName := fmt.Sprintf("terratest-route-egress-nat-%s", strings.ToLower(random.UniqueId()))
	proxyRouteName := fmt.Sprintf("terratest-route-proxy-vm-%s", strings.ToLower(random.UniqueId()))
	networkName := fmt.Sprintf("terratest-routes-vpc-%s", strings.ToLower(random.UniqueId()))
	cidrName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	compute.CreateNetwork(t, projectID, networkName)
	cidr := "10.50.10.0/24"
	subnetName := "subnet-01"
	subnet := compute.CreateSubnetwork(t, projectID, region, networkName, subnetName, cidr, cidrName)

	// create a vm instance with tag app-proxy
	name := "app-proxy"
	compute.CreateVMInstance(t, projectID, zone, name, "f1-micro", subnet.SelfLink, "debian", 10, []string{name})
	// custom routes
	customRoutes := []map[string]string{
		{
			"name":              routeName,
			"description":       "route through IGW to access internet",
			"destination_range": "0.0.0.0/0",
			"tags":              "egress-inet",
			"next_hop_internet": "true",
		}, {
			"name":                   proxyRouteName,
			"description":            "route through proxy to reach app",
			"destination_range":      cidr,
			"tags":                   name,
			"next_hop_instance":      name,
			"next_hop_instance_zone": zone,
		},
	}
	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"project_id":   projectID,
			"network_name": "default",
			"routes":       customRoutes,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	defer compute.DeleteNetwork(t, projectID, networkName)
	defer compute.DeleteSubnetwork(t, projectID, region, subnetName)
	defer compute.DeleteVMInstance(t, projectID, zone, name)

	terraform.InitAndApply(t, terraformOptions)

	// terraform.OutputAll is a map of string interface which is used to validate the actual resources are created
	terraformOutputs := terraform.OutputAll(t, terraformOptions)["output"].(map[string]interface{})

	validateInstanceRoute(t, terraformOptions, terraformOutputs, projectID, "default")
}

func validateInstanceRoute(t *testing.T, opts *terraform.Options, routes map[string]interface{}, projectID, networkName string) {
	networkLink := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", projectID, networkName)
	computeBaseURL := "https://www.googleapis.com/compute/v1"
	routesFound := 0
	var (
		nextHopInstance string
		nilSlice        []string
	)

	actualRoutesPerNetwork := compute.ListRoutes(t, projectID, networkLink)

	// for each of the terraform output routes build a route object
	for name, route := range routes {
		tags := []string{}
		route := route.(map[string]interface{})

		for _, t := range route["tags"].([]interface{}) {
			tags = append(tags, t.(string))
		}

		n := cast.ToString(route["next_hop_instance"])
		if n != "" {
			nextHopInstance = fmt.Sprintf("%s/%s", computeBaseURL, n)
		}

		r := &gcompute.Route{
			Name:             name,
			Description:      cast.ToString(route["description"]),
			DestRange:        cast.ToString(route["dest_range"]),
			Network:          cast.ToString(route["network"]),
			Tags:             tags,
			NextHopGateway:   cast.ToString(route["next_hop_gateway"]),
			NextHopIlb:       cast.ToString(route["next_hop_ilb"]),
			NextHopInstance:  nextHopInstance,
			NextHopIp:        cast.ToString(route["next_hop_ip"]),
			NextHopNetwork:   cast.ToString(route["next_hop_network"]),
			NextHopPeering:   cast.ToString(route["next_hop_peering"]),
			NextHopVpnTunnel: cast.ToString(route["next_hop_vpn_tunneling"]),
			Priority:         int64(cast.ToFloat64(route["priority"])),
			SelfLink:         cast.ToString(route["self_link"]),
			ForceSendFields:  nilSlice,
			NullFields:       nilSlice,
		}

		// check if name match that the outputs and actual are matching
		for _, actualRoute := range actualRoutesPerNetwork {
			// skip if expected route is not this route
			if r.Name != actualRoute.Name {
				continue
			}

			// align attributes
			r.CreationTimestamp = actualRoute.CreationTimestamp
			r.Id = actualRoute.Id
			r.SelfLink = actualRoute.SelfLink
			r.ServerResponse = actualRoute.ServerResponse
			r.Warnings = actualRoute.Warnings
			r.Kind = actualRoute.Kind

			if reflect.DeepEqual(r, actualRoute) {
				t.Logf("Route %s found and match the actual route", r.Name)
				routesFound++
			} else {
				log.Errorf("%s error! actual route and expected route does not match:\n\n%#v\n\n", r.Name, pretty.Diff(actualRoute, r))
				t.FailNow()
			}
		}
	}
}
