package spinnaker

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func datasourceApplication() *schema.Resource {
	return &schema.Resource{
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
				Type:     schema.TypeString,
				Computed: true,
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
		Read: resourceSpinnakerProjectRead,
	}
}
