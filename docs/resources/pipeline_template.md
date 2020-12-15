# spinnaker_pipeline_template Resource

Provides a Spinnaker pipeline resource.

## Example Usage

```hcl
# Create a new Spinnaker managed pipeline template
resource "spinnaker_pipeline_template" "pipeline_template" {
    template = file("pipelines/example.json")
}
```

## Argument Reference

The following arguments are supported:

* `template` - (Required) Pipeline JSON content.

## Import

Applications can be imported using their Spinnaker managed pipeline template ID, e.g.

```
$ terraform import spinnaker_pipeline_template.pipeline_template my-template
```
