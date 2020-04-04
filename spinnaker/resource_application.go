package spinnaker

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mercari/terraform-provider-spinnaker/spinnaker/api"
)

const (
	defaultInstancePort = 80
)

func resourceSpinnakerApplication() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "Name of the Application",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateSpinnakerApplicationName,
			},
			"email": {
				Description: "Email of the owner",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cloud_providers": {
				Description: "Cloud providers that is used by the application",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instance_port": {
				Description: "Default port of the Spinnaker generated links",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     defaultInstancePort,
			},
		},
		Create: resourceSpinnakerApplicationCreate,
		Read:   resourceSpinnakerApplicationRead,
		Update: resourceSpinnakerApplicationUpdate,
		Delete: resourceSpinnakerApplicationDelete,
		Exists: resourceSpinnakerApplicationExists,
		Importer: &schema.ResourceImporter{
			State: resourceSpinnakerApplicationImport,
		},
	}
}

type applicationRead struct {
	Name       string                 `json:"name"`
	Attributes *applicationAttributes `json:"attributes"`
}

type applicationAttributes struct {
	Accounts       string `json:"accounts"`
	CloudProviders string `json:"cloudproviders"`
	Email          string `json:"email"`
	InstancePort   int    `json:"instancePort"`
}

func resourceSpinnakerApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	appName := d.Get("name").(string)

	task, err := api.NewCreateApplicationTask(d)
	if err != nil {
		return err
	}

	if err := api.CreateApplication(client, task); err != nil {
		return err
	}

	d.SetId(appName)
	return resourceSpinnakerApplicationRead(d, meta)
}

func resourceSpinnakerApplicationRead(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	appName := d.Get("name").(string)
	app := &applicationRead{}
	if err := api.GetApplication(client, appName, app); err != nil {
		return err
	}

	if app == nil {
		d.SetId("")
		return nil
	}

	if v := app.Attributes.Accounts; v != "" {
		d.Set("accounts", v)
	}
	if v := app.Attributes.CloudProviders; v != "" {
		d.Set("cloud_providers", strings.Split(v, ","))
	}
	if v := app.Attributes.InstancePort; v != 0 {
		d.Set("instance_port", v)
	}

	return nil
}

func resourceSpinnakerApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	task, err := api.NewCreateApplicationTask(d)
	if err != nil {
		return err
	}

	if err := api.CreateApplication(client, task); err != nil {
		return err
	}
	return resourceSpinnakerApplicationRead(d, meta)
}

func resourceSpinnakerApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	appName := d.Get("name").(string)

	if err := api.DeleteApplication(client, appName); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSpinnakerApplicationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
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

func resourceSpinnakerApplicationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSpinnakerApplicationRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func flattenCloudProviders(input string) []string {
	return strings.Split(input, ",")
}
