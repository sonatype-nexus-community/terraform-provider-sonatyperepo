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
	resourceTypeAptProxy = "sonatyperepo_repository_apt_proxy"
)

var (
	resourceAptProxyName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeAptProxy)
)

func TestAccRepositoryAptProxyResourceNoReplication(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: getRepositoryAptProxyResourceMinimalConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify minimal config
					resource.TestCheckResourceAttr(resourceAptProxyName, "name", fmt.Sprintf("apt-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceAptProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceAptProxyName, "url"),
					resource.TestCheckResourceAttr(resourceAptProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceAptProxyName, "proxy.remote_url", "https://archive.ubuntu.com/ubuntu/"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "apt.distribution", "bionic"),
				),
			},
			// Update to full configuration
			{
				Config: getRepositoryAptProxyResourceConfig(randomString, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify full config
					resource.TestCheckResourceAttr(resourceAptProxyName, "name", fmt.Sprintf("apt-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceAptProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceAptProxyName, "url"),
					resource.TestCheckResourceAttr(resourceAptProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceAptProxyName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "proxy.remote_url", "https://archive.ubuntu.com/ubuntu/"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "proxy.content_max_age", "1442"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "proxy.metadata_max_age", "1400"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceAptProxyName, "http_client.connection.timeout", "999"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryAptProxyResourceMinimalConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "apt-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://archive.ubuntu.com/ubuntu/"
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
  apt = {
	distribution = "bionic"
  }
}
`, resourceTypeAptProxy, randomString)
}

func getRepositoryAptProxyResourceConfig(randomString string, includeReplication bool) string {
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
  name = "apt-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://archive.ubuntu.com/ubuntu/"
    content_max_age = 1442
    metadata_max_age = 1400
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
  apt = {
	distribution = "bionic"
  }
  %s
}
`, resourceTypeAptProxy, randomString, replicationConfig)
}

func TestAccRepositoryAptProxyImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("apt-proxy-import-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
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
    remote_url = "https://archive.ubuntu.com/ubuntu/"
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
  apt = {
    distribution = "bionic"
  }
}
`, resourceTypeAptProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAptProxyName, "name", repoName),
					resource.TestCheckResourceAttr(resourceAptProxyName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceAptProxyName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryAptProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "apt-proxy-repo-%s"
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
  apt = {
    distribution = "bionic"
  }
}
`, resourceTypeAptProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageInvalidRemoteUrl),
			},
		},
	})
}

func TestAccRepositoryAptProxyInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "apt-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://archive.ubuntu.com/ubuntu/"
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
  apt = {
    distribution = "bionic"
  }
}
`, resourceTypeAptProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageBlobStoreNotFound),
			},
		},
	})
}

func TestAccRepositoryAptProxyMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "apt-proxy-repo-%s"
  online = true
  # Missing storage block
  proxy = {
    remote_url = "https://archive.ubuntu.com/ubuntu/"
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
  apt = {
    distribution = "bionic"
  }
}
`, resourceTypeAptProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageStorageRequired),
			},
		},
	})
}

func TestAccRepositoryAptProxyInvalidTimeoutTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too large, max is 3600)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "apt-proxy-repo-timeout-%s"
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
  apt = {
    distribution = "bionic"
  }
}
`, resourceTypeAptProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryAptProxyInvalidTimeoutTooSmall(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too small, min is 1)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "apt-proxy-repo-timeout-small-%s"
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
  apt = {
    distribution = "bionic"
  }
}
`, resourceTypeAptProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryAptProxyInvalidRetriesTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (too large, max is 10)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "apt-proxy-repo-retries-%s"
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
  apt = {
    distribution = "bionic"
  }
}
`, resourceTypeAptProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryAptProxyInvalidRetriesNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "apt-proxy-repo-retries-neg-%s"
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
  apt = {
    distribution = "bionic"
  }
}
`, resourceTypeAptProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryAptProxyInvalidTimeToLiveNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid time_to_live (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "apt-proxy-repo-ttl-%s"
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
  apt = {
    distribution = "bionic"
  }
}
`, resourceTypeAptProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageNegativeCacheTimeoutValue),
			},
		},
	})
}
