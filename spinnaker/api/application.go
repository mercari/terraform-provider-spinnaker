package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

// CreateApplicationTask represents the Spinnaker createApplication Application API object
type CreateApplicationTask map[string]interface{}

// NewCreateApplicationTask returns a Spinanker createApplication Application API object
// by passed resource data configued
func NewCreateApplicationTask(d *schema.ResourceData) CreateApplicationTask {
	app := map[string]interface{}{}
	app["name"] = d.Get("name").(string)
	app["email"] = d.Get("email").(string)
	app["instancePort"] = d.Get("instance_port").(int)

	createAppTask := map[string]interface{}{
		"job":         []interface{}{map[string]interface{}{"type": "createApplication", "application": app}},
		"application": app["name"],
		"description": fmt.Sprintf("Create Application: %s", app["name"]),
	}

	return createAppTask
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
