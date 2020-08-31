package spinnaker

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mercari/terraform-provider-spinnaker/spinnaker/api"
)

const (
	defaultInstancePort = 80
)

func resourceSpinnakerApplication() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Spinnaker application resourc",
		Schema: map[string]*schema.Schema{
			"application": {
				Description:   "Name of the Application",
				Type:          schema.TypeString,
				Deprecated:    "use `name` instead",
				Optional:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validateSpinnakerApplicationName,
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
		CreateContext: resourceSpinnakerApplicationCreate,
		ReadContext:   resourceSpinnakerApplicationRead,
		UpdateContext: resourceSpinnakerApplicationUpdate,
		DeleteContext: resourceSpinnakerApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSpinnakerApplicationImport,
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

func resourceSpinnakerApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clientConfig := meta.(gateConfig)
	var diags diag.Diagnostics
	client := clientConfig.client
	appName := api.GetApplicationName(d)

	task, err := api.NewCreateApplicationTask(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := api.CreateApplication(client, task); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(appName)
	return diags
}

func resourceSpinnakerApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clientConfig := meta.(gateConfig)
	var diags diag.Diagnostics
	client := clientConfig.client
	appName := api.GetApplicationName(d)

	app := &applicationRead{}
	if err := api.GetApplication(client, appName, app); err != nil {
		return diag.FromErr(err)
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
		tfPermissions, err := buildTerraformPermissions(v)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("permissions", tfPermissions)
	}

	return diags
}

func resourceSpinnakerApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	task, err := api.NewCreateApplicationTask(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := api.CreateApplication(client, task); err != nil {
		return diag.FromErr(err)
	}
	return resourceSpinnakerApplicationRead(ctx, d, meta)
}

func resourceSpinnakerApplicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clientConfig := meta.(gateConfig)
	var diags diag.Diagnostics
	client := clientConfig.client
	appName := api.GetApplicationName(d)

	if err := api.DeleteApplication(client, appName); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func resourceSpinnakerApplicationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if diags := resourceSpinnakerApplicationRead(context.Background(), d, meta); diags.HasError() {
		return nil, fmt.Errorf("failed to read spinnaker application")
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

func validateSpinnakerApplicationName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("Only alphanumeric characters or '-' allowed in %q", k))
	}

	return
}
