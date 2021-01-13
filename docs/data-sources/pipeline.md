# spinnaker_pipeline Data Source

Use this data source to retrieve information about Spinnaker pipeline.

## Example Usage

```
data "spinnaker_pipeline" "pipeline" {}
```

## Attributes Reference

 * `application` - Name of the application which belongs to
 * `name` - Name of the pipeline
 * `pipeline` - JSON encoded pipeline content
 * `pipeline_id` - ID of the pipeline
 
