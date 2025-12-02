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
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceNameCocoaPodsProxy = "sonatyperepo_repository_cocoapods_proxy.repo"
	resourceTypeCocoaPodsProxy = "sonatyperepo_repository_cocoapods_proxy"
)

func TestAccRepositoryCocoaPodsProxyResourceNoReplication(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: repositoryCocoaPodsProxyResourceMinimalConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify minimal config
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "name", fmt.Sprintf("cocoapods-proxy-repo-minimal-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNameCocoaPodsProxy, "url"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "proxy.remote_url", "https://cdn.cocoapods.org/"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "proxy.content_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.auto_block", "true"),
				),
			},
			// Update to full configuration
			{
				Config: repositoryCocoaPodsProxyResourceConfig(randomString, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify full config
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "name", fmt.Sprintf("cocoapods-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNameCocoaPodsProxy, "url"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "proxy.remote_url", "https://cdn.cocoapods.org/"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceNameCocoaPodsProxy, "routing_rule"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceNameCocoaPodsProxy, "replication.asset_path_regex"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func repositoryCocoaPodsProxyResourceMinimalConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-minimal-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://cdn.cocoapods.org/"
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
}
`, resourceTypeCocoaPodsProxy, randomString)
}

func repositoryCocoaPodsProxyResourceConfig(randomString string, includeReplication bool) string {
	var replicationConfig = ""
	if includeReplication {
		replicationConfig = `
	replication = {
		preemptive_pull_enabled = true
		asset_path_regex = "some-value"
	}	
`
	}
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://cdn.cocoapods.org/"
    content_max_age = 1441
    metadata_max_age = 1440
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
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
  %s
}
`, resourceTypeCocoaPodsProxy, randomString, replicationConfig)
}

func TestAccRepositoryCocoaPodsProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-url-without-protocol"
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
}
`, resourceTypeCocoaPodsProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageInvalidRemoteUrl),
			},
		},
	})
}

func TestAccRepositoryCocoapodsProxyInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://cdn.cocoapods.org/"
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
}
`, "sonatyperepo_repository_cocoapods_proxy", randomString),
				ExpectError: regexp.MustCompile(errorMessageBlobStoreNotFound),
			},
		},
	})
}

func TestAccRepositoryCocoapodsProxyMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-%s"
  online = true
  # Missing storage block
  proxy = {
    remote_url = "https://cdn.cocoapods.org/"
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
}
`, "sonatyperepo_repository_cocoapods_proxy", randomString),
				ExpectError: regexp.MustCompile(errorMessageStorageRequired),
			},
		},
	})
}

func TestAccRepositoryCocoapodsProxyInvalidTimeoutTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too large, max is 3600)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-timeout-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      timeout = 3601
    }
  }
}
`, "sonatyperepo_repository_cocoapods_proxy", randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryCocoapodsProxyInvalidTimeoutTooSmall(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too small, min is 1)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-timeout-small-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      timeout = 0
    }
  }
}
`, "sonatyperepo_repository_cocoapods_proxy", randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryCocoapodsProxyInvalidRetriesTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (too large, max is 10)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-retries-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      retries = 11
    }
  }
}
`, "sonatyperepo_repository_cocoapods_proxy", randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryCocoapodsProxyInvalidRetriesNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-retries-neg-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      retries = -1
    }
  }
}
`, "sonatyperepo_repository_cocoapods_proxy", randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryCocoapodsProxyInvalidTimeToLiveNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid time_to_live (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-ttl-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
    content_max_age = 1440
    metadata_max_age = 1440
  }
  negative_cache = {
    enabled = true
    time_to_live = -1
  }
  http_client = {
    blocked = false
    auto_block = true
  }
}
`, "sonatyperepo_repository_cocoapods_proxy", randomString),
				ExpectError: regexp.MustCompile(errorMessageNegativeCacheTimeoutValue),
			},
		},
	})
}
