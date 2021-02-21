package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/random"
	"os"
	"strings"

	"testing"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	api "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/api"
	iam "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/iam"
	log "github.com/sirupsen/logrus"
	v1 "google.golang.org/api/cloudresourcemanager/v1"
)

// CreateGoogleProvider create provider.tf at the given filepath
func CreateGoogleProvider(t *testing.T, projectID, path string) {
	if err := CreateGoogleProviderE(t, projectID, path); err != nil {
		log.Error(err)
		t.FailNow()
	}
}

// CreateGoogleProviderE create provider.tf at the given filepathor return an error
func CreateGoogleProviderE(t *testing.T, projectID, filePath string) error {
	path := fmt.Sprintf("%s/provider.tf", filePath)
	fmt.Printf("creating provider.tf ffor google-beta in %s", filePath)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file provider.tf: %v", err)
	}
	l, err := f.WriteString("provider \"google-beta\" {project = var.project_id}")
	if err != nil {
		f.Close()
		return fmt.Errorf("error writing the file provider.tf: %v", err)
	}
	fmt.Println(l, "bytes written successfully")

	err = f.Close()
	if err != nil {
		return fmt.Errorf("error closing the file provider.tf: %v", err)
	}
	return nil
}

func setup(t *testing.T) {
	log.Info("Checking required environment variabels...")
	// Create a project
	orgDomain := os.Getenv("ORG_DOMAIN")
	if orgDomain == "" {
		log.Fatal("error org_domain env var is missing, add it with export ORG_DOMAIN=domain and re-run")
	}
	// Get BillingAccount
	billingAccount := os.Getenv("GCP_BILLING_ACCOUNT")
	if billingAccount == "" {
		log.Fatal("You must set billing account from environment variable: export GCP_BILLING_ACCOUNT=\"01A29B-000000-000000\"")
	}

	parent := &v1.ResourceId{
		Id:   strings.Split(iam.GetOrgID(t, orgDomain), "/")[1],
		Type: "organization",
	}

	projectID := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	log.Info("Setup new project and enable api's")
	project := iam.CreateProject(t, projectID, projectID, parent, billingAccount)

	// enable all api's
	services := []string{
		"compute.googleapis.com",
		"servicenetworking.googleapis.com",
		"sql-component.googleapis.com",
		"sqladmin.googleapis.com",
		"redis.googleapis.com",
		"secretmanager.googleapis.com",
		"dns.googleapis.com",
	}

	api.EnableAPI(t, project, services)

	os.Setenv("GOOGLE_PROJECT", projectID)
}

func teardown(t *testing.T) {
	// delete project
	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)

	if !strings.HasPrefix(projectID, "terratest") {
		log.Warnf("Refusing to delete project %s as it does nto start with terratest prefix", projectID)
		return
	}
	iam.DeleteProject(t, projectID)
}
