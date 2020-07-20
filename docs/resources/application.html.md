# spinnaker_application Resource

Provides a Spinnaker application resource.

## Example Usage

```hcl
# Create a new Spinnaker application
resource "spinnaker_application" "my_app" {
    application = "my-app"
    email       = "keisuke.yamashita@mercari.com"
}
```

## Argument Reference

The following arguments are supported.

* `application` - (Required) The Name of the application.
* `email` - (Required) Email of the owner.
* `cloud_providers` - (Optional) List of the cloud providers.
* `instance_port` - (Optional) Port of the Spinnaker generated links. Default to `80`.
* `permission` - (Optional) Nested block describing a application permission configuration. You have to enable [Authorization(RBAC)](https://spinnaker.io/setup/security/authorization/) for your Spinnaker to use this feature.
  
## Attribute Reference 

~> **Be careful!** You can accidentally lock yourself out of your Spinnaker application using `permission` attribute. One user or team should obtain `write` permission to edit the application after creation.

* `permission` - this block will have the following structure.
    * `user` - (Required) ID of the user. The ID type depends on the authorization methods. For example, the ID will be the email address if you use G Suite. Also, if you use GitHub Teams the ID will be the team name.   
    * `accesses` - (Required) List of the access permission. The options are `READ`, `EXECUTE` and `WRITE`.
  
## Import

Applications can be imported using their Spinnaker application name, e.g.

```
$ terraform import spinnaker_application.my_app my_app
```
