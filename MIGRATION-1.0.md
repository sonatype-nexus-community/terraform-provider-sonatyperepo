# Migrating from v0.x.x to v1.x.x

Release 1.x.x is a breaking-change release. This means that some ways of working / names and conventions have
changed and are not backwards compatible with 0.x.x releases.

Below is a summary of the key breaking changes with information on how you can update your Terraform HCL to work with 1.x.x releases.

- [Breaking Changes in 1.0.0](#breaking-changes-in-100)
  - [Configuration for Sonatype Repository Firewall](#configuration-for-sonatype-repository-firewall)
  - [Resource Renaming](#resource-renaming)
  - [Resources Deprecated](#resources-deprecated)
  - [Other Resource Schema Changes](#other-resource-schema-changes)
- [Improvements](#improvements)
  - [Resources now supporting Import](#resources-now-supporting-import)

## Breaking Changes in 1.0.0

### Configuration for Sonatype Repository Firewall

Prior to 1.0.0 of this provider, the consumer of this provider had to manage both a _proxy_ repository and a _capability_ in order to configure 
the Sonatype Repository Firewall for that proxy repository. This was only possible when running Sonatype Nexus Repository 3.84.0 or newer.

This was not ideal for two reasons:
1. You could configure Sonatype Repository Firewall for a proxy repository without Sonatype Nexus Repository being connected to a valid Sonatype IQ Server - hence it wouldn't actually function
2. The delcarative configuration required for Terraform did not shield users from the internal requirements sufficiently

Since 1.0.0 - configuration of Sonatype Repository Firewall is now handled within the `sonatyperepo_repository_*_proxy` resources themselves and there is no requirement to manage a separate capability resource. The `sonatyperepo_capability_firewall_audit_and_quarantine` resource has been deprecated.

Additionally - it is now required that a valid Sonatype IQ Connection is configured **_PRIOR_** to managing repository resources with Sonatype Repository Firewall configuration - use the `sonatyperepo_system_iq_connection` resource to ensure this is configured.

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
    remote_url = "https://npm.server.tld"
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

  # This config related to Repository Firewall too!
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
    remote_url = "https://npm.server.tld"
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

#### Migration Guide

You have **two options** for handling these resource renames:

##### Option 1: Continue Using Deprecated Names (No Action Required)

The old resource names are maintained as **deprecated aliases** for backward compatibility. Your existing Terraform configurations will continue to work without any changes. However, we recommend migrating to the new names (Option 2) as the deprecated names will be removed in a future major version.

##### Option 2: Migrate to New Resource Names (Recommended)

To migrate your existing resources to the new names, use Terraform's [`moved` blocks](https://developer.hashicorp.com/terraform/language/modules/develop/refactoring).

**Prerequisites:**
- Terraform 1.8 or later (required for `moved` blocks)
- Provider version v1.0.1 or later

**Complete Resource Name Mapping:**

| Old Name (Deprecated) | New Name |
|----------------------|----------|
| `sonatyperepo_repository_maven_group` | `sonatyperepo_repository_maven2_group` |
| `sonatyperepo_repository_maven_hosted` | `sonatyperepo_repository_maven2_hosted` |
| `sonatyperepo_repository_maven_proxy` | `sonatyperepo_repository_maven2_proxy` |
| `sonatyperepo_repository_ruby_gems_group` | `sonatyperepo_repository_rubygems_group` |
| `sonatyperepo_repository_ruby_gems_hosted` | `sonatyperepo_repository_rubygems_hosted` |
| `sonatyperepo_repository_ruby_gems_proxy` | `sonatyperepo_repository_rubygems_proxy` |

**Step-by-Step Migration Process:**

**Example: Migrating a Maven Hosted Repository**

1. **Current Configuration** (before migration):
   ```hcl
   resource "sonatyperepo_repository_maven_hosted" "my_repo" {
     name = "my-maven-repo"
     online = true
     storage = {
       blob_store_name = "default"
       strict_content_type_validation = true
       write_policy = "ALLOW"
     }
   }
   ```

2. **Update Resource Type** - Change the resource type to the new name:
   ```hcl
   resource "sonatyperepo_repository_maven2_hosted" "my_repo" {
     name = "my-maven-repo"
     online = true
     storage = {
       blob_store_name = "default"
       strict_content_type_validation = true
       write_policy = "ALLOW"
     }
   }
   ```

3. **Add `moved` Block** - Add this block to inform Terraform about the rename:
   ```hcl
   moved {
     from = sonatyperepo_repository_maven_hosted.my_repo
     to   = sonatyperepo_repository_maven2_hosted.my_repo
   }
   ```

4. **Verify Migration** - Run `terraform plan`:
   ```bash
   terraform plan
   ```

   You should see output similar to:
   ```
   # sonatyperepo_repository_maven_hosted.my_repo has moved to sonatyperepo_repository_maven2_hosted.my_repo
   resource "sonatyperepo_repository_maven2_hosted" "my_repo" {
     name = "my-maven-repo"
     # ... (no changes)
   }
   ```

   **Important:** Terraform should indicate the resource will be **moved**, not destroyed and recreated. If you see destroy/create operations, do not proceed - review your configuration.

5. **Apply Migration** - Execute the state migration:
   ```bash
   terraform apply
   ```

6. **Cleanup** - After successful migration, remove the `moved` block from your configuration. The block is no longer needed once the state has been migrated.

**Example: Migrating Multiple Resources**

If you have multiple repositories to migrate:

```hcl
# Maven resources
resource "sonatyperepo_repository_maven2_hosted" "releases" {
  name = "maven-releases"
  # ... configuration ...
}

resource "sonatyperepo_repository_maven2_proxy" "central" {
  name = "maven-central"
  # ... configuration ...
}

# RubyGems resources
resource "sonatyperepo_repository_rubygems_hosted" "gems" {
  name = "rubygems-hosted"
  # ... configuration ...
}

# Add moved blocks for all renamed resources
moved {
  from = sonatyperepo_repository_maven_hosted.releases
  to   = sonatyperepo_repository_maven2_hosted.releases
}

moved {
  from = sonatyperepo_repository_maven_proxy.central
  to   = sonatyperepo_repository_maven2_proxy.central
}

moved {
  from = sonatyperepo_repository_ruby_gems_hosted.gems
  to   = sonatyperepo_repository_rubygems_hosted.gems
}
```

**Troubleshooting:**

- **Error: "No resource schema found"** - Ensure you're using provider version 1.0.0 or later
- **Destroy/Create instead of Move** - Verify your `moved` block syntax matches the examples above
- **Terraform version error** - Upgrade to Terraform 1.8+ to use `moved` blocks
- **State already migrated** - If you see "moved to an address that doesn't exist", the migration may have already completed. Remove the `moved` block and run `terraform plan` again.

**Important Notes:**

- State migration is **safe** - no changes are made to your actual Nexus Repository resources
- The `moved` block only updates Terraform state, not the infrastructure
- You can migrate resources incrementally (one at a time) or all at once
- The deprecated resource names will be removed in a future major version (v2.0.0 or later)
- We recommend completing the migration during your next maintenance window

### Resources Deprecated

- `sonatyperepo_capability_firewall_audit_and_quarantine`


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
