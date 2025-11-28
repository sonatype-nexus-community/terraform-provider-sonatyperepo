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
	resourceTypePypiGroup  = "sonatyperepo_repository_pypi_group"
	resourceTypePypiHosted = "sonatyperepo_repository_pypi_hosted"
	resourceTypePypiProxy  = "sonatyperepo_repository_pypi_proxy"
)

var (
	resourcePypiGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypePypiGroup)
	resourcePypiHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypePypiHosted)
	resourcePypiProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypePypiProxy)
)

func TestAccRepositoryPyPiResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pypi-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
}
`, resourceTypePypiGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pypi-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}

resource "%s" "repo" {
  name = "pypi-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://pypi.org/"
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
  pypi = {
    remove_quarrantined = true
  }
}

resource "%s" "repo" {
  name = "pypi-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["pypi-proxy-repo-%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypePypiHosted, randomString, resourceTypePypiProxy, randomString, resourceTypePypiGroup, randomString, randomString, resourceTypePypiProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourcePypiHostedName, "name", fmt.Sprintf("pypi-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourcePypiHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourcePypiHostedName, "url"),
					resource.TestCheckResourceAttr(resourcePypiHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourcePypiHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourcePypiHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourcePypiHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourcePypiHostedName, "cleanup"),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourcePypiProxyName, "name", fmt.Sprintf("pypi-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourcePypiProxyName, "url"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "proxy.remote_url", "https://pypi.org/"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "proxy.content_max_age", "1442"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "proxy.metadata_max_age", "1400"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourcePypiProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "pypi.remove_quarrantined", "true"),

					// Verify Group
					resource.TestCheckResourceAttr(resourcePypiGroupName, "name", fmt.Sprintf("pypi-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourcePypiGroupName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourcePypiGroupName, "url"),
					resource.TestCheckResourceAttr(resourcePypiGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourcePypiGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryPyPiHostedImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("pypi-hosted-import-%s", randomString)

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
    write_policy = "ALLOW_ONCE"
  }
}
`, resourceTypePypiHosted, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePypiHostedName, "name", repoName),
					resource.TestCheckResourceAttr(resourcePypiHostedName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourcePypiHostedName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryPyPiProxyImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("pypi-proxy-import-%s", randomString)

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
    remote_url = "https://pypi.org/"
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
  pypi = {
    remove_quarrantined = true
  }
}
`, resourceTypePypiProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePypiProxyName, "name", repoName),
					resource.TestCheckResourceAttr(resourcePypiProxyName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourcePypiProxyName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryPyPiGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("pypi-group-import-%s", randomString)
	memberName := fmt.Sprintf("pypi-hosted-member-%s", randomString)

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
    write_policy = "ALLOW"
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
`, resourceTypePypiHosted, memberName, resourceTypePypiGroup, repoName, memberName, resourceTypePypiHosted),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePypiGroupName, "name", repoName),
					resource.TestCheckResourceAttr(resourcePypiGroupName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourcePypiGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
func TestAccRepositoryPypiProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pypi-proxy-repo-%s"
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
  pypi = {
    index_type = "SIMPLE"
  }
}
`, resourceTypePypiProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageInvalidRemoteUrl),
			},
		},
	})
}

func TestAccRepositoryPypiHostedInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pypi-hosted-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  pypi = {}
}
`, resourceTypePypiHosted, randomString),
				ExpectError: regexp.MustCompile(errorMessageBlobStoreNotFound),
			},
		},
	})
}

func TestAccRepositoryPypiHostedMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pypi-hosted-repo-%s"
  online = true
  # Missing storage block
}
`, resourceTypePypiHosted, randomString),
				ExpectError: regexp.MustCompile(errorMessageStorageRequired),
			},
		},
	})
}

func TestAccRepositoryPypiProxyInvalidTimeoutTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too large, max is 3600)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pypi-proxy-repo-timeout-%s"
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
`, resourceTypePypiProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryPypiProxyInvalidTimeoutTooSmall(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too small, min is 1)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pypi-proxy-repo-timeout-small-%s"
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
`, resourceTypePypiProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryPypiProxyInvalidRetriesTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (too large, max is 10)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pypi-proxy-repo-retries-%s"
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
`, resourceTypePypiProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryPypiProxyInvalidRetriesNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pypi-proxy-repo-retries-neg-%s"
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
`, resourceTypePypiProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryPypiProxyInvalidTimeToLiveNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid time_to_live (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pypi-proxy-repo-ttl-%s"
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
`, resourceTypePypiProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageNegativeCacheTimeoutValue),
			},
		},
	})
}
