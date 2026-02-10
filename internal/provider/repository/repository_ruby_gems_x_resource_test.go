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
	resourceTypeRubyGemsGroup  = "sonatyperepo_repository_rubygems_group"
	resourceTypeRubyGemsHosted = "sonatyperepo_repository_rubygems_hosted"
	resourceTypeRubyGemsProxy  = "sonatyperepo_repository_rubygems_proxy"
)

var (
	resourceRubyGemsGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRubyGemsGroup)
	resourceRubyGemsHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRubyGemsHosted)
	resourceRubyGemsProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRubyGemsProxy)
)

func TestAccRepositoryRubyGemsResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "ruby-gems-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
}
`, resourceTypeRubyGemsGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "ruby-gems-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}

resource "%s" "repo" {
  name = "ruby-gems-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://rubygems.org"
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
  name = "ruby-gems-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["ruby-gems-proxy-repo-%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeRubyGemsHosted, randomString, resourceTypeRubyGemsProxy, randomString, resourceTypeRubyGemsGroup, randomString, randomString, resourceTypeRubyGemsProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceRubyGemsHostedName, RES_ATTR_NAME, fmt.Sprintf("ruby-gems-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRubyGemsHostedName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceRubyGemsHostedName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceRubyGemsHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRubyGemsHostedName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceRubyGemsHostedName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceRubyGemsHostedName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
					resource.TestCheckNoResourceAttr(resourceRubyGemsHostedName, RES_ATTR_CLEANUP),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, RES_ATTR_NAME, fmt.Sprintf("ruby-gems-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceRubyGemsProxyName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "proxy.remote_url", "https://rubygems.org"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceRubyGemsProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceRubyGemsProxyName, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceRubyGemsProxyName, "replication.asset_path_regex"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceRubyGemsGroupName, RES_ATTR_NAME, fmt.Sprintf("ruby-gems-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRubyGemsGroupName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceRubyGemsGroupName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceRubyGemsGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRubyGemsGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryRubyGemsGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("ruby-gems-group-import-%s", randomString)
	memberName := fmt.Sprintf("ruby-gems-hosted-member-%s", randomString)
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
`, resourceTypeRubyGemsHosted, memberName, resourceTypeRubyGemsGroup, repoName, memberName, resourceTypeRubyGemsHosted),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRubyGemsGroupName, RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceRubyGemsGroupName, RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceRubyGemsGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
