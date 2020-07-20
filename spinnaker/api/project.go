package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

var defaultProjectConfig = map[string]interface{}{
	"applications":    []string{},
	"clusters":        []interface{}{},
	"pipelineConfigs": []interface{}{},
}

// UpseartApplicationTask represents the Spinnaker createApplication Application API object
type UpsertApplicationTask map[string]interface{}

// NewUpsertApplicationTask returns a Spinnaker upsertApplication Application API object
// by passed resource data configured
func NewUpsertApplicationTask(d *schema.ResourceData) (UpsertApplicationTask, error) {
	project := map[string]interface{}{}
	project["name"] = d.Get("name").(string)
	project["email"] = d.Get("email").(string)
	project["config"] = defaultProjectConfig

	if input, ok := d.GetOkExists("config"); ok {
		config := map[string]interface{}{
			"applications":    []interface{}{},
			"clusters":        []interface{}{},
			"pipelineConfigs": []interface{}{},
		}

		for _, v := range input.([]interface{}) {
			configInputs, ok := v.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("config is not map type isn't %T", v)
			}

			for k, configInput := range configInputs {
				switch k {
				case "applications":
					log.Println(configInput)
					applications := configInput.([]interface{})
					if !ok {
						return nil, fmt.Errorf("can't convert applications to list of string, got:%T", configInput)
					}

					config["applications"] = convToStringArray(applications)
				case "cluster":
					clusterInputs, ok := configInput.([]interface{})
					if !ok {
						return nil, fmt.Errorf("can't convert clusters to list of string, got:%T", configInput)
					}

					clusters := convToMapArray(clusterInputs)

					res := make([]map[string]interface{}, len(clusters))
					for i, cluster := range clusters {
						var detail, stack interface{}
						if v := cluster["details"]; v == "" {
							detail = "*"
						}

						if v := cluster["stack"]; v == "" {
							stack = v
						}

						res[i] = map[string]interface{}{
							"account": cluster["account"],
							"detail":  detail,
							"stack":   stack,
						}
					}

					config["cluster"] = res
				case "pipeline_configuration":
				}
			}
		}

		project["config"] = config
	}

	upsertProjectTask := map[string]interface{}{
		"job":         []interface{}{map[string]interface{}{"type": "upsertProject", "project": project}},
		"application": "spinnaker",
		"description": fmt.Sprintf("Create project: %s", project["name"]),
	}

	return upsertProjectTask, nil
}

// GetApplication gets an application from Spinnaker Gate
func GetProject(client *gate.GatewayClient, projectName string, dest interface{}) error {
	project, resp, err := client.ProjectControllerApi.GetUsingGET1(client.Context, projectName)
	if resp != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Project '%s' not found", projectName)
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Encountered an error getting application, status code: %data-resources", resp.StatusCode)
		}
	}

	if err != nil {
		log.Println("project", project)
		log.Println("resp", resp)
		log.Println("err", err)
		return err
	}

	if err := mapstructure.Decode(project, dest); err != nil {
		return err
	}

	return nil
}

// CreateProject creates passed project
func CreateProject(client *gate.GatewayClient, upsertProjectTask UpsertApplicationTask) error {
	ref, _, err := client.TaskControllerApi.TaskUsingPOST1(client.Context, upsertProjectTask)
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
		return fmt.Errorf("Encountered an error saving application, status code: %data-resources", resp.StatusCode)
	}
	if !taskSucceeded(task) {
		return fmt.Errorf("Encountered an error saving application, task output was: %v", task)
	}

	return nil
}

// DeleteProject deletes a project by project name
func DeleteProject(client *gate.GatewayClient, id string, projectName string) error {
	jobSpec := map[string]interface{}{
		"type": "deleteProject",
		"project": map[string]interface{}{
			"id": id,
		},
	}

	deleteAppTask := map[string]interface{}{
		"job":         []interface{}{jobSpec},
		"project":     projectName,
		"description": fmt.Sprintf("Delete project id: %s", id),
	}

	_, resp, err := client.TaskControllerApi.TaskUsingPOST1(client.Context, deleteAppTask)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error deleting application, status code: %data-resources", resp.StatusCode)
	}

	return nil
}

func convToMapArray(inputs []interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, len(inputs))
	for i, input := range inputs {
		res[i] = input.(map[string]interface{})
	}

	return res
}
