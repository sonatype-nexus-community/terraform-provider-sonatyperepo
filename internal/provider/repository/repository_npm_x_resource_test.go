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

func TestAccRepositoryNpmResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceTypeGroup := "sonatyperepo_repository_npm_group"
	resourceTypeHosted := "sonatyperepo_repository_npm_hosted"
	resourceTypeProxy := "sonatyperepo_repository_npm_proxy"
	resourceGroupName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeGroup)
	resourceHostedName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeHosted)
	resourceProxyName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeProxy)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "npm-group-repo-%s"
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
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "npm-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}

resource "%s" "repo" {
  name = "npm-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://registry.npmjs.org"
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
  npm = {
	remove_quarrantined = true
  }
}

resource "%s" "repo" {
  name = "npm-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["npm-proxy-repo-%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeHosted, randomString, resourceTypeProxy, randomString, resourceTypeGroup, randomString, randomString, resourceTypeProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceHostedName, "name", fmt.Sprintf("npm-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceHostedName, "url"),
					resource.TestCheckResourceAttr(resourceHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceHostedName, "cleanup"),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceProxyName, "name", fmt.Sprintf("npm-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceProxyName, "url"),
					resource.TestCheckResourceAttr(resourceProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceProxyName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "proxy.remote_url", "https://registry.npmjs.org"),
					resource.TestCheckResourceAttr(resourceProxyName, "proxy.content_max_age", "1442"),
					resource.TestCheckResourceAttr(resourceProxyName, "proxy.metadata_max_age", "1400"),
					resource.TestCheckResourceAttr(resourceProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckResourceAttr(resourceProxyName, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceProxyName, "replication.asset_path_regex"),
					resource.TestCheckNoResourceAttr(resourceProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceProxyName, "npm.remove_quarrantined", "true"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceGroupName, "name", fmt.Sprintf("npm-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceGroupName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceGroupName, "url"),
					resource.TestCheckResourceAttr(resourceGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryNpmHostedImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceType := "sonatyperepo_repository_npm_hosted"
	resourceName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceType)
	repoName := fmt.Sprintf("npm-hosted-import-%s", randomString)

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
`, resourceType, repoName),
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

func TestAccRepositoryNpmProxyImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceType := "sonatyperepo_repository_npm_proxy"
	resourceName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceType)
	repoName := fmt.Sprintf("npm-proxy-import-%s", randomString)

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
    remote_url = "https://registry.npmjs.org"
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
`, resourceType, repoName),
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

func TestAccRepositoryNpmGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceType := "sonatyperepo_repository_npm_group"
	resourceTypeHosted := "sonatyperepo_repository_npm_hosted"
	resourceName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceType)
	repoName := fmt.Sprintf("npm-group-import-%s", randomString)
	memberName := fmt.Sprintf("npm-hosted-member-%s", randomString)

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
func TestAccRepositoryNpmProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "npm-proxy-repo-%s"
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
`, "sonatyperepo_repository_npm_proxy", randomString),
				ExpectError: regexp.MustCompile("must be a valid URL|must be a valid HTTP URL"),
			},
		},
	})
}

func TestAccRepositoryNpmHostedInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "npm-hosted-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  npm = {
    package_name_normalize = "lowercase"
  }
}
`, "sonatyperepo_repository_npm_hosted", randomString),
				ExpectError: regexp.MustCompile("Blob store.*not found|Blob store.*does not exist"),
			},
		},
	})
}

func TestAccRepositoryNpmHostedMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "npm-hosted-repo-%s"
  online = true
  # Missing storage block
}
`, "sonatyperepo_repository_npm_hosted", randomString),
				ExpectError: regexp.MustCompile("Attribute storage is required"),
			},
		},
	})
}



func TestAccRepositoryNpmProxyInvalidTimeoutTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too large, max is 3600)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "npm-proxy-repo-timeout-%s"
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
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be less than or equal to 3600"),
			},
		},
	})
}

func TestAccRepositoryNpmProxyInvalidTimeoutTooSmall(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too small, min is 1)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "npm-proxy-repo-timeout-small-%s"
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
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be greater than or equal to 1"),
			},
		},
	})
}

func TestAccRepositoryNpmProxyInvalidRetriesTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (too large, max is 10)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "npm-proxy-repo-retries-%s"
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
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be less than or equal to 10"),
			},
		},
	})
}

func TestAccRepositoryNpmProxyInvalidRetriesNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "npm-proxy-repo-retries-neg-%s"
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
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be greater than or equal to 0"),
			},
		},
	})
}

func TestAccRepositoryNpmProxyInvalidMaxAgeNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid content_max_age (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "npm-proxy-repo-maxage-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
    content_max_age = -1
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
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be greater than or equal to|cannot be negative"),
			},
		},
	})
}

func TestAccRepositoryNpmProxyInvalidTimeToLiveNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid time_to_live (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "npm-proxy-repo-ttl-%s"
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
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be greater than or equal to|cannot be negative"),
			},
		},
	})
}

