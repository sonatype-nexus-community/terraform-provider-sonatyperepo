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
					resource.TestCheckResourceAttr(resourceNugetHostedName, RES_ATTR_NAME, fmt.Sprintf("nuget-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNugetHostedName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceNugetHostedName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceNugetHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceNugetHostedName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceNugetHostedName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceNugetHostedName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
					resource.TestCheckNoResourceAttr(resourceNugetHostedName, RES_ATTR_CLEANUP),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceNugetProxyName, RES_ATTR_NAME, fmt.Sprintf("nuget-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNugetProxyName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceNugetProxyName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceNugetProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceNugetProxyName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
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
					resource.TestCheckResourceAttr(resourceNugetGroupName, RES_ATTR_NAME, fmt.Sprintf("nuget-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNugetGroupName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceNugetGroupName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceNugetGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceNugetGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
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
					resource.TestCheckResourceAttr(resourceNugetGroupName, RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceNugetGroupName, RES_ATTR_ONLINE, "true"),
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
