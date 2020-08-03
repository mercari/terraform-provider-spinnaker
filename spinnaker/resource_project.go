package spinnaker

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mercari/terraform-provider-spinnaker/spinnaker/api"
)

func resourceSpinnakerProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "Name of the project",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateSpinnakerApplicationName,
			},
			"email": {
				Description: "Email of the owner",
				Type:        schema.TypeString,
				Required:    true,
			},
			"config": {
				Description: "Configuration of the project",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: getProjectConfigSchema(),
				},
			},
		},
		Create: resourceSpinnakerProjectCreate,
		Read:   resourceSpinnakerProjectRead,
		Update: resourceSpinnakerProjectUpdate,
		Delete: resourceSpinnakerProjectDelete,
		Exists: resourceSpinnakerProjectExists,
		Importer: &schema.ResourceImporter{
			State: resourceSpinnakerProjectImport,
		},
	}
}

type projectRead struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Config *Config `json:"config"`
}

type Config struct {
	Applications    []string         `json:"applications"`
	Clusters        []Clusters       `json:"clusters"`
	PipelineConfigs []PipelineConfig `json:"pipelineConfigs"`
}

type Clusters struct {
	Account string `json:"account"`
	Detail  string `json:"detail"`
	Stack   string `json:"stack"`
}

type PipelineConfig struct {
	Application string `json:"application"`
	ID          string `json:"pipelineConfigId"`
}

func resourceSpinnakerProjectCreate(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	projectName := d.Get("name").(string)

	task, err := api.NewUpsertApplicationTask(d)
	if err != nil {
		return err
	}

	if err := api.CreateProject(client, task); err != nil {
		return err
	}

	d.SetId(projectName)
	return resourceSpinnakerProjectRead(d, meta)
}

func resourceSpinnakerProjectRead(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	projectName := d.Get("name").(string)
	if projectName == "" {
		projectName = d.Id()
	}

	app := &projectRead{}
	if err := api.GetProject(client, projectName, app); err != nil {
		return err
	}

	if app == nil {
		d.SetId("")
		return nil
	}

	if v := app.ID; v != "" {
		d.Set("id", v)
	}

	if v := app.Name; v != "" {
		d.Set("name", v)
	}
	if v := app.Email; v != "" {
		d.Set("email", v)
	}
	if v := app.Config; v != nil {
		d.Set("config", v)
	}

	return nil
}

func resourceSpinnakerProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	task, err := api.NewCreateApplicationTask(d)
	if err != nil {
		return err
	}

	if err := api.CreateApplication(client, task); err != nil {
		return err
	}
	return resourceSpinnakerProjectRead(d, meta)
}

func resourceSpinnakerProjectDelete(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	id := d.Id()
	appName := d.Get("name").(string)

	if err := api.DeleteProject(client, id, appName); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSpinnakerProjectExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	appName := d.Get("name").(string)

	var app applicationRead
	if err := api.GetApplication(client, appName, &app); err != nil {
		errmsg := err.Error()
		if strings.Contains(errmsg, "not found") {
			return false, nil
		}
		return false, err
	}

	if app.Name == "" {
		return false, nil
	}

	return true, nil
}

func resourceSpinnakerProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSpinnakerProjectRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func getProjectConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"applications": {
			Type:        schema.TypeList,
			Description: "List of belonging Spinnaker Applications",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"cluster": {
			Type:        schema.TypeList,
			Description: "List of clusters",
			Elem: &schema.Resource{
				Schema: getProjectClusterSchema(),
			},
			Optional: true,
		},
		"pipeline_config": {
			Type:        schema.TypeList,
			Description: "List of pipeline configurations",
			Elem: &schema.Resource{
				Schema: getProjectPipelineConfigSchema(),
			},
			Optional: true,
		},
	}
}

func getProjectClusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account": {
			Type:        schema.TypeString,
			Description: "Name of the account",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"detail": {
			Type:        schema.TypeList,
			Description: "Detail of the cluster",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
		},
		"stack": {
			Type:        schema.TypeList,
			Description: "Stack of the cluster",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
		},
	}
}

func getProjectPipelineConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application": {
			Type:        schema.TypeList,
			Description: "Application which to refer the pipeline",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"pipeline_config_id": {
			Type:        schema.TypeList,
			Description: "ID of the pipeline config",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Required:    true,
		},
	}
}
