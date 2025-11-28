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
	resourceTypeGoGroup = "sonatyperepo_repository_go_group"
	resourceTypeGoProxy = "sonatyperepo_repository_go_proxy"
)

var (
	resourceGoGroupName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeGoGroup)
	resourceGoProxyName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeGoProxy)
)

func TestAccRepositoryGoResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "go-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
}
`, resourceTypeGoGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "go-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://proxy.golang.org/"
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
}

resource "%s" "repo" {
  name = "go-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["go-proxy-repo-%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeGoProxy, randomString, resourceTypeGoGroup, randomString, randomString, resourceTypeGoProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Proxy
					resource.TestCheckResourceAttr(resourceGoProxyName, "name", fmt.Sprintf("go-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceGoProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceGoProxyName, "url"),
					resource.TestCheckResourceAttr(resourceGoProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceGoProxyName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "proxy.remote_url", "https://proxy.golang.org/"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceGoProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceGoProxyName, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceGoProxyName, "replication.asset_path_regex"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceGoGroupName, "name", fmt.Sprintf("go-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceGoGroupName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceGoGroupName, "url"),
					resource.TestCheckResourceAttr(resourceGoGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceGoGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryGoProxyImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeGoProxy)
	repoName := fmt.Sprintf("go-proxy-import-%s", randomString)

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
    remote_url = "https://proxy.golang.org/"
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
`, resourceTypeGoProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", repoName),
					resource.TestCheckResourceAttr(resourceName, "online", "true"),
				),
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

func TestAccRepositoryGoGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("go-group-import-%s", randomString)
	memberName := fmt.Sprintf("go-proxy-member-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "member" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://proxy.golang.org/"
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

resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = ["%s"]
  }
  depends_on = [%s.member]
}
`, resourceTypeGoProxy, memberName, resourceTypeGoGroup, repoName, memberName, resourceTypeGoProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceGoGroupName, "name", repoName),
					resource.TestCheckResourceAttr(resourceGoGroupName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceGoGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
func TestAccRepositoryGoProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "go-proxy-repo-%s"
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
`, resourceTypeGoProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageInvalidRemoteUrl),
			},
		},
	})
}

func TestAccRepositoryGoProxyInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "go-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  go = {}
}
`, resourceTypeGoProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageBlobStoreNotFound),
			},
		},
	})
}

func TestAccRepositoryGoProxyMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "go-proxy-repo-%s"
  online = true
  # Missing storage block
}
`, resourceTypeGoGroup, randomString),
				ExpectError: regexp.MustCompile(errorMessageStorageRequired),
			},
		},
	})
}

func TestAccRepositoryGoProxyInvalidTimeoutTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too large, max is 3600)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "go-proxy-repo-timeout-%s"
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
`, resourceTypeGoProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryGoProxyInvalidTimeoutTooSmall(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too small, min is 1)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "go-proxy-repo-timeout-small-%s"
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
`, resourceTypeGoProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryGoProxyInvalidRetriesTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (too large, max is 10)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {s
  name = "go-proxy-repo-retries-%s"
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
`, resourceTypeGoProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryGoProxyInvalidRetriesNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "go-proxy-repo-retries-neg-%s"
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
`, resourceTypeGoProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryGoProxyInvalidTimeToLiveNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid time_to_live (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "go-proxy-repo-ttl-%s"
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
`, resourceTypeGoProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageNegativeCacheTimeoutValue),
			},
		},
	})
}
