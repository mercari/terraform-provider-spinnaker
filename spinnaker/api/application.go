package api

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

var (
	// CloudProviders ...
	// See details in Spinnaker Orca
	// ref: https://github.com/spinnaker/orca/blob/master/orca-applications/src/main/groovy/com/netflix/spinnaker/orca/applications/utils/ApplicationNameValidator.groovy
	CloudProviders = map[string]applicationNameConstraint{
		"appengine":    applicationNameConstraint{58, `^[a-z0-9]*$`},
		"aws":          applicationNameConstraint{250, `^[a-zA-Z_0-9.]*$`},
		"dcos":         applicationNameConstraint{127, `^[a-z0-9]*$`},
		"kubernetes":   applicationNameConstraint{63, `^([a-zA-Z][a-zA-Z0-9-]*)$`},
		"gce":          applicationNameConstraint{63, `^([a-zA-Z][a-zA-Z0-9]*)?$`},
		"openstack":    applicationNameConstraint{250, `^[a-zA-Z_0-9.]*$`},
		"tencentcloud": applicationNameConstraint{50, `^[a-zA-Z_0-9.\u4e00-\u9fa5-]*$`},
		"titus":        applicationNameConstraint{250, `^[a-zA-Z_0-9.]*$`},
	}
)

// applicationNameConstraint ...
type applicationNameConstraint struct {
	maxLength int
	regex     string
}

// CreateApplicationTask represents the Spinnaker createApplication Application API object
type CreateApplicationTask map[string]interface{}

// NewCreateApplicationTask returns a Spinanker createApplication Application API object
// by passed resource data configued
func NewCreateApplicationTask(d *schema.ResourceData) (CreateApplicationTask, error) {
	app := map[string]interface{}{}
	app["name"] = d.Get("name").(string)
	app["email"] = d.Get("email").(string)
	app["instancePort"] = d.Get("instance_port").(int)

	if v, ok := d.GetOkExists("cloud_providers"); ok {
		input := v.([]interface{})
		cloudProviders := make([]string, len(input))
		for k, input := range v.([]interface{}) {
			cloudProvider := input.(string)
			if err := validateSpinnakerApplicationNameByCloudProvider(d.Get("name").(string), cloudProvider); err != nil {
				return nil, err
			}

			cloudProviders[k] = cloudProvider
		}

		app["cloudProviders"] = strings.Join(cloudProviders, ",")
	}

	createAppTask := map[string]interface{}{
		"job":         []interface{}{map[string]interface{}{"type": "createApplication", "application": app}},
		"application": app["name"],
		"description": fmt.Sprintf("Create Application: %s", app["name"]),
	}

	return createAppTask, nil
}

// GetApplication gets an application from Spinnaker Gate
func GetApplication(client *gate.GatewayClient, appName string, dest interface{}) error {
	app, resp, err := client.ApplicationControllerApi.GetApplicationUsingGET(client.Context, appName, map[string]interface{}{})
	if resp != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Application '%s' not found", appName)
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Encountered an error getting application, status code: %d", resp.StatusCode)
		}
	}

	if err != nil {
		log.Println("app", app)
		log.Println("resp", resp)
		log.Println("err", err)
		return err
	}

	if err := mapstructure.Decode(app, dest); err != nil {
		return err
	}

	return nil
}

// CreateApplication creates passed application
func CreateApplication(client *gate.GatewayClient, createAppTask CreateApplicationTask) error {
	ref, _, err := client.TaskControllerApi.TaskUsingPOST1(client.Context, createAppTask)
	if err != nil {
		return err
	}

	toks := strings.Split(ref["ref"].(string), "/")
	id := toks[len(toks)-1]

	task, resp, err := client.TaskControllerApi.GetTaskUsingGET1(client.Context, id)
	attempts := 0
	for (task == nil || !taskCompleted(task)) && attempts < 5 {
		toks := strings.Split(ref["ref"].(string), "/")
		id := toks[len(toks)-1]

		task, resp, err = client.TaskControllerApi.GetTaskUsingGET1(client.Context, id)
		attempts++
		time.Sleep(time.Duration(attempts*attempts) * time.Second)
	}

	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Encountered an error saving application, status code: %d", resp.StatusCode)
	}
	if !taskSucceeded(task) {
		return fmt.Errorf("Encountered an error saving application, task output was: %v", task)
	}

	return nil
}

// DeleteApplication deletes an application by application name
func DeleteApplication(client *gate.GatewayClient, appName string) error {
	jobSpec := map[string]interface{}{
		"type": "deleteApplication",
		"application": map[string]interface{}{
			"name": appName,
		},
	}

	deleteAppTask := map[string]interface{}{
		"job":         []interface{}{jobSpec},
		"application": appName,
		"description": fmt.Sprintf("Delete Application: %s", appName),
	}

	_, resp, err := client.TaskControllerApi.TaskUsingPOST1(client.Context, deleteAppTask)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error deleting application, status code: %d", resp.StatusCode)
	}

	return nil
}

func taskCompleted(task map[string]interface{}) bool {
	taskStatus, exists := task["status"]
	if !exists {
		return false
	}

	COMPLETED := [...]string{"SUCCEEDED", "STOPPED", "SKIPPED", "TERMINAL", "FAILED_CONTINUE"}
	for _, status := range COMPLETED {
		if taskStatus == status {
			return true
		}
	}
	return false
}

func taskSucceeded(task map[string]interface{}) bool {
	taskStatus, exists := task["status"]
	if !exists {
		return false
	}

	SUCCESSFUL := [...]string{"SUCCEEDED", "STOPPED", "SKIPPED"}
	for _, status := range SUCCESSFUL {
		if taskStatus == status {
			return true
		}
	}
	return false
}

func validateSpinnakerApplicationNameByCloudProvider(appName, provider string) error {
	if regex, ok := CloudProviders[provider]; ok {
		if !regexp.MustCompile(regex.regex).MatchString(appName) {
			return fmt.Errorf("application name %s for cloud provider %s doesn't match %s", appName, provider, regex.regex)
		}

		if c := len(appName); c > regex.maxLength {
			return fmt.Errorf("application name %s for cloud provider %s is more than its limit %d", appName, provider, c)
		}

		return nil
	}

	return fmt.Errorf("cloud provider %s is not supported", provider)
}
