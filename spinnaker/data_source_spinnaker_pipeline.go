package spinnaker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourcePipeline() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Pipeline deprecated because is not HCL native. We plan to support in the next major release",
		Schema: map[string]*schema.Schema{
			"application": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSpinnakerApplicationName,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"pipeline": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipeline_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Read: resourcePipelineRead,
	}
}
