# spinnaker_canary_config Data Source

Provides a Spinnaker canary config data source.

## Example Usage

```hcl
# Create a new Spinnaker canary config
data "spinnaker_canary_config" "canary_config" {}
```

## Argument Reference

The following arguments are supported:

* `id` - Canary config ID.
* `name` - Name of the canary configuration.
* `description` - Description for the canary config.
* `applications` - List of the application which the canary config belongs.
* `metric` - List of the metric to analyze.
* `classifier` - Classification configuration.
  
## Attribute Reference 

* `metric` - Group of provider's metrics.
    * `name` - Human readable name of the metric.
    * `query` - Query config object for your metric source type.
    * `groups` - List of the group which this query belongs to.
    * `analysis_configurations` - (Optional) Analysis configuration, describes how to judge a given metric.
* `query` - Block for querying the metric from the service.
    * `type` - Type of the metric.
    * `service_type` - Type of the service that is providing the metric.
    * `per_series_aligner` -  Algorithm to align individual time series.
    * `cross_series_reducer` - Algorithm to group multiple time series together.
    * `group_by_fields` - Group by resource or metric labels to reduce the number of time series.
    * `resource_type` - Type of the resource
    * `metric_type` - Typo of the metric
 * `analysis_configurations` - Canary analysis configuration.
    * `canary` - Configuration for canary.
 * `canary` - Configuration for canary.
    * `direction` - Direction on how to judge the failure. Options are `increase`, `decrease` or `either`. Defaults to `increase`.
    * `nan_strategy` - How to handle NaN values which can occur if the metric does not return data for a particular time interval. Options are `remove` or `replace`. Default is `remove`.
    * `critical` - Fails on this metrics error or not.
    * `must_have_data` - Use to fail a metric if data is missing.
    * `effect_zize` - Controls how much different the metric needs to be to fail or fail critically.
    * `outliers` - Controls how to classify and handle outliers.
 * `effect_size` - Controls how much different the metric needs to be to fail or fail critically.
    * `allowed_increase` - The multiplier increase that must be met for the metric to fail. Default to `1`.
    * `allowed_decrease` - The multiplier decrease that must be met for the metric to fail. Default to `1`.
    * `critical_increase` - The multiplier increase that must be met for the metric to be a critical failure and fail the entire analysis with a score of 0. Default to `1`.
    * `critical_decrease` -  The multiplier decrease that must be met for the metric to be a critical failure and fail the entire analysis with a score of 0. Default to `1`.
 * `outliers` -  Controls how to classify and handle outliers.
    * `strategy` - Remove or keep outliers. Options are `remove` and `keep`. Default is `keep`.
    * `outlier_factor` The degree of significance a data point has to differ from other observations to be considered an outlier. Default is `3.0`.
 * `classifier` - Define how the metrics are classified.
    * `group_wieghts` - Define the weight for each groups.
 * `group_weights` - Define the weight for each groups by mapping the group and the weight.
