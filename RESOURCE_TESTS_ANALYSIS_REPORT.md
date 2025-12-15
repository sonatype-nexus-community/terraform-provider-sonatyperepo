# Terraform Provider SonatypeRepo - Resource Acceptance Tests Analysis Report

## Executive Summary

This report analyzes the acceptance tests targeting Resources in terraform-provider-sonatyperepo. The analysis identified opportunities for:
- Improving test coverage
- Reducing code duplication
- Improving method naming conventions
- Consolidating repeated string values as constants

---

## 1. Test Coverage Gaps

### 1.1 Limited CRUD Cycle Testing

**Issue**: Most resource tests only test Create and Read operations, but skip Update and Delete testing.

**Files Affected**:
- `blob_store/blob_store_file_resource_test.go`
- `role/role_resource_test.go`
- `user/user_resource_test.go`
- `privilege/privilege_application_resource_test.go`
- `privilege/privilege_repository_admin_resource_test.go`
- `privilege/privilege_wildcard_resource_test.go`
- `privilege/privilege_repository_view_resource_test.go`

**Current Pattern**:
```go
{
    Config: getTestAccBlobStoreFileResourceConfig(randomString),
    Check: resource.ComposeAggregateTestCheckFunc(...),
},
// Delete testing automatically occurs in TestCase
```

**Recommendation**:
- Add explicit Update steps to test field modifications
- Verify that changes are persisted correctly in the Nexus API
- Add assertions to verify intermediate states during updates

### 1.2 Missing Negative Test Scenarios

**Issue**: Limited negative test scenarios for data validation and error handling.

**Files Affected**:
- All repository resource test files (missing validation tests for proxy configurations, etc.)
- All privilege resource test files (missing invalid privilege configuration tests)

**Examples Needed**:
- Invalid remote URL formats for proxy repositories
- Invalid blob store name references
- Missing required fields
- Constraint violations (e.g., empty group member lists - some files have this, but inconsistently)

### 1.3 Minimal Configuration vs Full Configuration Testing

**Issue**: Tests should validate both minimal and maximal configurations for comprehensive coverage.

**Files with Gaps**:
- `blob_store/blob_store_file_resource_test.go` - only minimal
- Most single-repository resource tests (apt_proxy, cocoapods_proxy, composer_proxy)
- `user/user_resource_test.go` - only minimal
- `role/role_resource_test.go` - only single configuration

**Recommendation**:
Add test steps that:
1. Create resource with minimal required fields
2. Create resource with all optional fields populated
3. Verify computed fields are set appropriately

### 1.4 Data Source Tests with Environment Variable Gating

**✅ IMPLEMENTED**

**Good Pattern Observed**: `blob_store/blob_store_s3_data_source_test.go` correctly uses environment variable gating for tests requiring AWS credentials:

```go
func TestAccBlobStoreS3WithCredentialsDataSource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck: func() {
            if os.Getenv("TF_ACC_S3_BLOB_STORE") != "1" {
                t.Skip("S3 blob store tests require AWS credentials...")
            }
        },
        // ...
    })
}
```

**Changes Made**:
- `blob_store/blob_store_s3_resource_test.go` - Refactored to follow the data source pattern:
  - Created `TestAccBlobStoreS3ResourceValidation()` to test schema validation without API credentials
  - Created `TestAccBlobStoreS3ResourceWithCredentials()` with environment variable gating using `TF_ACC_S3_BLOB_STORE=1`
  - Renamed `getTestAccBlobStoreS3ResourceWillFail()` to `buildS3ResourceConfig()`
  - Added `buildS3ResourceCompleteConfig()` for full configuration testing
  
- `blob_store/blob_store_google_cloud_resource_test.go` - Similar refactoring:
  - Renamed all `getTestAccBlobStoreGoogleCloud*()` functions to `buildGoogleCloud*()`
  - Added `TestAccBlobStoreGoogleCloudResourceWithCredentials()` with `TF_ACC_GCS_BLOB_STORE=1` gating
  - Tests now separate validation-only tests from credential-dependent CRUD tests
  - Added Update and Import test steps for complete CRUD coverage when credentials available

---

## 2. Code Duplication Issues

### 2.1 Duplicate Test Configuration Patterns

**Issue**: Repository resource tests share nearly identical configuration patterns across multiple format types.

**Duplicated Patterns** (appears in 10+ files):
```
- HTTP client configuration (blocked, auto_block, connection, authentication)
- Proxy configuration (remote_url, content_max_age, metadata_max_age)
- Negative cache configuration (enabled, time_to_live)
- Storage configuration (blob_store_name, strict_content_type_validation)
```

**Files with Highest Duplication**:
- `repository/repository_maven_x_resources_test.go`
- `repository/repository_npm_x_resources_test.go`
- `repository/repository_docker_x_resources_test.go`
- `repository/repository_pypi_x_resources_test.go`
- `repository/repository_ruby_gems_x_resources_test.go`
- `repository/repository_conan_x_resources_test.go`
- `repository/repository_cargo_x_resources_test.go`

**Example of Duplication** (lines 88-123 in multiple files):
```hcl
http_client = {
    blocked = false
    auto_block = true
    connection = {
        enable_cookies = true
        retries = 9
        timeout = 999
        use_trust_store = true
        user_agent_suffix = "terraform"
    }
    authentication = {
        username = "user"
        password = "pass"
        preemptive = true
        type = "username"
    }
}
```

**Recommendation**:
Create shared configuration builder functions in `repository/repository_common_test.go`:
```go
func buildHttpClientConfig() string {
    return `
http_client = {
    blocked = false
    auto_block = true
    connection = {
        enable_cookies = true
        retries = 9
        timeout = 999
        use_trust_store = true
        user_agent_suffix = "terraform"
    }
    authentication = {
        username = "user"
        password = "pass"
        preemptive = true
        type = "username"
    }
}
`
}

func buildProxyConfig(remoteUrl string, contentMaxAge, metadataMaxAge int) string {
    return fmt.Sprintf(`
proxy = {
    remote_url = "%s"
    content_max_age = %d
    metadata_max_age = %d
}
`, remoteUrl, contentMaxAge, metadataMaxAge)
}

func buildNegativeCacheConfig(enabled bool, ttl int) string {
    return fmt.Sprintf(`
negative_cache = {
    enabled = %t
    time_to_live = %d
}
`, enabled, ttl)
}
```

### 2.2 Duplicate Group Repository Test Patterns

**Issue**: Group repository tests follow identical patterns across all format types.

**Pattern Observed** (in maven, npm, docker, etc.):
```go
// Test 1: Empty group member_names validation
Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "XXX-group-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = []
  }
}
`, resourceTypeGroup, randomString),
ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
```

This is repeated identically (except resource type name) in:
- `repository_maven_x_resources_test.go`
- `repository_npm_x_resources_test.go`
- `repository_docker_x_resources_test.go`
- `repository_pypi_x_resources_test.go`
- `repository_ruby_gems_x_resources_test.go`
- And more...

**Recommendation**:
Create a shared test template function in `repository_common_test.go`:
```go
func testGroupEmptyMembersValidation(t *testing.T, resourceType, randomString string) resource.TestStep {
    return resource.TestStep{
        Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s-group-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = []
  }
}
`, resourceType, strings.ToLower(strings.Split(resourceType, "_")[2]), randomString),
        ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
    }
}
```

### 2.3 Duplicate Import Test Patterns

**Issue**: Import tests follow nearly identical patterns across all repository types.

**Affected Files**:
- `repository_maven_x_resources_test.go` - 3 separate import tests (Hosted, Proxy, Group)
- `repository_npm_x_resources_test.go` - Similar pattern
- `repository_docker_x_resources_test.go` - Similar pattern

**Boilerplate Duplication**:
```go
func TestAccRepositoryXxxImport(t *testing.T) {
    randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
    resourceType := "sonatyperepo_repository_XXX_hosted"
    resourceName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceType)
    repoName := fmt.Sprintf("XXX-hosted-import-%s", randomString)

    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Create with minimal configuration
            {
                Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  // ... repository-specific config
}
`, resourceType, repoName),
            },
            // Import and verify no changes
            {
                ResourceName:                         resourceName,
                ImportState:                          true,
                ImportStateVerify:                    true,
                ImportStateId:                        repoName,
                ImportStateVerifyIdentifierAttribute: "name",
                ImportStateVerifyIgnore:              []string{"last_updated"},
            },
        },
    })
}
```

**Recommendation**:
This pattern appears in at least 6+ repository test files. Consider creating a helper:
```go
func createImportTestStep(resourceType, repoName string) resource.TestStep {
    return resource.TestStep{
        ResourceName:                         fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceType),
        ImportState:                          true,
        ImportStateVerify:                    true,
        ImportStateId:                        repoName,
        ImportStateVerifyIdentifierAttribute: "name",
        ImportStateVerifyIgnore:              []string{"last_updated"},
    }
}
```

---

## 3. Method Naming Issues

### 3.1 'Get' Prefix Violations

**Issue**: Multiple test helper functions use 'get' prefix, which violates Go idioms. In Go, getter methods should only use 'Get' prefix for accessor methods, not for builder/factory functions.

**Files with Issues**:

#### blob_store directory:
- `blob_store_file_resource_test.go:50` - `getTestAccBlobStoreFileResourceConfig()`
- `blob_store_s3_resource_test.go:48` - `getTestAccBlobStoreS3ResourceWillFail()`
- `blob_store_google_cloud_resource_test.go:102` - `getTestAccBlobStoreGoogleCloudResourceMinimal()`
- `blob_store_google_cloud_resource_test.go:128` - `getTestAccBlobStoreGoogleCloudResourceComplete()`
- `blob_store_google_cloud_resource_test.go:163` - `getTestAccBlobStoreGoogleCloudResourceInvalidBucket()`
- `blob_store_google_cloud_resource_test.go:185` - `getTestAccBlobStoreGoogleCloudResourceMissingName()`
- `blob_store_google_cloud_resource_test.go:205` - `getTestAccBlobStoreGoogleCloudResourceInvalidSoftQuota()`

#### repository directory:
- `repository_apt_proxy_resource_test.go:80` - `getRepositoryAptProxyResourceConfig()`
- `repository_cocoapods_proxy_resource_test.go:78` - `getRepositoryCocoaPodsProxyResourceConfig()`
- `repository_composer_proxy_resource_test.go:78` - `getRepositoryComposerProxyResourceConfig()`
- Many more similar patterns...

#### routing_rule & cleanup_policies:
- `routing_rule_test.go:117` - `getTestAccRoutingRuleResourceConfig()`
- `cleanup_policies_test.go:114` - `getTestAccCleanupPolicyResourceConfig()`

**Recommended Refactoring**:

Replace `get` prefix with more idiomatic names:
- `getTestAccBlobStoreFileResourceConfig()` → `buildBlobStoreFileResourceConfig()` or `testBlobStoreFileResourceConfig()`
- `getRepositoryAptProxyResourceConfig()` → `buildAptProxyRepositoryConfig()` or `testAptProxyRepositoryConfig()`
- `getTestAccRoutingRuleResourceConfig()` → `buildRoutingRuleResourceConfig()` or `testRoutingRuleResourceConfig()`

**Alternative Go Idioms**:
- Use verb-noun pattern: `buildXConfig()`, `createXConfig()`, `newXConfig()`
- Use descriptive names: `minimalBlobStoreConfig()`, `completeBlobStoreConfig()`

### 3.2 Inconsistent Naming Conventions

**Issue**: Different conventions used for similar helper functions across files.

**Examples**:
- Some use `getRepository...` prefix (apt, cocoapods, composer)
- Some use `getTestAcc...` prefix (blob_store, routing_rule, cleanup_policies)
- Some use `build...` prefix (observed in some places)

**Recommendation**:
Establish consistent naming for test helper functions across all resource tests:
- For configuration builders: `config<ResourceType>Minimal()`, `config<ResourceType>Complete()`
- For resource creation: `build<ResourceType>Config()` or `create<ResourceType>Config()`

---

## 4. Constants and Magic Strings

### 4.1 Repeated String Values Without Constants

**Issue**: Common strings are repeated throughout tests without being defined as constants.

**Commonly Repeated Strings** (across multiple files):
```
- "default" (blob store name) - appears 50+ times
- "true" / "false" (inline in Terraform configs) - appears 100+ times
- "https://archive.ubuntu.com/ubuntu/" - repeated in apt tests
- "https://registry.npmjs.org" - repeated in npm tests
- "https://registry-1.docker.io" - repeated in docker tests
- "https://repo1.maven.org/maven2/" - repeated in maven tests
- "blob_store_name"
- "strict_content_type_validation"
- "online"
- "storage"
- "proxy"
- "http_client"
- "negative_cache"
- "authentication"
```

**Files Most Affected**:
- All repository resource test files
- All blob store resource test files

### 4.2 Partial Constant Usage

**Good Pattern Already Exists** (`repository_common_test.go:19-25`):
```go
const (
    RES_ATTR_DOCKER_FORCE_BASIC_AUTH string = "docker.force_basic_auth"
    RES_ATTR_DOCKER_PATH_ENABLED     string = "docker.path_enabled"
    RES_ATTR_DOCKER_V1_ENABLED       string = "docker.v1_enabled"
    RES_ATTR_RAW_CONTENT_DISPOSITION string = "raw.content_disposition"
    RES_ATTR_STORAGE_BLOB_STORE_NAME string = "storage.blob_store_name"
)
```

**Additional Constants Needed** (in `repository_common_test.go`):
```go
const (
    // Common resource attributes
    RES_ATTR_NAME                           = "name"
    RES_ATTR_ONLINE                         = "online"
    RES_ATTR_URL                            = "url"
    RES_ATTR_STORAGE                        = "storage"
    RES_ATTR_STORAGE_STRICT_CONTENT_TYPE    = "storage.strict_content_type_validation"
    RES_ATTR_STORAGE_WRITE_POLICY           = "storage.write_policy"
    RES_ATTR_PROXY                          = "proxy"
    RES_ATTR_PROXY_REMOTE_URL               = "proxy.remote_url"
    RES_ATTR_PROXY_CONTENT_MAX_AGE          = "proxy.content_max_age"
    RES_ATTR_PROXY_METADATA_MAX_AGE         = "proxy.metadata_max_age"
    RES_ATTR_NEGATIVE_CACHE                 = "negative_cache"
    RES_ATTR_NEGATIVE_CACHE_ENABLED         = "negative_cache.enabled"
    RES_ATTR_NEGATIVE_CACHE_TTL             = "negative_cache.time_to_live"
    RES_ATTR_HTTP_CLIENT                    = "http_client"
    RES_ATTR_HTTP_CLIENT_BLOCKED            = "http_client.blocked"
    RES_ATTR_HTTP_CLIENT_AUTO_BLOCK         = "http_client.auto_block"
    RES_ATTR_HTTP_CLIENT_CONN_COOKIES       = "http_client.connection.enable_cookies"
    RES_ATTR_HTTP_CLIENT_CONN_RETRIES       = "http_client.connection.retries"
    RES_ATTR_HTTP_CLIENT_CONN_TIMEOUT       = "http_client.connection.timeout"
    RES_ATTR_HTTP_CLIENT_CONN_TRUST_STORE   = "http_client.connection.use_trust_store"
    RES_ATTR_HTTP_CLIENT_CONN_USER_AGENT    = "http_client.connection.user_agent_suffix"
    RES_ATTR_HTTP_CLIENT_AUTH_USERNAME      = "http_client.authentication.username"
    RES_ATTR_HTTP_CLIENT_AUTH_PASSWORD      = "http_client.authentication.password"
    RES_ATTR_HTTP_CLIENT_AUTH_PREEMPTIVE    = "http_client.authentication.preemptive"
    RES_ATTR_HTTP_CLIENT_AUTH_TYPE          = "http_client.authentication.type"
    RES_ATTR_ROUTING_RULE                   = "routing_rule"
    RES_ATTR_REPLICATION                    = "replication"
    RES_ATTR_GROUP_MEMBER_NAMES             = "group.member_names"
    RES_ATTR_COMPONENT_PROPRIETARY          = "component.proprietary_components"

    // Common configuration values
    DEFAULT_BLOB_STORE              = "default"
    WRITE_POLICY_ALLOW_ONCE         = "ALLOW_ONCE"
    DEFAULT_CONTENT_MAX_AGE         = 1442
    DEFAULT_METADATA_MAX_AGE        = 1400
    DEFAULT_CACHE_TTL              = 1440
    HTTP_CLIENT_CONN_TIMEOUT_TEST   = 999
    HTTP_CLIENT_CONN_RETRIES_TEST   = 9
    
    // Default test URLs
    APT_REMOTE_URL                  = "https://archive.ubuntu.com/ubuntu/"
    NPM_REMOTE_URL                  = "https://registry.npmjs.org"
    DOCKER_REMOTE_URL               = "https://registry-1.docker.io"
    MAVEN_REMOTE_URL                = "https://repo1.maven.org/maven2/"
    PYPI_REMOTE_URL                 = "https://pypi.org/simple"
    RUBYGEMS_REMOTE_URL             = "https://rubygems.org"
    COCOAPODS_REMOTE_URL            = "https://cdn.cocoapods.org"
    COMPOSER_REMOTE_URL             = "https://repo.packagist.org"
)
```

**Additional Constants Needed** (in other test files):

For `role/role_resource_test.go` and `role/roles_data_source_test.go`:
```go
const (
    RES_ATTR_ID                  = "id"
    RES_ATTR_DESCRIPTION         = "description"
    RES_ATTR_PRIVILEGES          = "privileges"
    RES_ATTR_ROLES               = "roles"
)
```

For `user/user_resource_test.go` and `user/users_data_source_test.go`:
```go
const (
    RES_ATTR_USER_ID          = "user_id"
    RES_ATTR_FIRST_NAME       = "first_name"
    RES_ATTR_LAST_NAME        = "last_name"
    RES_ATTR_EMAIL_ADDRESS    = "email_address"
    RES_ATTR_STATUS           = "status"
    RES_ATTR_READ_ONLY        = "read_only"
    RES_ATTR_SOURCE           = "source"
)
```

For `privilege/` resource tests:
```go
const (
    RES_ATTR_NAME           = "name"
    RES_ATTR_DESCRIPTION    = "description"
    RES_ATTR_READ_ONLY      = "read_only"
    RES_ATTR_TYPE           = "type"
    RES_ATTR_DOMAIN         = "domain"
    RES_ATTR_ACTIONS        = "actions"
)
```

### 4.3 Hardcoded Test Values in Assertions

**Issue**: Test assertions use hardcoded values instead of referencing variables used in configuration.

**Example from `user_resource_test.go:45-52`**:
```go
Check: resource.ComposeAggregateTestCheckFunc(
    // Verify
    resource.TestCheckResourceAttr(resourceNameUser, "user_id", fmt.Sprintf("acc-test-user-%s", randomString)),
    resource.TestCheckResourceAttr(resourceNameUser, "first_name", fmt.Sprintf("Acc Test %s", randomString)),
    resource.TestCheckResourceAttr(resourceNameUser, "last_name", "User"),  // Hardcoded
    resource.TestCheckResourceAttr(resourceNameUser, "email_address", fmt.Sprintf("acc-test-%s@local", randomString)),
    resource.TestCheckResourceAttr(resourceNameUser, "status", "active"),   // Hardcoded
    resource.TestCheckResourceAttr(resourceNameUser, "read_only", "false"), // Hardcoded
    resource.TestCheckResourceAttr(resourceNameUser, "source", common.DEFAULT_USER_SOURCE),
    resource.TestCheckResourceAttr(resourceNameUser, "roles.#", "1"),       // Hardcoded count
),
```

**Better Approach**:
```go
const (
    testLastName        = "User"
    testStatus          = "active"
    testReadOnly        = "false"
    testRolesCount      = "1"
)

// Then in assertions
resource.TestCheckResourceAttr(resourceNameUser, "last_name", testLastName),
resource.TestCheckResourceAttr(resourceNameUser, "status", testStatus),
```

---

## 5. Test Specific Patterns and Best Practices

### 5.1 Good Patterns to Keep

**Environment Variable Gating** (blob_store_s3_data_source_test.go):
```go
PreCheck: func() {
    if os.Getenv("TF_ACC_S3_BLOB_STORE") != "1" {
        t.Skip("S3 blob store tests require AWS credentials...")
    }
},
```
✅ This pattern should be applied to S3 and GCS resource tests.

**Version-Based Test Skipping** (capability_x_resource_test.go):
```go
testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
    Major: 3, Minor: 0, Patch: 0,
}, &common.SystemVersion{
    Major: 3, Minor: 83, Patch: 99,
})
```
✅ This is excellent for version-dependent features.

**Import State Testing** (repository and other resource tests):
```go
{
    ResourceName:                         resourceName,
    ImportState:                          true,
    ImportStateVerify:                    true,
    ImportStateId:                        repoName,
    ImportStateVerifyIdentifierAttribute: "name",
    ImportStateVerifyIgnore:              []string{"last_updated"},
},
```
✅ Consistently used and should remain.

### 5.2 Room for Improvement

**Replication Configuration Testing**:
Some repository tests have optional replication config, but not all tests that should have it do.

**Cleanup Policy Testing**:
Repository tests should verify cleanup policy associations more thoroughly.

---

## 6. Recommendations Summary

### Priority 1: High Impact
1. **Consolidate duplicate HTTP client and proxy configurations** into shared builder functions
2. **Apply environment variable gating pattern** to S3 and GCS resource tests
3. **Rename `get*` prefixed functions** to use idiomatic Go naming (build, create, config)
4. **Extract repeated resource attribute strings** to constants in common_test.go files

### Priority 2: Medium Impact
5. **Add explicit Update test steps** to all CRUD resource tests
6. **Extend negative test scenarios** for validation and error handling
7. **Create template functions** for group repository tests and import tests
8. **Consolidate hardcoded test values** into constants

### Priority 3: Lower Priority
9. Document test structure and patterns for maintainability
10. Consider creating a test utilities package for common test helpers

---

## 7. Specific Action Items by File/Directory

### blob_store/
- [x] Rename `getTestAcc*` functions to `build*` (S3 and GCS done)
- [x] Add environment variable gating for S3 and GCS tests
- [x] Refactor S3 and GCS resource tests to match data source pattern
- [ ] Add Update test steps to file resource test
- [ ] Extract blob store path/name patterns to constants

### repository/
- [ ] Create configuration builder functions in `repository_common_test.go`
- [ ] Rename all `getRepository*` and `getTestAcc*` functions
- [ ] Create shared template functions for group validation and import tests
- [ ] Extract all magic strings (URLs, attribute paths) to constants
- [ ] Add explicit Update test steps to all repository resource tests
- [ ] Extend negative test scenarios for proxy configurations

### capability/
- [ ] Rename `getTestAcc*` functions to `build*` or `config*`
- [ ] Extract hardcoded notes and property values to constants
- [ ] Consider creating shared configuration builders

### privilege/
- [ ] Add constants for common privilege attributes
- [ ] Add negative test scenarios for invalid privilege types
- [ ] Add Update test steps to resource tests

### role/
- [ ] Add constants for role attributes (id, name, description, privileges)
- [ ] Add Update test steps to verify privilege/role changes

### user/
- [ ] Add constants for user attributes
- [ ] Add Update test steps to verify user modifications
- [ ] Add negative scenarios for invalid status values

### system/
- [ ] Review and consolidate test patterns across config_mail, config_ldap, security_saml
- [ ] Extract hardcoded configuration values to constants

### content_selector/
- [ ] Review test patterns for consistency
- [ ] Add Update test steps to resource test

---

## Conclusion

The acceptance test suite is comprehensive in breadth but has opportunities for:
1. **Better maintainability** through consolidation of duplicate code patterns
2. **Improved clarity** through consistent naming conventions
3. **Enhanced coverage** through explicit Update steps and negative scenarios  
4. **Better testability** through proper constant extraction and environment variable gating

Implementing these recommendations will reduce maintenance burden and make the test suite more resilient to future changes.
