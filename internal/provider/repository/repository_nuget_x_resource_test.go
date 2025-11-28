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
	resourceTypeNugetGroup  = "sonatyperepo_repository_nuget_group"
	resourceTypeNugetHosted = "sonatyperepo_repository_nuget_hosted"
	resourceTypeNugetProxy  = "sonatyperepo_repository_nuget_proxy"
)

var (
	resourceNugetGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeNugetGroup)
	resourceNugetHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeNugetHosted)
	resourceNugetProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeNugetProxy)
)

func TestAccRepositoryNugetResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "nuget-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
}
`, resourceTypeNugetGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "nuget-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}

resource "%s" "repo" {
  name = "nuget-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://api.nuget.org/v3/index.json"
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
  nuget_proxy = {
    nuget_version = "V2"
  }
}

resource "%s" "repo" {
  name = "nuget-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["nuget-proxy-repo-%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeNugetHosted, randomString, resourceTypeNugetProxy, randomString, resourceTypeNugetGroup, randomString, randomString, resourceTypeNugetProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceNugetHostedName, "name", fmt.Sprintf("nuget-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNugetHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNugetHostedName, "url"),
					resource.TestCheckResourceAttr(resourceNugetHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceNugetHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNugetHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceNugetHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceNugetHostedName, "cleanup"),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceNugetProxyName, "name", fmt.Sprintf("nuget-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNugetProxyName, "url"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "proxy.remote_url", "https://api.nuget.org/v3/index.json"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "proxy.content_max_age", "1442"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "proxy.metadata_max_age", "1400"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceNugetProxyName, "routing_rule"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceNugetGroupName, "name", fmt.Sprintf("nuget-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNugetGroupName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNugetGroupName, "url"),
					resource.TestCheckResourceAttr(resourceNugetGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceNugetGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryNugetHostedImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("nuget-hosted-import-%s", randomString)

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
`, resourceTypeNugetHosted, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNugetHostedName, "name", repoName),
					resource.TestCheckResourceAttr(resourceNugetHostedName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceNugetHostedName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryNugetProxyImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("nuget-proxy-import-%s", randomString)

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
    remote_url = "https://api.nuget.org/v3/index.json"
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
  nuget_proxy = {
    nuget_version = "V2"
  }
}
`, resourceTypeNugetProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNugetProxyName, "name", repoName),
					resource.TestCheckResourceAttr(resourceNugetProxyName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceNugetProxyName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryNugetGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("nuget-group-import-%s", randomString)
	memberName := fmt.Sprintf("nuget-hosted-member-%s", randomString)

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
`, resourceTypeNugetHosted, memberName, resourceTypeNugetGroup, repoName, memberName, resourceTypeNugetHosted),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNugetGroupName, "name", repoName),
					resource.TestCheckResourceAttr(resourceNugetGroupName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceNugetGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
func TestAccRepositoryNugetProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "nuget-proxy-repo-%s"
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
`, resourceTypeNugetProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageInvalidRemoteUrl),
			},
		},
	})
}

func TestAccRepositoryNugetHostedInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "nuget-hosted-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  nuget = {}
}
`, resourceTypeNugetHosted, randomString),
				ExpectError: regexp.MustCompile(errorMessageBlobStoreNotFound),
			},
		},
	})
}

func TestAccRepositoryNugetHostedMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "nuget-hosted-repo-%s"
  online = true
  # Missing storage block
}
`, resourceTypeNugetHosted, randomString),
				ExpectError: regexp.MustCompile(errorMessageStorageRequired),
			},
		},
	})
}

func TestAccRepositoryNugetProxyInvalidTimeoutTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too large, max is 3600)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "nuget-proxy-repo-timeout-%s"
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
`, resourceTypeNugetProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryNugetProxyInvalidTimeoutTooSmall(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too small, min is 1)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "nuget-proxy-repo-timeout-small-%s"
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
`, resourceTypeNugetProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryNugetProxyInvalidRetriesTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (too large, max is 10)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "nuget-proxy-repo-retries-%s"
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
`, resourceTypeNugetProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryNugetProxyInvalidRetriesNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "nuget-proxy-repo-retries-neg-%s"
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
`, resourceTypeNugetProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryNugetProxyInvalidTimeToLiveNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid time_to_live (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "nuget-proxy-repo-ttl-%s"
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
`, resourceTypeNugetProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageNegativeCacheTimeoutValue),
			},
		},
	})
}
