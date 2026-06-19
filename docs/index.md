---
layout: ""
page_title: "Provider: Sonatype Nexus Repository"
description: |-
  The Sonatype Nexus Repository provider provides resources to interact with a Sonatype Nexus Repository installation.
---

# Sonatype Nexus Repository Provider

The `sonatyperepo` provider is used to interact with resources supported by [Sonatype Nexus Repository](https://www.sonatype.com/products/sonatype-nexus-repository). 

The provider needs to be configured with the proper credentials before it can be used.

## Compatability

This Provider is tested on Sonatype Nexus Repository Manager versions that have not yet entered Extendend Maintenance. 

See [Sonatype Nexus Repository 3 Versions Status](https://help.sonatype.com/en/sonatype-nexus-repository-3-versions-status.html) for details.

Sonatype Nexus Repository must not be in read-only mode in order to use this Provider. This will be checked. 
		
Some resources and features depend on the version of Sonatype Nexus Repository you are running. See individual Data Source and Resource documentaiton for details.

## Example Usage

```terraform
# Simplest Configuration
provider "sonatyperepo" {
  url      = "https://my-sonatype-nexus-repository.tld:port"
  username = "username"
  password = "password"
}

# Using environment variables for credentials (useful for CI/CD)
# Set NXRM_SERVER_URL, NXRM_SERVER_USERNAME, and NXRM_SERVER_PASSWORD
provider "sonatyperepo" {
  # Credentials provided via environment variables
}

# Mix environment variables with explicit configuration
provider "sonatyperepo" {
  username = "terraform-user"
  password = "terraform-password"
  # URL provided via NXRM_SERVER_URL environment variable
}

# If you run with a base path, you can add it:
provider "sonatyperepo" {
  url           = "https://my-sonatype-nexus-repository.tld:port"
  username      = "username"
  password      = "password"
  api_base_path = "/my-custom-base/service/rest"
}

# If you access via a Load Balancer or service that strips the `Server` header
# you can provide a hint as to the version of Sonatype Nexus Repository:
provider "sonatyperepo" {
  url          = "https://my-sonatype-nexus-repository.tld:port"
  username     = "username"
  password     = "password"
  version_hint = "3.89.1-01 (PRO)"
}
```

## Environment Variables

The provider supports the following environment variables for authentication and configuration. These are particularly useful for CI/CD scenarios where you don't want to store credentials in Terraform configuration files.

### Supported Environment Variables

| Environment Variable | Description | Provider Argument |
|---------------------|-------------|-------------------|
| `NXRM_SERVER_URL` | Sonatype Nexus Repository Server URL | `url` |
| `NXRM_SERVER_USERNAME` | Username for authentication | `username` |
| `NXRM_SERVER_PASSWORD` | Password for authentication | `password` |

### Precedence

Environment variables are evaluated first and will be overridden if the corresponding provider argument is explicitly set in your Terraform configuration. This allows you to:

- Set defaults via environment variables
- Override specific values in your Terraform configuration when needed

### CI/CD Example

In your CI/CD pipeline, set environment variables:

```bash
export NXRM_SERVER_URL="https://nexus.example.com"
export NXRM_SERVER_USERNAME="${NEXUS_USERNAME}"  # from CI/CD secret
export NXRM_SERVER_PASSWORD="${NEXUS_PASSWORD}"  # from CI/CD secret
terraform plan
```

> [!TIP]
> When using environment variables for all required fields, you can leave the provider block empty or omit credentials entirely. The provider will use environment variable values automatically.

## Required Privileges

The user account used to authenticate with Sonatype Nexus Repository must have appropriate privileges. Different Terraform operations require different privilege levels.

### Provider Initialization

When the provider initializes (during any Terraform operation), it performs these checks:

1. **Writable Status Check** — Validates the NXRM instance is not in read-only mode
   - Endpoint: `GET /service/rest/v1/status/writable`
   - Privilege: Access to status API (available to authenticated users)

2. **Cluster Node Detection** — Detects if running against an HA cluster to apply synchronization delays
   - Endpoint: `GET /service/rest/v1/status/check`
   - **Required Privilege:** `nx-metrics-all`

> [!IMPORTANT]
> The `nx-metrics-all` privilege is required even for read-only `terraform plan` operations because the provider needs to detect cluster topology during initialization.

### Privileges by Operation

#### terraform plan (Read-Only)

To run `terraform plan`, the user needs:
- `nx-metrics-all` — For cluster status detection
- Read privileges for each resource type being examined (e.g., `nx-repository-view-*-read`, `nx-blobstore-read`, `nx-security-read`)
- No write/edit/delete privileges required

#### terraform apply (Write Operations)

To run `terraform apply`, the user needs:
- All privileges required for `terraform plan`
- Write privileges for resources being managed (e.g., `nx-repository-view-*-edit`, `nx-blobstore-all`, `nx-security-all`)
- Administrative privileges for system-level configurations

### Creating a Read-Only Service Account

For CI/CD pipelines running `terraform plan` in pull request checks, create a dedicated service account:

**Step 1: Create Role**

Navigate to **Administration → Security → Roles → Create role**:

| Field | Value |
|-------|-------|
| Role ID | `terraform-plan-readonly` |
| Name | `Terraform Plan Read-Only` |
| Description | `Minimal privileges for terraform plan operations` |

**Step 2: Add Privileges**

Add these privileges to the role:
- `nx-metrics-all` — Required for cluster detection
- `nx-repository-view-*-read` — Read repository configurations
- `nx-blobstore-read` — Read blob store configurations
- `nx-security-read` — Read users, roles, privileges
- `nx-settings-read` — Read system settings
- `nx-routingrule-read` — Read routing rules
- `nx-contentselector-read` — Read content selectors

**Step 3: Create User**

Navigate to **Administration → Security → Users → Create user** and assign the `terraform-plan-readonly` role.

**Step 4: Use in CI/CD**

```bash
# Example: Set environment variables in your CI/CD pipeline
export NXRM_SERVER_URL="https://nexus.example.com"
export NXRM_SERVER_USERNAME="${NXRM_TERRAFORM_PLAN_USER}"  # from CI/CD secrets
export NXRM_SERVER_PASSWORD="${NXRM_TERRAFORM_PLAN_PASSWORD}"  # from CI/CD secrets

# Then run terraform plan
terraform plan
```

### Troubleshooting Permission Issues

#### Error: "Unable to check Sonatype Nexus Repository Cluster Node Count"

**Cause:** User lacks the `nx-metrics-all` privilege.

**Solution:** Add `nx-metrics-all` to the user's role:
1. Navigate to **Administration → Security → Roles**
2. Select the role assigned to your Terraform user
3. Add the `nx-metrics-all` privilege
4. Save the role

#### Error: "403 Forbidden" for specific resources

**Cause:** User lacks read (for plan) or write (for apply) privileges for that resource type.

**Solution:** Add the appropriate privilege based on the resource type:

| Resource Type | Read Privilege | Write Privileges |
|--------------|----------------|------------------|
| Repositories | `nx-repository-view-<format>-read` | `nx-repository-view-<format>-add/edit/delete` |
| Blob Stores | `nx-blobstore-read` | `nx-blobstore-all` |
| Users/Roles/Privileges | `nx-security-read` | `nx-security-all` |
| System Settings | `nx-settings-read` | `nx-settings-all` |
| Routing Rules | `nx-routingrule-read` | `nx-routingrule-all` |
| Content Selectors | `nx-contentselector-read` | `nx-contentselector-all` |

Replace `<format>` with the repository format (e.g., `maven2`, `npm`, `docker`) or use `*` for all formats.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `password` (String, Sensitive) Password for your user for Sonatype Nexus Repository Server. Can also be set using the `NXRM_SERVER_PASSWORD` environment variable.
- `url` (String) Sonatype Nexus Repository Server URL. Can also be set using the `NXRM_SERVER_URL` environment variable.
- `username` (String) Username for Sonatype Nexus Repository Server, requires role/permissions scoped to the resources you wish to manage. Can also be set using the `NXRM_SERVER_USERNAME` environment variable.

### Optional

- `api_base_path` (String) Base Path at which the API is present - defaults to `/service/rest`. This only needs to be set if you run Sonatype Nexus Repository at a Base Path that is not `/`.
- `cluster_stabilisation_delay_ms` (Number) Delay after write requests to allow for multi-node Cluster events to be processed by all Nodes before read requests. Only applies when running against a cluster with >1 active node.
				
> [!NOTE]
> Only set this if you are experiencing issues - the default value (10000) should suffice for most scenarios.
- `version_hint` (String) You can set this to the full version string (e.g. "3.85.0-03 (PRO)" or "3.80.0-06 (OSS)") of Sonatype Nexus Repository that you are connecting to.

> [!NOTE] 
> You can find the full version string in _Admin -> Support -> System Information_.
>				
> By default, this provider will attempt to automatically determine the version of Sonatype Nexus Repository you are connected to - but in some 
> real world cases, a Load Balancer or such may strip the HTTP Header that contians this information (the _Server_ header).

> [!TIP]
> If you receive an error such as `Plan is not supported for Sonatype Nexus Repository Manager: 0.0.0-0 (PRO=false)` then you should set 
> this attribute - otherwise, do not supply this attribute.