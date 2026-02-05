<!-- See https://developer.hashicorp.com/terraform/plugin/best-practices/versioning#changelog-specification -->

## UNRELEASED

NOTES:
* Tested against [Sonatype Nexus Repository Manager 3.88.0](https://help.sonatype.com/en/sonatype-nexus-repository-3-88-0-release-notes.html)
  * There is an API regression in NXRM 3.88 that prevents use of the `sonatyperepo_system_config_ldap_connection` resource

## 0.18.1 January 29, 2026

BUG FIXES:
- Resolved inability to update resource `sonatyperepo_system_config_ldap_connection` [GH-271]
- Unable to create resource `sonatyperepo_blob_store_s3` [GH-273]

## 0.18.0 January 29, 2026

ENHANCEMENTS:
- Support for S3 Pre-signed URLs on `sonatyperepo_blob_store_s3` Data Source and Resource [GH-248]
- Resource `sonatyperepo_role` now supports managing roles consisting of only privileges or only roles [GH-259]

BUG FIXES:
- The following resources could not be updated (missing `id` in API call) [GH-261]:
  - `sonatyperepo_capability_base_url`
  - `sonatyperepo_capability_custom_s3_regions`
  - `sonatyperepo_capability_default_role`
  - `sonatyperepo_capability_firewall_audit_and_quarantine`
  - `sonatyperepo_capability_healthcheck`
  - `sonatyperepo_capability_outreach_management`
  - `sonatyperepo_capability_rut_auth`
  - `sonatyperepo_capability_storage_settings`
  - `sonatyperepo_capability_ui_branding`
  - `sonatyperepo_capability_ui_settings`
  - `sonatyperepo_capability_webhook_global`
  - `sonatyperepo_capability_webhook_repository`

NOTES:
* Tested against [Sonatype Nexus Repository Manager 3.87.1](https://help.sonatype.com/en/sonatype-nexus-repository-3-87-0-release-notes.html)

## 0.17.0 January 06, 2026

ENHANCEMENTS:
- Updated the minimum Terraform version to 1.7.0 [GH-253]
  
  _Earlier versions may well work fine, but are now not supported or tested._

BUG FIXES:
- Unable to update `sonatyperepo_capability_audit` resource after import [GH-250]
- Unable to update capabilities (`sonatyperepo_capability_*`) due to missing ID in request [GH-251]

## 0.16.0 December 17, 2025

ENHANCEMENTS:
* **New Resource:** `sonatyperepo_blob_store_group` [GH-136]
* **New Resource:** `sonatyperepo_system_config_http` [GH-246]
* **New Resource:** `sonatyperepo_task_license_expiration_notification` - see [help docs](https://help.sonatype.com/en/license-management.html#license-expiration-notifications) [GH-223]

BUG FIXES:
* Resource `sonatyperepo_capability_audit` could not be imported [GH-237]
* Resource `sonatyperepo_system_iq_connection` always showed plan changes due to `properties` and `last_updated` fields [GH-236]
* Proxy Repository Resources could produce inconsistent results after apply [GH-232] (thanks @yfougeray-euphoria)

NOTES:
* Tested against [Sonatype Nexus Repository Manager 3.86.2](https://help.sonatype.com/en/sonatype-nexus-repository-3-86-0-release-notes.html)

## 0.15.0 December 15, 2025

ENHANCEMENTS:
* Refactored common code and patterns into a shared library to improve mantainability and consistency [GH-208], [GH-201], [GH-199]

  You will see slight changes (improvements, minor corrections and consistency) in descriptions of fields for data sources and resources as a result of this, but there **should be NO breaking changes** - all schemas have been verified as matching version the previous release.

BUG FIXES:
* Could not create `sonatyperepo_repository_conan_proxy` resource [GH-224]

## 0.14.1 December 02, 2025

BUG FIXES:
* Resources `sonatyperepo_capability_*` crashed provider when `notes` was not provided [GH-222]

## 0.14.0 November 25, 2025

ENHANCEMENTS:
* HuggingFace now a supported format for `sonatyperepo_cleanup_policy` [GH-218] (thanks @yfougeray-euphoria)

BUG FIXES:
* Fixed "Provider produced inconsistent result after apply" error when `routing_rule` was provided for Maven, PyPi, R or Ruby Gems repositories [GH-216] (thanks @yfougeray-euphoria)

## 0.13.0 November 24, 2025

ENHANCEMENTS:
* `sonatyperepo_content_selector` resources now support import [GH-209]

BUG FIXES:
* Temporary workaround for `sonatyperepo_repository_docker_hosted` as API does not return `storage.latest_policy` [GH-210]


## 0.12.1 November 19, 2025

BUG FIXES:
* Fixed regression that immpacted all `sonatyperepo_repository_*_proxy` resources where replication was not enabled [GH-206]

## 0.12.0 November 14, 2025

ENHANCEMENTS:
* `sonatyperepo_repository_apt_*` resources now support import [GH-146] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_conda_*` resources now support import [GH-147] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_go_*` resources now support import [GH-148] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_helm_*` resources now support import [GH-145] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_huggingface_proxy` resource now support import [GH-149] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_maven_*` resources now support import [GH-138] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_npm_*` resources now support import [GH-139] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_nuget_*` resources now support import [GH-140] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_p2_proxy` resource now support import [GH-150] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_pypi_*` resources now support import [GH-141] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_raw_*` resources now support import [GH-142] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_rubygems_*` resources now support import [GH-143] (thanks @yfougeray-euphoria)
* `sonatyperepo_repository_yum_*` resources now support import [GH-144] (thanks @yfougeray-euphoria)

## 0.11.1 November 11, 2025

BUG FIXES:
* Fixed issue with import for `sonatyperepo_repository_*_proxy` resources [GH-153] (thanks @yfougeray-euphoria)

## 0.11.0 November 11, 2025

FEATURES:
* **New Data Source:** `sonatyperepo_capabilities` [GH-157]
* **New Data Source:** `sonatyperepo_routing_rule` [GH-152] (thanks @yfougeray-euphoria)
* **New Data Source:** `sonatyperepo_routing_rules` [GH-152] (thanks @yfougeray-euphoria)
* **New Data Source:** `sonatyperepo_user_tokens` [GH-151] (thanks @yfougeray-euphoria)
* **New Resource:** `sonatyperepo_capability_audit` [GH-159]
* **New Resource:** `sonatyperepo_capability_base_url` [GH-156]
* **New Resource:** `sonatyperepo_capability_custom_s3_regions` [GH-186]
* **New Resource:** `sonatyperepo_capability_default_role` [GH-188]
* **New Resource:** `sonatyperepo_capability_firewall_audit_and_quarantine` [GH-163]
* **New Resource:** `sonatyperepo_capability_healthcheck` [GH-163]
* **New Resource:** `sonatyperepo_capability_outreach_management` [GH-166]
* **New Resource:** `sonatyperepo_capability_rut_auth` [GH-195]
* **New Resource:** `sonatyperepo_capability_storage_settings` [GH-194]
* **New Resource:** `sonatyperepo_capability_ui_branding` [GH-168]
* **New Resource:** `sonatyperepo_capability_ui_settings` [GH-169]
* **New Resource:** `sonatyperepo_capability_webhook_global` [GH-183]
* **New Resource:** `sonatyperepo_capability_webhook_repository` [GH-184]
* **New Resource:** `sonatyperepo_routing_rule` [GH-152] (thanks @yfougeray-euphoria)
* **New Resource:** `sonatyperepo_user_tokens` [GH-151] (thanks @yfougeray-euphoria)

## 0.10.0 November 4, 2025

FEATURES:
* **New Resource:** `sonatyperepo_task_blobstore_compact` [GH-92]
* **New Resource:** `sonatyperepo_task_malware_remediator` [GH-92]
* **New Resource:** `sonatyperepo_task_repair_create_browse_nodes` [GH-92]
* **New Resource:** `sonatyperepo_task_repository_maven_remove_snapshots` [GH-92]
* **New Resource:** `sonatyperepo_task_repository_docker_gc` [GH-92]
* **New Resource:** `sonatyperepo_task_repository_docker_upload_purge` [GH-92]

ENHANCEMENTS:
* You can now supply a `version_hint` in the Provider configuration to work around scenarios where the `Server` header is stripped from responses. This can be typical when using a reverse proxy or load balancer in front of Sonatype Nexus Repository. [GH-109]

## 0.9.0 October 30, 2025

ENHANCEMENTS:
* Resource `sonatyperepo_repository_docker_group` now supports import [GH-126] (thanks @yfougeray-euphoria)
* Resource `sonatyperepo_repository_docker_hosted` now supports import [GH-126] (thanks @yfougeray-euphoria)
* Resource `sonatyperepo_repository_docker_proxy` now supports import [GH-126] (thanks @yfougeray-euphoria)
* Resource `sonatyperepo_role` now supports import [GH-128] (thanks @yfougeray-euphoria)
* Resource `sonatyperepo_system_config_mail` now supports import [GH-119] (thanks @yfougeray-euphoria)
* Resource `sonatyperepo_security_realms` now supports import [GH-118] (thanks @yfougeray-euphoria)
* Resource `sonatyperepo_security_saml` now supports import [GH-120] (thanks @yfougeray-euphoria)
* Resource `sonatyperepo_system_anonymous_access` now supports import [GH-117] (thanks @yfougeray-euphoria)
* Resource `sonatyperepo_user` now supports import [GH-97]

## 0.8.2 October 23, 2025

ENHANCEMENTS:
* A number of resources have had their schemas updated to use `Set` rather than `List` to avoid unncessary plan changes due to ordering of data. The impacted data sources and resources are:
  * Data Sources:
    * `sonatyperepo_blob_store_group`
    * `sonatyperepo_roles`
    * `sonatyperepo_users`
  * Resources:
    * `sonatyperepo_privilege_*` (i.e. all Privilege resources)
    * `sonatyperepo_repository_*_hosted` (i.e. all Hosted Repository resources)
    * `sonatyperepo_repository_*_proxy` (i.e. all Proxy Repository resources)
    * `sonatyperepo_role`
    * `sonatyperepo_user`

BUG FIXES:
* **Bug:** Read for resource `sonatyperepo_user` might return incorrect User breaking state [GH-111]
* **Bug:** Subsequent plans for `sonatyperepo_role` showed change due to ordering of privileges [GH-112]

## 0.8.1 October 21, 2025

BUG FIXES:
* **Bug:** Further improve connectivity errors to Sonatype Nexus Repository without the provider crashing [GH-105]

## 0.8.0 October 18, 2025

BUG FIXES:
* **Bug:** `userMemberOfAttribute` was not mapped in API requests for `sonatyperepo_system_config_ldap_connection` resource [GH-101]
* **Bug:** `prefix` for `sonatyperepo_blob_store_s3` was not stored in Terraform State correctly leading to *inconsistent results after apply* error [GH-103]

NOTES:
* Tested against Sonatype Nexus Repository Manager 3.85.0

## 0.7.0 October 17, 2025

FEATURES:

* **New Data Source:** `sonatyperepo_content_selector` [GH-84]
* **New Data Source:** `sonatyperepo_content_selectors` [GH-84]
* **New Resource:** `sonatyperepo_content_selector` [GH-84]

BUG FIXES:

* **Bug:** Changing the `id` on a `sonatyperepo_role` did not force a recreation of the resource - failed with `409 - Conflict` [GH-88]
* **Bug:** Handle connectivity errors to Sonatype Nexus Repository without the provider crashing [GH-94]

## 0.6.2 October 14, 2025

NOTES:

* **Docs:** Clarified that `sonatyperepo_role` resource can be used to create Roles that auto-map to LDAP or SAML groups

## 0.6.1 October 14, 2025

BUG FIXES:

* **Bug:** Unable to create Docker Registry on NXRM 3.84+ as `latest_policy` not returned in API create response

## 0.6.0 September 22, 2025

FEATURES:

* **New Resource:** `sonatyperepo_blob_store_google_cloud` [GH-64]

## 0.5.0 September 22, 2025

FEATURES:

* **New Resource:** `sonatyperepo_cleanup_policy` [GH-58]
* **New Resource:** `sonatyperepo_security_saml` [GH-63]

ENHANCEMENTS:
* resource/sonatyperepo_repository_docker_group: Add `path_enabled` attribute [GH-75]
* resource/sonatyperepo_repository_docker_hosted: Add `path_enabled` attribute [GH-75]
* resource/sonatyperepo_repository_docker_proxy: Add `path_enabled` attribute [GH-75]

NOTES:
* Tested against Sonatype Nexus Repository Manager 3.79.1 through 3.84.1

## 0.4.0 September 3, 2025

FEATURES:

* **New Resource:** `sonatyperepo_security_realms` [GH-60]

## 0.3.0 July 29, 2025

ENHANCEMENTS:

* Confirm Sonatype Nexus Repository Manager is WRITABLE [GH-54]
* Determine version of Sonatype Nexus Repository Manager for future use in this Provider [GH-54]

NOTES:
* Tested against [Sonatype Nexus Repository Manager 3.82.0](https://help.sonatype.com/en/sonatype-nexus-repository-3-82-0-release-notes.html)

## 0.2.0 July 17, 2025

FEATURES:

* **New Resource:** `sonatyperepo_repository_apt_hosted` [GH-35]
* **New Resource:** `sonatyperepo_repository_apt_proxy` [GH-35]
* **New Resource:** `sonatyperepo_repository_cargo_group` [GH-50]
* **New Resource:** `sonatyperepo_repository_cargo_hosted` [GH-50]
* **New Resource:** `sonatyperepo_repository_cargo_proxy` [GH-50]
* **New Resource:** `sonatyperepo_repository_conan_group` [GH-51]
* **New Resource:** `sonatyperepo_repository_conan_hosted` [GH-51]
* **New Resource:** `sonatyperepo_repository_conan_proxy` [GH-51]
* **New Resource:** `sonatyperepo_repository_cocoapods_proxy` [GH-39]
* **New Resource:** `sonatyperepo_repository_composer_proxy` [GH-40]
* **New Resource:** `sonatyperepo_repository_conda_proxy` [GH-45]
* **New Resource:** `sonatyperepo_repository_docker_group` [GH-36]
* **New Resource:** `sonatyperepo_repository_docker_hosted` [GH-36]
* **New Resource:** `sonatyperepo_repository_docker_proxy` [GH-36]
* **New Resource:** `sonatyperepo_repository_gitlfs_hosted` [GH-46]
* **New Resource:** `sonatyperepo_repository_go_group` [GH-41]
* **New Resource:** `sonatyperepo_repository_go_hosted` [GH-41]
* **New Resource:** `sonatyperepo_repository_helm_hosted` [GH-42]
* **New Resource:** `sonatyperepo_repository_helm_proxy` [GH-42]
* **New Resource:** `sonatyperepo_repository_huggingface_proxy` [GH-43]
* **New Resource:** `sonatyperepo_repository_nuget_group` [GH-38]
* **New Resource:** `sonatyperepo_repository_nuget_hosted` [GH-38]
* **New Resource:** `sonatyperepo_repository_nuget_proxy` [GH-38]
* **New Resource:** `sonatyperepo_repository_p2_proxy` [GH-47]
* **New Resource:** `sonatyperepo_repository_pypi_group` [GH-37]
* **New Resource:** `sonatyperepo_repository_pypi_hosted` [GH-37]
* **New Resource:** `sonatyperepo_repository_pypi_proxy` [GH-37]
* **New Resource:** `sonatyperepo_repository_r_group` [GH-49]
* **New Resource:** `sonatyperepo_repository_r_hosted` [GH-49]
* **New Resource:** `sonatyperepo_repository_r_proxy` [GH-49]
* **New Resource:** `sonatyperepo_repository_ruby_gems_group` [GH-48]
* **New Resource:** `sonatyperepo_repository_ruby_gems_hosted` [GH-48]
* **New Resource:** `sonatyperepo_repository_ruby_gems_proxy` [GH-48]
* **New Resource:** `sonatyperepo_repository_yum_group` [GH-52]
* **New Resource:** `sonatyperepo_repository_yum_hosted` [GH-52]
* **New Resource:** `sonatyperepo_repository_yum_proxy` [GH-52]

## 0.1.0 July 11, 2025

This is the first MVP release.

FEATURES:

* Resources and Datasources defined in [GH-1]
