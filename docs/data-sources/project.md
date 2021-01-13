# spinnaker_application Data Source

Use this data source to retrieve information about Spinnaker project.

## Example Usage

```
data "spinnaker_project" "my_proj" {}
```

## Attributes Reference

* `id` - Project ID
* `name` - Name of the project
* `email` - Email of the owner
* `config` - Detail configuration.

## Attribute Reference 

* `config` - Configures these items for the project.
    * `applications` - (Optional) List of the applications which belongs to this project.
    * `cluster` - (Optional) List of [Cluster](https://spinnaker.io/concepts/clusters/#clusters) configuration.
* `cluster` - [Cluster](https://spinnaker.io/concepts/clusters/#clusters) configuration.
    * `account`- (Required): Cluster account.
    * `detail` - (Optional) Detail option of the cluster. Default value is `*`(all).
    * `stack` - (Optional) Stack option of the cluster. Default value is `*`(all).
* `pipeline_config` - [Pipeline configuration](https://spinnaker.io/concepts/pipelines/#pipeline-configuration)
    * `application` - (Required) Application of the pipeline config.
    * `pipeline_config_id` - (Required) ID of the pipeline
