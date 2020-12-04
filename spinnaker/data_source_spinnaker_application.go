package spinnaker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceApplication() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Spinnaker application resource",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSpinnakerApplicationName,
			},
			"email": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"accounts": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instance_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
		ReadContext: resourceSpinnakerApplicationRead,
	}
}
