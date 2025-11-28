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
	resourceTypeRawGroup  = "sonatyperepo_repository_raw_group"
	resourceTypeRawHosted = "sonatyperepo_repository_raw_hosted"
	resourceTypeRawProxy  = "sonatyperepo_repository_raw_proxy"
)

var (
	resourceRawGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRawGroup)
	resourceRawHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRawHosted)
	resourceRawProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRawProxy)
)

func TestAccRepositoryRawResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
  raw = {
	content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  raw = {
	content_disposition = "ATTACHMENT"
  }
}

resource "%s" "repo" {
  name = "raw-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://nodejs.org/dist/"
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
  raw = {
	content_disposition = "ATTACHMENT"
  }
}

resource "%s" "repo" {
  name = "raw-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["raw-proxy-repo-%s"]
  }
  raw = {
	content_disposition = "ATTACHMENT"
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeRawHosted, randomString, resourceTypeRawProxy, randomString, resourceTypeRawGroup, randomString, randomString, resourceTypeRawProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceRawHostedName, "name", fmt.Sprintf("raw-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRawHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceRawHostedName, "url"),
					resource.TestCheckResourceAttr(resourceRawHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRawHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceRawHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceRawHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceRawHostedName, "cleanup"),
					resource.TestCheckResourceAttr(resourceRawHostedName, RES_ATTR_RAW_CONTENT_DISPOSITION, common.CONTENT_DISPOSITION_ATTACHMENT),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceRawProxyName, "name", fmt.Sprintf("raw-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRawProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceRawProxyName, "url"),
					resource.TestCheckResourceAttr(resourceRawProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRawProxyName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "proxy.remote_url", "https://nodejs.org/dist/"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "proxy.content_max_age", "1442"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "proxy.metadata_max_age", "1400"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceRawProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceRawProxyName, RES_ATTR_RAW_CONTENT_DISPOSITION, common.CONTENT_DISPOSITION_ATTACHMENT),

					// Verify Group
					resource.TestCheckResourceAttr(resourceRawGroupName, "name", fmt.Sprintf("raw-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRawGroupName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceRawGroupName, "url"),
					resource.TestCheckResourceAttr(resourceRawGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRawGroupName, "group.member_names.#", "1"),
					resource.TestCheckResourceAttr(resourceRawGroupName, RES_ATTR_RAW_CONTENT_DISPOSITION, common.CONTENT_DISPOSITION_ATTACHMENT),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryRawHostedImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("raw-hosted-import-%s", randomString)
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
  raw = {
    content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawHosted, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRawHostedName, "name", repoName),
					resource.TestCheckResourceAttr(resourceRawHostedName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceRawHostedName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryRawProxyImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("raw-proxy-import-%s", randomString)
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
    remote_url = "https://nodejs.org/dist/"
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
  raw = {
    content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRawProxyName, "name", repoName),
					resource.TestCheckResourceAttr(resourceRawProxyName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceRawProxyName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryRawGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceType := "sonatyperepo_repository_raw_group"
	resourceTypeHosted := "sonatyperepo_repository_raw_hosted"
	resourceName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceType)
	repoName := fmt.Sprintf("raw-group-import-%s", randomString)
	memberName := fmt.Sprintf("raw-hosted-member-%s", randomString)

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
  raw = {
    content_disposition = "ATTACHMENT"
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
  raw = {
    content_disposition = "ATTACHMENT"
  }
  depends_on = [%s.member]
}
`, resourceTypeHosted, memberName, resourceType, repoName, memberName, resourceTypeHosted),
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
func TestAccRepositoryRawProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-proxy-repo-%s"
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
  raw = {
    content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageInvalidRemoteUrl),
			},
		},
	})
}

func TestAccRepositoryRawHostedInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-hosted-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
    write_policy = "ALLOW"
  }
  raw = {
    content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawHosted, randomString),
				ExpectError: regexp.MustCompile(errorMessageBlobStoreNotFound),
			},
		},
	})
}

func TestAccRepositoryRawHostedMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-hosted-repo-%s"
  online = true
  # Missing storage block
}
`, resourceTypeRawHosted, randomString),
				ExpectError: regexp.MustCompile(errorMessageStorageRequired),
			},
		},
	})
}

func TestAccRepositoryRawProxyInvalidTimeoutTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too large, max is 3600)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-proxy-repo-timeout-%s"
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
  raw = {
    content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryRawProxyInvalidTimeoutTooSmall(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too small, min is 1)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-proxy-repo-timeout-small-%s"
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
  raw = {
    content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryRawProxyInvalidRetriesTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (too large, max is 10)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-proxy-repo-retries-%s"
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
  raw = {
    content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryRawProxyInvalidRetriesNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-proxy-repo-retries-neg-%s"
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
  raw = {
    content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryRawProxyInvalidTimeToLiveNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid time_to_live (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-proxy-repo-ttl-%s"
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
  raw = {
    content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageNegativeCacheTimeoutValue),
			},
		},
	})
}
