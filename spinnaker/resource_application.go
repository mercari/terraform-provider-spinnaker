package spinnaker

import (
	"fmt"
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
			"application": {
				Description:   "Name of the Application",
				Type:          schema.TypeString,
				Deprecated:    "use `name` instead",
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Description:  "Name of the Application",
				Type:         schema.TypeString,
				Optional:     true,
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
			"permission": {
				Description: "Application level permissions",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: getApplicationPermissionSchema(),
				},
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
	Accounts       string       `json:"accounts"`
	CloudProviders string       `json:"cloudproviders"`
	Email          string       `json:"email"`
	InstancePort   int          `json:"instancePort"`
	Permissions    *Permissions `json:"permissions"`
}

type Permissions struct {
	Read    []string `json:"READ"`
	Execute []string `json:"EXECUTE"`
	Write   []string `json:"WRITE"`
}

func resourceSpinnakerApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	appName := api.GetApplicationName(d)

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
	var clientConfig gateConfig
	var ok bool
	if clientConfig, ok = meta.(gateConfig); !ok {
		// TODO(@KeisukeYamashita): Remove this. This is mercari specific problem
		// once we fixed the CI, we will remove this
		// This is only nil when the gate client failed to initialize
		// Moreover, this patch will not be released to the main stream
		return fmt.Errorf("gate client failed to instant")
	}

	client := clientConfig.client
	appName := api.GetApplicationName(d)

	app := &applicationRead{}
	if err := api.GetApplication(client, appName, app); err != nil {
		return err
	}

	if app == nil {
		d.SetId("")
		return nil
	}

	if v := app.Name; v != "" {
		if _, deprecated := d.GetOk("application"); deprecated {
			d.Set("name", nil)
			d.Set("application", v)
		} else {
			d.Set("name", v)
			d.Set("application", nil)
		}
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
	if v := app.Attributes.Permissions; v != nil {
		terraformPermissions, err := buildTerraformPermissions(v)
		if err != nil {
			return err
		}

		d.Set("permissions", terraformPermissions)
	}

	return nil
}

func resourceSpinnakerApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	var clientConfig gateConfig
	var ok bool
	if clientConfig, ok = meta.(gateConfig); !ok {
		// TODO(@KeisukeYamashita): Remove this. This is mercari specific problem
		// once we fixed the CI, we will remove this
		// This is only nil when the gate client failed to initialize
		// Moreover, this patch will not be released to the main stream
		return fmt.Errorf("gate client failed to instant")
	}
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
	var clientConfig gateConfig
	var ok bool
	if clientConfig, ok = meta.(gateConfig); !ok {
		// TODO(@KeisukeYamashita): Remove this. This is mercari specific problem
		// once we fixed the CI, we will remove this
		// This is only nil when the gate client failed to initialize
		// Moreover, this patch will not be released to the main stream
		return fmt.Errorf("gate client failed to instant")
	}

	client := clientConfig.client
	appName := api.GetApplicationName(d)

	if err := api.DeleteApplication(client, appName); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSpinnakerApplicationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	var clientConfig gateConfig
	var ok bool
	if clientConfig, ok = meta.(gateConfig); !ok {
		// TODO(@KeisukeYamashita): Remove this. This is mercari specific problem
		// once we fixed the CI, we will remove this
		// This is only nil when the gate client failed to initialize
		// Moreover, this patch will not be released to the main stream
		return ok, fmt.Errorf("gate client failed to instant")
	}

	client := clientConfig.client
	appName := api.GetApplicationName(d)

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

func getApplicationPermissionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user": {
			Type:        schema.TypeString,
			Description: "User ID",
			Required:    true,
		},
		"accesses": {
			Type:        schema.TypeList,
			Description: "List of access",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Required:    true,
		},
	}
}

func buildTerraformPermissions(permissions *Permissions) (*map[string][]string, error) {
	users := map[string][]string{}
	for _, rUser := range permissions.Read {
		users[rUser] = append(users[rUser], "READ")
	}

	for _, xUser := range permissions.Execute {
		users[xUser] = append(users[xUser], "EXECUTE")
	}

	for _, wUser := range permissions.Read {
		users[wUser] = append(users[wUser], "WRITE")
	}

	for user, accesses := range users {
		if len(accesses) > 3 {
			return nil, fmt.Errorf("more than 3 access granted for %s", user)
		}
	}

	return &users, nil
}
