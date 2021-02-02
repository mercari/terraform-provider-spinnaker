# spinnaker_project Resource

Provides a Spinnaker project resource.

## Example Usage

```hcl
# Create a new Spinnaker application
resource "spinnaker_project" "my_project" {
    name   = "my-project"
    email  = "keisuke.yamashita@mercari.com"
}
```

## Argument Reference

* `name` - (Required) Name of the project.
* `email` - (Required) Email of the owner.
* `config` - (Optional) Detail configuration.

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
  
## Import

Applications can be imported using their Spinnaker application name, e.g.

```
$ terraform import spinnaker_project.my_proj my-project
```
