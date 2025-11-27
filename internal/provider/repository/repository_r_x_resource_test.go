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
	resourceTypeRGroup  = "sonatyperepo_repository_r_group"
	resourceTypeRHosted = "sonatyperepo_repository_r_hosted"
	resourceTypeRProxy  = "sonatyperepo_repository_r_proxy"
)

var (
	resourceRGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRGroup)
	resourceRHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRHosted)
	resourceRProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRProxy)
)

func TestAccRepositoryRResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "r-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
}
`, resourceTypeRGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "r-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}

resource "%s" "repo" {
  name = "r-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://cran.r-project.org/"
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
  name = "r-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["r-proxy-repo-%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeRHosted, randomString, resourceTypeRProxy, randomString, resourceTypeRGroup, randomString, randomString, resourceTypeRProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceRHostedName, "name", fmt.Sprintf("r-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceRHostedName, "url"),
					resource.TestCheckResourceAttr(resourceRHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceRHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceRHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceRHostedName, "cleanup"),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceRProxyName, "name", fmt.Sprintf("r-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceRProxyName, "url"),
					resource.TestCheckResourceAttr(resourceRProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRProxyName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceRProxyName, "proxy.remote_url", "https://cran.r-project.org/"),
					resource.TestCheckResourceAttr(resourceRProxyName, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceRProxyName, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceRProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceRProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceRProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceRProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceRProxyName, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceRProxyName, "replication.asset_path_regex"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceRGroupName, "name", fmt.Sprintf("r-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRGroupName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceRGroupName, "url"),
					resource.TestCheckResourceAttr(resourceRGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccRepositoryRProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "r-proxy-repo-%s"
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
`, resourceTypeRProxy, randomString),
				ExpectError: regexp.MustCompile("must be a valid URL|must be a valid HTTP URL"),
			},
		},
	})
}

func TestAccRepositoryRProxyInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "r-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  r = {}
}
`, resourceTypeRProxy, randomString),
				ExpectError: regexp.MustCompile("Blob store.*not found|Blob store.*does not exist"),
			},
		},
	})
}

func TestAccRepositoryRHostedMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "r-hosted-repo-%s"
  online = true
  # Missing storage block
}
`, resourceTypeRHosted, randomString),
				ExpectError: regexp.MustCompile("Attribute storage is required"),
			},
		},
	})
}

func TestAccRepositoryRProxyInvalidTimeoutTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too large, max is 3600)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "r-proxy-repo-timeout-%s"
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
`, resourceTypeRProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be less than or equal to 3600"),
			},
		},
	})
}

func TestAccRepositoryRProxyInvalidTimeoutTooSmall(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too small, min is 1)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "r-proxy-repo-timeout-small-%s"
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
`, resourceTypeRProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be greater than or equal to 1"),
			},
		},
	})
}

func TestAccRepositoryRProxyInvalidRetriesTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (too large, max is 10)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "r-proxy-repo-retries-%s"
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
`, resourceTypeRProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be less than or equal to 10"),
			},
		},
	})
}

func TestAccRepositoryRProxyInvalidRetriesNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "r-proxy-repo-retries-neg-%s"
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
`, resourceTypeRProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be greater than or equal to 0"),
			},
		},
	})
}

func TestAccRepositoryRProxyInvalidMaxAgeNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid content_max_age (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "r-proxy-repo-maxage-%s"
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
`, resourceTypeRProxy, randomString),
				ExpectError: regexp.MustCompile("must be greater than or equal to|cannot be negative"),
			},
		},
	})
}

func TestAccRepositoryRProxyInvalidTimeToLiveNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid time_to_live (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "r-proxy-repo-ttl-%s"
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
`, resourceTypeRProxy, randomString),
				ExpectError: regexp.MustCompile("must be greater than or equal to|cannot be negative"),
			},
		},
	})
}
