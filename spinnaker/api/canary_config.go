package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

var (
	AllowedAligners = []string{
		"ALIGN_NONE",
		"ALIGN_DELTA",
		"ALIGN_RATE",
		"ALIGN_INTERPORATE",
		"ALIGN_NEXT_ORDER",
		"ALIGN_MIN",
		"ALIGN_MAX",
		"ALIGN_MEAN",
		"ALIGN_COUNT",
		"ALIGN_SUM",
		"ALIGN_STDDEV",
		"ALIGN_COUNT_TRUE",
		"ALIGN_COUNT_FALSE",
		"ALIGN_FRACTION_TRUE",
		"ALIGN_PERCENTILE_99",
		"ALIGN_PERCENTILE_95",
		"ALIGN_PERCENTILE_50",
		"ALIGN_PERCENTILE_05",
		"ALIGN_PERCENTILE_CHANGE"}
	AllowedReducers = []string{
		"REDUCE_NONE",
		"REDUCE_MEAN",
		"REDUCE_MAX",
		"REDUCE_MIN",
		"REDUCE_COUNT_TRUE",
		"REDUCE_COUNT_FALSE",
		"REDUCE_FRACTION_TRUE",
		"ALIGN_FRACTION_99",
		"ALIGN_FRACTION_95",
		"ALIGN_FRACTION_50",
		"ALIGN_FRACTION_05",
	}
	AllowedDirections = []string{
		"increase",
		"decrease",
		"either",
	}
	AllowedNaNStrategies = []string{
		"remove",
		"replace",
	}
	AllowedOutlierStrategies = []string{
		"remove",
		"keep",
	}
)

type CanaryConfig map[string]interface{}
type Metrics []map[string]interface{}
type Metric map[string]interface{}
type Classifier = map[string]interface{}
type Query = map[string]interface{}

func NewCanaryConfig(d *schema.ResourceData) (CanaryConfig, error) {
	cfg := map[string]interface{}{}
	cfg["name"] = d.Get("name").(string)
	cfg["description"] = d.Get("description").(string)

	if v, ok := d.GetOkExists("applications"); ok {
		input := v.([]interface{})
		apps := make([]string, len(input))
		for i, input := range v.([]interface{}) {
			app := input.(string)
			apps[i] = app
		}

		cfg["applications"] = apps
	}

	classifiers := convToMapArray(d.Get("classifier").([]interface{}))
	if len(classifiers) > 1 {
		return nil, fmt.Errorf("no more than one classifier block")
	}

	if len(classifiers) == 1 {
		classifier, err := newCanaryConfigClassifier(classifiers[0])
		if err != nil {
			return nil, err
		}

		cfg["classifier"] = classifier
	}

	metrics, err := newCanaryConfigMetrics(convToMapArray(d.Get("metric").([]interface{})))
	if err != nil {
		return nil, err
	}

	cfg["metrics"] = metrics
	//cfg["metrics"] = map[string]interface{}{}

	// Memo(@KeisukeYamashita)
	// "templates" should be always empty
	// See details: CanaryMetricConfigObject(https://github.com/spinnaker/kayenta/blob/master/docs/canary-config.md#canarymetricconfig-object)
	cfg["templates"] = map[string]interface{}{}
	cfg["judge"] = map[string]interface{}{
		"name":                "NetflixACAJudge-v1.0",
		"judgeConfigurations": map[string]interface{}{},
	}
	cfg["configVersion"] = "1"

	return cfg, nil
}

func newCanaryConfigMetrics(ds []map[string]interface{}) (Metrics, error) {
	ms := []map[string]interface{}{}
	for _, d := range ds {
		m, err := newCanaryConfigMetric(d)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func newCanaryConfigMetric(d map[string]interface{}) (Metric, error) {
	m := map[string]interface{}{}
	m["name"] = d["name"].(string)
	query, err := newCanaryConfigQuery(convToMapArray(d["query"].([]interface{})))
	if err != nil {
		return nil, err
	}

	m["query"] = query

	analysisConfiguration := map[string]interface{}{}
	m["analysisConfigurations"] = analysisConfiguration

	groups := convToStringArray(d["groups"].([]interface{}))
	m["groups"] = groups

	// Memo(@KeisukeYamashita)
	// "scopeName" should be always "default
	// See details: CanaryMetricConfigObject(https://github.com/spinnaker/kayenta/blob/master/docs/canary-config.md#canarymetricconfig-object)
	m["scopeName"] = "default"
	return m, nil
}

func newCanaryConfigQuery(ds []map[string]interface{}) (Query, error) {
	if len(ds) != 1 {
		return nil, fmt.Errorf("no more than one query in metric block")
	}

	d := ds[0]
	q := map[string]interface{}{}
	q["type"] = d["type"].(string)
	q["serviceType"] = d["service_type"].(string)

	if v := d["per_series_aligner"].(string); v != "" {
		q["perSeriesAligner"] = v
	}

	q["resourceType"] = d["resource_type"].(string)
	q["metricType"] = d["metric_type"].(string)
	return q, nil
}

func newCanaryConfigAnalysisConfiguration(ds []map[string]interface{}) (map[string]interface{}, error) {
	if len(ds) != 1 {
		return nil, fmt.Errorf("no more than one analysis_configurations block")
	}

	c := map[string]interface{}{}
	d := ds[0]
	canary, err := newCanaryConfigAnalysisConfigurationCanary(d["canary"].([]map[string]interface{}))
	if err != nil {
		return nil, err
	}

	c["canary"] = canary
	return c, nil
}

func newCanaryConfigAnalysisConfigurationCanary(ds []map[string]interface{}) (map[string]interface{}, error) {
	if len(ds) != 1 {
		return nil, fmt.Errorf("no more than one canary block")
	}

	c := map[string]interface{}{}
	d := ds[0]
	if v, ok := d["direction"].(string); ok {
		c["direction"] = v
	}

	if v, ok := d["nan_strategy"].(string); ok {
		c["nan_strategy"] = v
	}

	if v, ok := d["critical"].(bool); ok {
		c["critical"] = v
	}

	if v, ok := d["must_have_data"].(bool); ok {
		c["must_have_data"] = v
	}

	if vs, ok := d["effect_size"].([]map[string]interface{}); ok {
		if len(vs) != 1 {
			return nil, fmt.Errorf("no more than one effect_size block")
		}

		v := vs[0]
		c["effect_size"] = newCanaryConfigEffectSize(v)
	}

	if vs, ok := d["outliers"].([]map[string]interface{}); ok {
		if len(vs) != 1 {
			return nil, fmt.Errorf("no more than one outliers block")
		}

		v := vs[0]
		c["outliers"] = newCanaryConfigOutliers(v)
	}

	return c, nil
}

func newCanaryConfigEffectSize(d map[string]interface{}) map[string]interface{} {
	es := map[string]interface{}{}

	if v, ok := d["allowed_increase"]; ok {
		es["allowed_increase"] = v
	}

	if v, ok := d["allowed_decrease"]; ok {
		es["allowed_decrease"] = v
	}

	if v, ok := d["critical_increase"]; ok {
		es["critical_increase"] = v
	}

	if v, ok := d["critical_decrease"]; ok {
		es["critical_decrease"] = v
	}

	return es
}

func newCanaryConfigOutliers(d map[string]interface{}) map[string]interface{} {
	o := map[string]interface{}{}

	if v, ok := d["strategy"]; ok {
		o["strategy"] = v
	}

	if v, ok := d["outlier_factor"]; ok {
		o["outlier_factor"] = v
	}

	return o
}

func newCanaryConfigClassifier(d map[string]interface{}) (Classifier, error) {
	c := map[string]interface{}{}

	input := d["group_weights"].(map[string]interface{})
	groupWeights := map[string]int{}
	for k, v := range input {
		weight, err := strconv.Atoi(v.(string))
		if err != nil {
			return nil, err
		}

		groupWeights[k] = weight
	}

	c["groupWeights"] = groupWeights
	return c, nil
}

// CreateCanaryConfig creates passed canary config
func CreateCanaryConfig(client *gate.GatewayClient, config CanaryConfig) (string, error) {
	ref, resp, err := client.V2CanaryConfigControllerApi.CreateCanaryConfigUsingPOST(context.Background(), config, map[string]interface{}{})
	if err != nil {
		log.Println(fmt.Sprintf("%#v", config))
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("Encountered an error saving canary config, status code: %d", resp.StatusCode)
	}

	canaryConfig := ref.(map[string]interface{})
	return canaryConfig["canaryConfigId"].(string), nil
}

func GetCanaryConfig(client *gate.GatewayClient, id string, dest interface{}) error {
	conf, resp, err := client.V2CanaryConfigControllerApi.GetCanaryConfigUsingGET(context.Background(), id, map[string]interface{}{})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error getting canary config with id %s, status code: %d\n", id, resp.StatusCode)
	}

	if err := mapstructure.Decode(conf, dest); err != nil {
		return err
	}

	return nil
}

func DeleteCanaryConfig(client *gate.GatewayClient, id string) error {
	resp, err := client.V2CanaryConfigControllerApi.DeleteCanaryConfigUsingDELETE(context.Background(), id, map[string]interface{}{})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error getting canary config with id %s, status code: %d\n", id, resp.StatusCode)
	}

	return nil
}

func UpdateCanaryConfig(client *gate.GatewayClient, id string, config CanaryConfig) error {
	_, resp, err := client.V2CanaryConfigControllerApi.UpdateCanaryConfigUsingPUT(context.Background(), id, config, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Encountered an error saving canary config, status code: %d", resp.StatusCode)
	}

	return nil
}
