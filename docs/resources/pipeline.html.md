# spinnaker_pipeline Resource

Provides a Spinnaker pipeline resource.

## Example Usage

```hcl
# Create a new Spinnaker pipeline
resource "spinnaker_pipeline" "pipeline" {
    application = "${spinnaker_application.my_app.application}"
    name        = "Example Pipeline"
    pipeline    = file("pipelines/example.json")
}
```

## Argument Reference

The following arguments are supported:

* `application` - (Required) The Name of the application.
* `name` - (Required) Pipeline name.
* `pipeline` - (Required) Pipeline JSON content.

## Import

Applications can be imported using their Spinnaker application and pipeline name, e.g.

```
$ terraform import spinnaker_pipeline.pipeline my_app.pipeline
```
