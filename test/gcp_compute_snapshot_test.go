package test

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	terratestGCP "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/gruntwork-io/terratest/modules/test-structure"
	compute "github.com/rocketlawyer/ops-terraform-terratest/modules/gcp/compute"
	"github.com/spf13/cast"
	gcompute "google.golang.org/api/compute/v1"

	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
)

func TestGCPComputeSnapshoptResourcePolicyDailyAsDefault(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/disk_schedule_snapshot_policy"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	policyPrefixName := "daily-snapshot-policy"
	policyName := fmt.Sprintf("terratest-%s-%s", policyPrefixName, strings.ToLower(random.UniqueId()))
	region := "us-central1"

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"policy_name": policyName,
			"project_id":  projectID,
			"region":      region,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutput := terraform.OutputAll(t, terraformOptions)
	terraformOutput = terraformOutput["output"].(map[string]interface{})

	validateSnapshotResourcePolicyIsCorrect(t, terraformOptions, terraformOutput)
}

func TestGCPComputeSnapshoptResourcePolicyHourly(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/disk_schedule_snapshot_policy"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	policyPrefixName := "daily-snapshot-policy"
	policyName := fmt.Sprintf("terratest-%s-%s", policyPrefixName, strings.ToLower(random.UniqueId()))
	region := "us-central1"

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"policy_name":    policyName,
			"project_id":     projectID,
			"region":         region,
			"hours_in_cycle": 12,
			"start_time":     "10:00",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutput := terraform.OutputAll(t, terraformOptions)
	terraformOutput = terraformOutput["output"].(map[string]interface{})

	validateSnapshotResourcePolicyIsCorrect(t, terraformOptions, terraformOutput)
}

func TestGCPComputeSnapshoptResourcePolicyWeekly(t *testing.T) {
	t.Parallel()

	rootFolder := ".."
	terraformFolderRelativeToRoot := "gcp/compute/disk_schedule_snapshot_policy"
	tempTestFolder := utils.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

	if os.Getenv("NO_TEMP_FOLDER") != "" {
		tempTestFolder = "../"
	}

	projectID := terratestGCP.GetGoogleProjectIDFromEnvVar(t)
	policyPrefixName := "daily-snapshot-policy"
	policyName := fmt.Sprintf("terratest-%s-%s", policyPrefixName, strings.ToLower(random.UniqueId()))
	region := "us-central1"

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,

		Vars: map[string]interface{}{
			"policy_name": policyName,
			"project_id":  projectID,
			"region":      region,
			"day_of_weeks": []map[string]string{
				{"day": "FRIDAY", "start_time": "15:00"},
				{"day": "MONDAY", "start_time": "20:00"},
				{"start_time": "10:00", "day": "SUNDAY"},
			},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOutput := terraform.OutputAll(t, terraformOptions)
	terraformOutput = terraformOutput["output"].(map[string]interface{})

	validateSnapshotResourcePolicyIsCorrect(t, terraformOptions, terraformOutput)
}

// validates snapshot policy is correct for either daily, hourly, or weekly
func validateSnapshotResourcePolicyIsCorrect(t *testing.T, opts *terraform.Options, out map[string]interface{}) {
	resourcePolicyName := out["name"].(string)
	projectID := out["project"].(string)
	regionURL, err := url.Parse(out["region"].(string))
	if err != nil {
		log.Errorf("error getting ergion from policy URL %s: %v", regionURL, err)
	}
	region := path.Base(regionURL.Path)

	expectedPolicy := gcompute.ResourcePolicy{
		SnapshotSchedulePolicy: &gcompute.ResourcePolicySnapshotSchedulePolicy{
			RetentionPolicy: &gcompute.ResourcePolicySnapshotSchedulePolicyRetentionPolicy{},
			Schedule: &gcompute.ResourcePolicySnapshotSchedulePolicySchedule{
				DailySchedule:  &gcompute.ResourcePolicyDailyCycle{},
				HourlySchedule: &gcompute.ResourcePolicyHourlyCycle{},
				WeeklySchedule: &gcompute.ResourcePolicyWeeklyCycle{},
			},
			SnapshotProperties: &gcompute.ResourcePolicySnapshotSchedulePolicySnapshotProperties{},
		},
	}
	actualPolicy := compute.GetResourcePolicy(t, projectID, region, resourcePolicyName)

	// snapshot_schedule_policy
	snapshotPolicyOut := out["snapshot_schedule_policy"].([]interface{})[0].(map[string]interface{})

	// retention_policy
	retentionPolicyOut := snapshotPolicyOut["retention_policy"].([]interface{})[0].(map[string]interface{})
	maxRetentionDays := int64(retentionPolicyOut["max_retention_days"].(float64))
	onSourceDelete := retentionPolicyOut["on_source_disk_delete"].(string)

	// schedule
	schedule := snapshotPolicyOut["schedule"].([]interface{})[0].(map[string]interface{})
	dailySchedule := schedule["daily_schedule"].([]interface{})
	hourlySchedule := schedule["hourly_schedule"].([]interface{})
	weeklySchedule := schedule["weekly_schedule"].([]interface{})

	// snapshot properties
	prop := snapshotPolicyOut["snapshot_properties"].([]interface{})[0].(map[string]interface{})
	guestFlush := prop["guest_flush"].(bool)
	labels := cast.ToStringMapString(prop["labels"])
	storageLocations := cast.ToStringSlice(prop["storage_locations"])

	var (
		daysInCycle  int64
		hoursInCycle int64
		startTime    string
		dayOfWeeks   []*gcompute.ResourcePolicyWeeklyCycleDayOfWeek
	)

	if len(dailySchedule) > 0 {
		attr := dailySchedule[0].(map[string]interface{})
		daysInCycle = int64(attr["days_in_cycle"].(float64))
		startTime = attr["start_time"].(string)
	}

	if len(hourlySchedule) > 0 {
		attr := hourlySchedule[0].(map[string]interface{})
		hoursInCycle = int64(attr["hours_in_cycle"].(float64))
		startTime = attr["start_time"].(string)
	}

	if len(weeklySchedule) > 0 {
		// assuming all days have the same duration
		duration := actualPolicy.SnapshotSchedulePolicy.Schedule.WeeklySchedule.DayOfWeeks[0].Duration

		for _, day := range weeklySchedule[0].(map[string]interface{})["day_of_weeks"].([]interface{}) {
			d := day.(map[string]interface{})
			dayBlock := &gcompute.ResourcePolicyWeeklyCycleDayOfWeek{
				Day:       d["day"].(string),
				StartTime: d["start_time"].(string),
				Duration:  duration,
			}
			dayOfWeeks = append(dayOfWeeks, dayBlock)
		}
	}

	expectedPolicy.SnapshotSchedulePolicy.RetentionPolicy.MaxRetentionDays = maxRetentionDays
	expectedPolicy.SnapshotSchedulePolicy.RetentionPolicy.OnSourceDiskDelete = onSourceDelete

	// Daily
	if actualPolicy.SnapshotSchedulePolicy.Schedule.DailySchedule != nil {
		expectedPolicy.SnapshotSchedulePolicy.Schedule.DailySchedule.DaysInCycle = daysInCycle
		expectedPolicy.SnapshotSchedulePolicy.Schedule.DailySchedule.StartTime = startTime
		expectedPolicy.SnapshotSchedulePolicy.Schedule.DailySchedule.Duration = actualPolicy.SnapshotSchedulePolicy.Schedule.DailySchedule.Duration
	}

	// Hourly
	if actualPolicy.SnapshotSchedulePolicy.Schedule.HourlySchedule != nil {
		expectedPolicy.SnapshotSchedulePolicy.Schedule.HourlySchedule.HoursInCycle = hoursInCycle
		expectedPolicy.SnapshotSchedulePolicy.Schedule.HourlySchedule.StartTime = startTime
		expectedPolicy.SnapshotSchedulePolicy.Schedule.HourlySchedule.Duration = actualPolicy.SnapshotSchedulePolicy.Schedule.HourlySchedule.Duration
	}

	// Weekly
	if actualPolicy.SnapshotSchedulePolicy.Schedule.WeeklySchedule != nil {
		expectedPolicy.SnapshotSchedulePolicy.Schedule.WeeklySchedule.DayOfWeeks = dayOfWeeks
	}

	expectedPolicy.SnapshotSchedulePolicy.SnapshotProperties.GuestFlush = guestFlush
	expectedPolicy.SnapshotSchedulePolicy.SnapshotProperties.Labels = labels
	expectedPolicy.SnapshotSchedulePolicy.SnapshotProperties.StorageLocations = storageLocations

	if reflect.DeepEqual(expectedPolicy, actualPolicy.SnapshotSchedulePolicy) {
		log.Errorf("error! actual policy and expected policy do not match:\n\n%#v\n\n", pretty.Diff(actualPolicy.SnapshotSchedulePolicy, expectedPolicy))

		t.Logf("expected policy:\n\n %#v\n\n", expectedPolicy)
		t.Logf("actual policy:\n\n %#v\n\n", actualPolicy)
		t.FailNow()
	}
}
