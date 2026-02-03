/*
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package repository_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ------------------------------------------------------------
// Test Data Scenarios
// ------------------------------------------------------------
var proxyTestData = []repositoryProxyTestData{
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_APT_DISTRIBUTION, "bionic"),
			}
		},
		RemoteUrl:  TEST_DATA_APT_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_APT,
		SchemaFunc: repositoryProxyResourceConfig,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_CARGO_REQUIRE_AUTHENTICATION, "true"),
			}
		},
		RemoteUrl:  TEST_DATA_CARGO_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_CARGO,
		SchemaFunc: repositoryProxyResourceConfig,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_COCOAPODS_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_COCOAPODS,
		SchemaFunc: repositoryProxyResourceConfig,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_COMPOSER_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_COMPOSER,
		SchemaFunc: repositoryProxyResourceConfig,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_CONAN_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_CONAN,
		SchemaFunc: repositoryProxyResourceConfig,
	},
}

// ------------------------------------------------------------
// PROXY REPO TESTING (GENERIC)
// ------------------------------------------------------------
func TestAccRepositoryGenericProxyByFormat(t *testing.T) {
	for _, td := range proxyTestData {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(td.RepoFormat))
		resourceName := fmt.Sprintf(repoNameFString, resourceType)
		repoName := fmt.Sprintf(proxyNameFString, td.RepoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// 1. Create with minimal configuration relying on defaults
				{
					Config: td.SchemaFunc(resourceType, repoName, td.RepoFormat, td.RemoteUrl, randomString, false),
					Check: resource.ComposeAggregateTestCheckFunc(
						append(
							// Test Case Specific Checks
							td.CheckFunc(resourceName),

							// Generic Checks
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_NAME, repoName),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_ONLINE, "true"),
							resource.TestCheckResourceAttrSet(resourceName, RES_ATTR_URL),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
							// resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_CLEANUP_POLICY_COUNT, "0"),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_REMOTE_URL, td.RemoteUrl),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_CONTENT_MAX_AGE, fmt.Sprintf("%d", common.DEFAULT_PROXY_CONTENT_MAX_AGE)),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_METADATA_MAX_AGE, fmt.Sprintf("%d", common.DEFAULT_PROXY_METADATA_MAX_AGE)),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_NEGATIVE_CACHE_ENABLED, fmt.Sprintf("%t", common.DEFAULT_PROXY_NEGATIVE_CACHE_ENABLED)),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, fmt.Sprintf("%d", common.DEFAULT_PROXY_NEGATIVE_CACHE_TTL)),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_BLOCKED, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_BLOCKED)),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_AUTO_BLOCK)),
							resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS)),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES)),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, fmt.Sprintf("%d", common.DEFAULT_HTTP_CLIENT_CONNECTION_RETRIES)),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, fmt.Sprintf("%d", common.DEFAULT_HTTP_CLIENT_CONNECTION_TIMEOUT)),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE)),
							resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX),
							resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_ROUTING_RULE_NAME),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED, "false"),
						)...,
					),
				},
				// 2. Update to use full config
				{
					Config: td.SchemaFunc(resourceType, repoName, td.RepoFormat, td.RemoteUrl, randomString, true),
					Check: resource.ComposeAggregateTestCheckFunc(
						append(
							// Test Case Specific Checks
							td.CheckFunc(resourceName),

							// Generic Checks
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_NAME, repoName),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_ONLINE, "true"),
							resource.TestCheckResourceAttrSet(resourceName, RES_ATTR_URL),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
							// resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_CLEANUP_POLICY_COUNT, "0"),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_REMOTE_URL, td.RemoteUrl),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_CONTENT_MAX_AGE, TEST_DATA_TIMEOUT),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_METADATA_MAX_AGE, TEST_DATA_TIMEOUT),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_NEGATIVE_CACHE_ENABLED, "false"),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, TEST_DATA_TIMEOUT),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_BLOCKED, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_BLOCKED)),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, "false"),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, "true"),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, "true"),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, "2"),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, "59"),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, "true"),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX, "custom-suffix"),
							resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_ROUTING_RULE_NAME),
							resource.TestCheckResourceAttr(resourceName, RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED, "false"),
							// resource.TestCheckResourceAttr(resourceName, RES_ATTR_REPLICATION_ASSET_PATH_REGEX, ".*"),
						)...,
					),
				},
				// 3. Import and verify no changes
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
}

func TestAccRepositoryGenericProxyInvalidRemoteUrl(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(proxyNameFString, repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-remote-url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageInvalidRemoteUrl),
				},
			},
		})
	}
}

func TestAccRepositoryGenericProxyInvalidBlobStore(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(proxyNameFString, repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://some.source.url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageBlobStoreNotFound),
				},
			},
		})
	}
}

func TestAccRepositoryGenericProxyInvalidHttpConnectionRetries(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(proxyNameFString, repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// HTTP Connection Timeout to large
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-remote-url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
    connection = {
      retries = 11
    }
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
				},
				// HTTP Connection Timeout to small
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-remote-url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
    connection = {
      retries = -1
    }
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
				},
			},
		})
	}
}

func TestAccRepositoryGenericProxyInvalidHttpConnectionTimeout(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(proxyNameFString, repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// HTTP Connection Timeout to large
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-remote-url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
    connection = {
      timeout = 3601
    }
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
				},
				// HTTP Connection Timeout to small
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-remote-url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
    connection = {
      timeout = 0
    }
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
				},
			},
		})
	}
}

func TestAccRepositoryGenericProxyInvalidNegativeCacheTtl(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(proxyNameFString, repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://some.source.url"
  }
  negative_cache = {
    enabled = true
    time_to_live = -1
  }
  http_client = {
    blocked = false
    auto_block = true
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageNegativeCacheTimeoutValue),
				},
			},
		})
	}
}

func TestAccRepositoryGenericProxyMissingStorage(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(proxyNameFString, repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  proxy = {
    remote_url = "https://some.source.url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
    connection = {
      timeout = 0
    }
  }
  %s
 }`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageStorageRequired),
				},
			},
		})
	}
}

const (
	configBlockProxyDefaultApt    string = "apt = { distribution = \"bionic\" }"
	configBlockProxyDefaultCargo  string = "cargo = { require_authentication = true }"
	configBlockProxyDefaultConan  string = "conan = { conan_version = \"V2\" }"
	configBlockProxyDefaultDocker string = "docker = { force_basic_auth = false\nv1_enabled = false }\ndocker_proxy = { }"
	configBlockProxyDefaultMaven  string = "maven = { layout_policy = \"PERMISSIVE\"\nversion_policy = \"RELEASE\" }"
	configBlockProxyDefaultNuget  string = "nuget_proxy = { nuget_version = \"V3\" }"
	configBlockProxyDefaultRaw    string = "raw = { content_disposition = \"ATTACHMENT\" }"
)

func formatSpecificProxyDefaultConfig(repoFormat string) string {
	switch repoFormat {
	case common.REPO_FORMAT_APT:
		return configBlockProxyDefaultApt
	case common.REPO_FORMAT_CARGO:
		return configBlockProxyDefaultCargo
	case common.REPO_FORMAT_CONAN:
		return configBlockProxyDefaultConan
	case common.REPO_FORMAT_DOCKER:
		return configBlockProxyDefaultDocker
	case common.REPO_FORMAT_MAVEN:
		return configBlockProxyDefaultMaven
	case common.REPO_FORMAT_NUGET:
		return configBlockProxyDefaultNuget
	case common.REPO_FORMAT_RAW:
		return configBlockProxyDefaultRaw

	default:
		return ""
	}
}

func repositoryProxyResourceFullConfig(resourceType, repoName, remoteUrl, formatSpecificConfig string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "%s"
    content_max_age = 1439
    metadata_max_age = 1439
  }
  negative_cache = {
    enabled = false
    time_to_live = 1439
  }
  http_client = {
    blocked = false
    auto_block = false
    connection = {
	  enable_circular_redirects = true
	  enable_cookies = true
	  retries = 2
	  timeout = 59
	  use_trust_store = true
	  user_agent_suffix = "custom-suffix"
	}
  }
  replication = {
    preemptive_pull_enabled = false
  }
  %s
 }
`, resourceType, repoName, remoteUrl, formatSpecificConfig)
}

func repositoryProxyResourceMinimalConfigWithDefaults(resourceType, repoName, remoteUrl, formatSpecificConfig string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
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
  %s
 }
`, resourceType, repoName, remoteUrl, formatSpecificConfig)
}

func repositoryProxyResourceConfig(resourceType, repoName, repoFormat, remoteUrl, randomString string, completeData bool) string {
	configBlock := formatSpecificProxyDefaultConfig(repoFormat)
	if completeData {
		return repositoryProxyResourceFullConfig(
			resourceType, repoName, remoteUrl, configBlock,
		)
	} else {
		return repositoryProxyResourceMinimalConfigWithDefaults(
			resourceType, repoName, remoteUrl, configBlock,
		)
	}
}
