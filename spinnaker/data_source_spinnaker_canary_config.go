package spinnaker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceCanaryConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Canary config id",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name of the canary config",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the canary config",
				Type:        schema.TypeString,
				Required:    true,
			},
			"applications": {
				Description: "Name of the Application",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"metric": {
				Description: "Detail of the metrics",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: getCanaryConfigMetricSchema(),
				},
			},
			"classifier": {
				Type:        schema.TypeList,
				Description: "Classifier of the metrics",
				Required:    true,
				Elem: &schema.Resource{
					Schema: getCanaryConfigMetricClassifier(),
				},
			},
		},
		Read: resourceSpinnakerCanaryConfigRead,
	}
}
