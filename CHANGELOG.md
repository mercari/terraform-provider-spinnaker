## v0.2.3 (Unreleased)

## v0.2.2 (Auguest 12, 2020)
### Improvements
* Add `spinnaker_canary_config` resource/data source for defining canary config. ([#55](https://github.com/mercari/terraform-provider-spinnaker/pull/55))
    * Note that this currently supports only Cloud Monitoring.
* Introduce [Terraform Plugin SDK](https://www.terraform.io/docs/extend/plugin-sdk.html). ([#56](https://github.com/mercari/terraform-provider-spinnaker/pull/56))

### Bug fixes
* Fixed `gate_endpoint` attribute is not loaded. ([#52](https://github.com/mercari/terraform-provider-spinnaker/pull/52))

### Breaking change
* Add prefix `v` from next release([#59](https://github.com/mercari/terraform-provider-spinnaker/pull/59)

## 0.2.1 (July 27, 2020)
### Bug fixes

* Fixed crash on changes when using deprecated argument `name` for `spinnaker_application`. ([#47](https://github.com/mercari/terraform-provider-spinnaker/pull/47))

## 0.2.0 (July 22, 2020)
### Improvements

* Add `cloud_providers`, `instance_port` to `spinnaker_applicaiton` resource ([#12](https://github.com/mercari/terraform-provider-spinnaker/pull/12))
* Add data resources for these resources ([#12](https://github.com/mercari/terraform-provider-spinnaker/pull/12))
    * `spinnaker_application`
    * `spinnaker_pipeline`
* Add application permission([#22](https://github.com/mercari/terraform-provider-spinnaker/pull/22))
* Migrate legacy provider docs([#24](https://github.com/mercari/terraform-provider-spinnaker/pull/24))
* Add `spinnaker_project` resource/data source([#25](https://github.com/mercari/terraform-provider-spinnaker/pull/25))

### Breaking changes

* Removed `spinnaker_pipeline_template_config` resources. ([#11](https://github.com/mercari/terraform-provider-spinnaker/pull/21))

## 0.1.0 (March 23th, 2020)

* Released the forked [armory-io/terraform-provider-spinnaker](https://github.com/armory-io/terraform-provider-spinnaker) as `0.1.0`
