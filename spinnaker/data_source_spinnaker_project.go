package spinnaker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceProject() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Spinnaker project resource",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getProjectConfigSchema(),
				},
			},
		},
		ReadContext: resourceSpinnakerProjectRead,
	}
}
