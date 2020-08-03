# spinnaker_canary_config Resource

Provides a Spinnaker canary config resource.

## Example Usage

```hcl
# Create a new Spinnaker canary config
resource "spinnaker_canary_config" "canary_config" {
    name         = "Test canary config"
    description  = "Canary config for the demo"
    applications = ["${spinnaker_application.my_app.application}"]
    
    metric {
      name = "CPU metric"
      
      query {
         type               = "stackdriver"
         service_type       = "stackdriver"
         per_series_aligner = "ALIGN_DELTA"
         resource_type      = "k8s_node"
         metric_type        = "kubernetes.io/cpu"
      }
      
      groups     = ["cpu"]
      scope_name = "default"
    }
          
    classifier {
      group_weights {
        cpu = 100
      }
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the canary configuration.
* `description` - (Required) Description for the canary config.
* `applications` - (Required) List of the application which the canary config belongs.
* `metric` - (Required) List of the metric to analyze.
* `classifier` - (Required) Classification configuration.
  
## Attribute Reference 

* `metric` - Group of provider's metrics.
    * `name` - (Required) Human readable name of the metric.
    * `query` - (Required) Query config object for your metric source type.
    * `groups` - (Required) List of the group which this query belongs to.
    * `analysis_configurations` - (Optional) Analysis configuration, describes how to judge a given metric.
* `query` - Block for querying the metric from the service.
    * `type` - (Required) Type of the metric.
    * `service_type` - (Required) Type of the service that is providing the metric.
    * `per_series_aligner` - (Optional) Algorithm to align individual time series.
    * `cross_series_reducer` - (Optional) Algorithm to group multiple time series together.
    * `group_by_fields` - (Optional) Group by resource or metric labels to reduce the number of time series.
    * `resource_type` - (Required) Type of the resource.
    * `metric_type` - (Required) Typo of the metric.
 * `analysis_configurations` - Canary analysis configuration.
    * `canary` - (Required) Configuration for canary.
 * `canary` - Configuration for canary.
    * `direction` - (Optional) Direction on how to judge the failure. Options are `increase`, `decrease` or `either`. Defaults to `increase`.
    * `nan_strategy` - (Optional) How to handle NaN values which can occur if the metric does not return data for a particular time interval. Options are `remove` or `replace`. Default is `remove`.
    * `critical` - (Optional) Fails on this metrics error or not.
    * `must_have_data` - (Optional) Used to fail a metric if data is missing.
    * `effect_zize` - (Optional) Controls how much different the metric needs to be to fail or fail critically.
    * `outliers` - (Optional) Controls how to classify and handle outliers.
 * `effect_size` - Controls how much different the metric needs to be to fail or fail critically.
    * `allowed_increase` - (Optional) The multiplier increase that must be met for the metric to fail. Default to `1`.
    * `allowed_decrease` - (Optional) The multiplier decrease that must be met for the metric to fail. Default to `1`.
    * `critical_increase` - (Optional) The multiplier increase that must be met for the metric to be a critical failure and fail the entire analysis with a score of 0. Default to `1`.
    * `critical_decrease` - (Optional)  The multiplier decrease that must be met for the metric to be a critical failure and fail the entire analysis with a score of 0. Default to `1`.
 * `outliers` -  Controls how to classify and handle outliers.
    * `strategy` - (Optional) Remove or keep outliers. Options are `remove` and `keep`. Default is `keep`.
    * `outlier_factor` - (Option) The degree of significance a data point has to differ from other observations to be considered an outlier. Default is `3.0`.
 * `classifier` - Define how the metrics are classified.
    * `group_wieghts` - (Required) Define the weight for each groups.
 * `group_weights` - Define the weight for each groups by mapping the group and the weight.

## Import

Canary config can be imported using their id, e.g.

```
$ terraform import spinnaker_canary  9753bd1b-3a5c-4104-99ea-26fbc7c78ead
```
