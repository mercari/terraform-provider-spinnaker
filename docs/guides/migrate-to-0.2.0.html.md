---
page_title: "Migrate provider to 0.2.0"
subcategory: "Migration guide"
---

# Migrate provider to 0.2.0

This project was orginally forked from [armory-io/terraform-provider-spinnaker](https://github.com/armory-io/terraform-provider-spinnaker) and released as `0.1.0`.
Since the project was stale, we decided to develop entirely independent.

We make many updates and we are going to release `0.2.0` containing breaking changes.
This guide describes the items you have to change before upgrading this Terraform provider.

## Provider schema

Change `server` field to `gate_endpoint`. This aims to clarify what the attribute is for.

```diff
provider "spinnaker" {
- server = "https://spinnaker-.xxx.com"
+ gate_endpoint = "https://spinnaker-.xxx.com"
}
```

## `spinnaker_application` resource

Change `application` field to `name`. This is much more easy to understand what the attribute is for.

```diff
resource "spinnaker_application" "my_app" {
- application = "keke-app"
+ name = "keke-app"
  email = "xxxx@mercari.com"
}
```

## Removed `spinnaker_pipeline_template_config` resource

We currently don't support this resource. 
The reason is that the initial implementation only supported V1 schema which is no more used widely.
We are planning to implement with V2 schema in the near future.

