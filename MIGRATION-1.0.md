# THIS IS TO BE FORMATTED

## Breaking Changes in 1.0.0

### Configuration for Sonatype Repository Firewall

Prior to 1.0.0 of this provider, the consumer of this provider had to manage both a _proxy_ repository and a _capability_ in order to configure 
the Sonatype Repository Firewall for that proxy repository. This was only possible when running Sonatype Nexus Repository 3.84.0 or newer.

This was not ideal for two reasons:
1. You could configure Sonatype Repository Firewall for a proxy repository without Sonatype Nexus Repository being connected to a valid Sonatype IQ Server - hence it wouldn't actually function
2. The delcarative configuration required for Terraform did not shield users from the internal requirements sufficiently

Since 1.0.0 - configuration of Sonatype Repository Firewall is now handled within the `sonatyperepo_repository_*_proxy` resources themselves and there is no requirement to manage a separate capability resource. The `sonatyperepo_capability_repository_firewall` resource has been deprecated.

Additionally - it is now required that a valid Sonatype IQ Connection is configured _PRIOR_ to mamnaging repository resources with Sonatype Repository Firewall configuration - use the `sonatyperepo_system_iq_connection` resource.

Example Terraform prior to 1.0.0:
```hcl
resource "sonatyperepo_repository_npm_proxy" "example" {
  name = "npm-proxy"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "%s"
    content_max_age = 1440
    metadata_max_age = 1440
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
  }

  # This config related to Repositort Firewall too!
  npm = {
    remove_quarrantined = true    
  }
}

resource "sonatyperepo_capability_firewall_audit_and_quarantine" "example" {
  notes      = "These are notes from Terraform"
  enabled    = true
  properties = {
    repository = sonatyperepo_repository_npm_proxy.example.name
    quarantine = true
  }
}
```

The equivalent in 1.0.0+ is now:

```hcl
resource "sonatyperepo_repository_npm_proxy" "example" {
  name = "npm-proxy"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "%s"
    content_max_age = 1440
    metadata_max_age = 1440
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
  }
  repository_firewall = {
    enabled = true
    quarantine = true
    pccs_enabled = true   # This replaces the `remove_quarrantined` property
  }
}
```

Not all proxy repository formats support Sonatype Repository Firewall or [Policy-Compliant Component Selection (PCCS)](https://help.sonatype.com/en/policy-compliant-component-selection.html) - 
see this provider's [documentation](https://registry.terraform.io/providers/sonatype-nexus-community/sonatyperepo/latest/) for more details.

Sonatype Nexus Repository 3.84.0 or newer is still required.

### Resource Renaming

The following resources have been renamed to improve consistency.

- Resource `sonatyperepo_repository_maven_group` has been renamed to `sonatyperepo_repository_maven2_group`
- Resource `sonatyperepo_repository_maven_hosted` has been renamed to `sonatyperepo_repository_maven2_hosted`
- Resource `sonatyperepo_repository_maven_proxy` has been renamed to `sonatyperepo_repository_maven2_proxy`
- Resource `sonatyperepo_repository_ruby_gems_group` has been renamed to `sonatyperepo_repository_rubygems_group`
- Resource `sonatyperepo_repository_ruby_gems_hosted` has been renamed to `sonatyperepo_repository_rubygems_hosted`
- Resource `sonatyperepo_repository_ruby_gems_proxy` has been renamed to `sonatyperepo_repository_rubygems_proxy`

### Resources Deprecated

- `sonatyperepo_capability_repository_firewall`


### Other Resource Schema Changes

- Resource `sonatyperepo_repository_apt_hosted`: 
  - The `apt_signing` block is now required
  - The `apt_signing.passphrase` field is now optional
- Resources `sonatyperepo_repository_*_proxy`: 
  - `proxy.content_max_age` is now optional and has a default value (`1440`)
  - `proxy.metadata_max_age` is now optional and has a default value (`1440`)
  - `negative_cache.enabled` is now optional and has a default value (`true`)
  - `negative_cache.time_to_live` is now optional and has a default value (`1440`)

## Improvements 

### Resources now supporting Import
- `sonatyperepo_repository_cargo_hosted`
- `sonatyperepo_repository_cargo_proxy`
- `sonatyperepo_repository_cocoapods_proxy`
- `sonatyperepo_repository_conan_hosted`
- `sonatyperepo_repository_composer_proxy`
- `sonatyperepo_repository_conan_proxy`
- `sonatyperepo_repository_gitlfs_hosted`
- `sonatyperepo_repository_r_proxy`
- `sonatyperepo_repository_r_hosted`
