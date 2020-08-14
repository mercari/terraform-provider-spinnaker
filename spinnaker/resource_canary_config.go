package spinnaker

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mercari/terraform-provider-spinnaker/spinnaker/api"
)

func resourceSpinnakerCanaryConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "Name of the canary configuration",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateSpinnakerCanaryConfigName,
			},
			"description": {
				Description: "Description for the canary config",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"applications": {
				Description: "List of the application which the canary config belongs",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"metric": {
				Description: "Metric to analyze",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: getCanaryConfigMetricSchema(),
				},
			},
			"classifier": {
				Type:        schema.TypeList,
				Description: "Classification configuration",
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: getCanaryConfigMetricClassifier(),
				},
			},
		},
		Create: resourceSpinnakerCanaryConfigCreate,
		Read:   resourceSpinnakerCanaryConfigRead,
		Update: resourceSpinnakerCanaryConfigUpdate,
		Delete: resourceSpinnakerCanaryConfigDelete,
		Exists: resourceSpinnakerCanaryConfigExists,
		Importer: &schema.ResourceImporter{
			State: resourceSpinnakerCanaryConfigImport,
		},
	}
}

type canaryConfigRead struct {
	id            string      `json:"id"`
	name          string      `json:"name"`
	description   string      `json:"description"`
	configVersion string      `json:"configVersion"`
	applications  []string    `json:"applications"`
	judge         *judge      `json:"judge"`
	metrics       []metric    `json:"metrics"`
	templates     *templates  `json:"templates"`
	classifier    *classifier `json:"classifier"`
}

type judge struct {
	name                string              `json:"name"`
	judgeConfigurations judgeConfigurations `json:"judgeConfigurations"`
}

type judgeConfigurations struct {
}

type metric struct {
	name string
}

type templates struct {
}

type classifier struct {
	groupWeights map[string]int `json:"groupWeights"`
}

func resourceSpinnakerCanaryConfigCreate(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	config, err := api.NewCanaryConfig(d)
	if err != nil {
		return err
	}

	id, err := api.CreateCanaryConfig(client, config)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceSpinnakerCanaryConfigRead(d, meta)
}

func resourceSpinnakerCanaryConfigRead(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	id := d.Id()

	config := &canaryConfigRead{}
	if err := api.GetCanaryConfig(client, id, config); err != nil {
		return err
	}

	if v := config.name; v != "" {
		d.Set("name", v)
	}

	if v := config.description; v != "" {
		d.Set("description", v)
	}

	if v := config.applications; v != nil {
		d.Set("applications", v)
	}

	if v := config.judge; v != nil {
		d.Set("judge", v)
	}

	if v := config.metrics; v != nil {
		tfMetrics := buildTerraformMetrics(v)
		d.Set("metrics", tfMetrics)
	}

	if v := config.templates; v != nil {
		d.Set("template", v)
	}

	if v := config.classifier; v != nil {
		d.Set("classifier", v)
	}

	return nil
}

func resourceSpinnakerCanaryConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	config, err := api.NewCanaryConfig(d)
	if err != nil {
		return err
	}

	id := d.Id()
	if err := api.UpdateCanaryConfig(client, id, config); err != nil {
		return err
	}

	return resourceSpinnakerApplicationRead(d, meta)
}

func resourceSpinnakerCanaryConfigDelete(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	id := d.Id()
	if err := api.DeleteCanaryConfig(client, id); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSpinnakerCanaryConfigExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	id := d.Id()

	var cfg canaryConfigRead
	if err := api.GetCanaryConfig(client, id, &cfg); err != nil {
		errmsg := err.Error()
		if strings.Contains(errmsg, "not found") {
			return false, nil
		}
		return false, err
	}

	if cfg.id == "" {
		return false, nil
	}

	return true, nil
}

func resourceSpinnakerCanaryConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSpinnakerCanaryConfigRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func getCanaryConfigMetricSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Human readable name of the metric",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"query": {
			Type:        schema.TypeList,
			Description: "Query config object for your metric source type",
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: getCanaryConfigMetricQuerySchema(),
			},
		},
		"analysis_configurations": {
			Type:        schema.TypeList,
			Description: "Analysis configuration",
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: getCanaryConfigMetricQueryAnalysisConfiguration(),
			},
		},
		"groups": {
			Type:        schema.TypeList,
			Description: "List of the group which this query belongs to",
			Required:    true,
			MaxItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func getCanaryConfigMetricQuerySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "Type of the metric",
			Required:    true,
		},
		"service_type": {
			Type:        schema.TypeString,
			Description: "Type of the service",
			Required:    true,
		},
		"per_series_aligner": {
			Type:         schema.TypeString,
			Description:  "Algorithm to align individual time series",
			Optional:     true,
			ValidateFunc: validateSpinnakerCanaryConfigQueryAligner,
		},
		"cross_series_reducer": {
			Type:         schema.TypeString,
			Description:  " Algorithm to group multiple time series together",
			Optional:     true,
			ValidateFunc: validateSpinnakerCanaryConfigQueryReducer,
		},
		"group_by_fields": {
			Type:        schema.TypeList,
			Description: "Group by resource or metric labels to reduce the number of time series",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"resource_type": {
			Type:        schema.TypeString,
			Description: "Type of the resource",
			Required:    true,
		},
		"metric_type": {
			Type:        schema.TypeString,
			Description: "Type of the metric",
			Required:    true,
		},
	}
}

func getCanaryConfigMetricQueryAnalysisConfiguration() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"canary": {
			Type:        schema.TypeList,
			Description: "Canary analysis configuration",
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: getCanaryConfigMetricQueryAnalysisConfigurationCanary(),
			},
		},
	}
}

func getCanaryConfigMetricQueryAnalysisConfigurationCanary() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"direction": {
			Type:         schema.TypeString,
			Description:  "Direction on how to judge the failure",
			Optional:     true,
			ValidateFunc: validateSpinnakerCanaryConfigDirection,
			Default:      "increase",
		},
		"nan_strategy": {
			Type:         schema.TypeString,
			Description:  " How to handle NaN values which can occur if the metric does not return data",
			Optional:     true,
			ValidateFunc: validateSpinnakerCanaryConfigNaNStrategy,
			Default:      "remove",
		},
		"critical": {
			Type:        schema.TypeBool,
			Description: "Fails on this metrics error or not",
			Optional:    true,
			Default:     false,
		},
		"must_have_data": {
			Type:        schema.TypeBool,
			Description: "Used to fail a metric if data is missing",
			Optional:    true,
			Default:     false,
		},
		"effect_size": {
			Type:        schema.TypeList,
			Description: "Controls how much different the metric needs to be",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"allowed_increase": {
						Type:        schema.TypeInt,
						Description: "The multiplier increase that must be met for the metric to fail",
						Optional:    true,
						Default:     1,
					},
					"allowed_decrease": {
						Type:        schema.TypeInt,
						Description: "The multiplier decrease that must be met for the metric to fail",
						Optional:    true,
						Default:     1,
					},
					"critical_increase": {
						Type:        schema.TypeInt,
						Description: "The multiplier increase that must be met for the metric to be a critical failure and fail the entire analysis with a score of 0",
						Optional:    true,
						Default:     1,
					},
					"critical_decrease": {
						Type:        schema.TypeInt,
						Description: "The multiplier decrease that must be met for the metric to be a critical failure and fail the entire analysis with a score of 0",
						Optional:    true,
						Default:     1,
					},
				},
			},
		},
		"outliers": {
			Type:        schema.TypeList,
			Description: "Controls how to classify and handle outliers",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"strategy": {
						Type:         schema.TypeString,
						Description:  "Remove or keep outliers",
						Optional:     true,
						Default:      "keep",
						ValidateFunc: validateSpinnakerCanaryConfigOutlierStrategy,
					},
					"outlier_factor": {
						Type:        schema.TypeFloat,
						Description: "The degree of significance a data point has to differ from other observations to be considered an outlier",
						Optional:    true,
						Default:     3.0,
					},
				},
			},
		},
	}
}

func getCanaryConfigMetricClassifier() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_weights": {
			Type:        schema.TypeMap,
			Description: "Weight for each groups",
			Required:    true,
		},
	}
}

func buildTerraformMetrics(metrics []metric) *map[string]interface{} {
	res := map[string]interface{}{}
	for _, metric := range metrics {
		r := map[string]interface{}{}
		res[metric.name] = r
	}

	return &res
}

func validateSpinnakerCanaryConfigName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9-_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("Only alphanumeric characters or '-' and '_' allowed in %q", k))
	}
	return
}

func validateSpinnakerCanaryConfigQueryAligner(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for _, allowedAligner := range api.AllowedAligners {
		if value == allowedAligner {
			return
		}
	}

	return nil, []error{fmt.Errorf("not supported aligner algorythm %s", value)}
}

func validateSpinnakerCanaryConfigQueryReducer(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for _, allowedReducer := range api.AllowedReducers {
		if value == allowedReducer {
			return
		}
	}

	return nil, []error{fmt.Errorf("not supported reducer algorythm %s", value)}
}

func validateSpinnakerCanaryConfigDirection(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for _, allowedDirection := range api.AllowedDirections {
		if value == allowedDirection {
			return
		}
	}

	return nil, []error{fmt.Errorf("not supported direction %s", value)}
}

func validateSpinnakerCanaryConfigNaNStrategy(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for _, allowedNaNStrategy := range api.AllowedNaNStrategies {
		if value == allowedNaNStrategy {
			return
		}
	}

	return nil, []error{fmt.Errorf("not supported NaN strategy %s", value)}
}

func validateSpinnakerCanaryConfigOutlierStrategy(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for _, allowedOutlierStrategy := range api.AllowedOutlierStrategies {
		if value == allowedOutlierStrategy {
			return
		}
	}

	return nil, []error{fmt.Errorf("not supported outlier strategy %s", value)}
}
