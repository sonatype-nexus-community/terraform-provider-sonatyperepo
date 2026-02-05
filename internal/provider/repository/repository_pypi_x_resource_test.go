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
					resource.TestCheckResourceAttr(resourcePypiHostedName, RES_ATTR_NAME, fmt.Sprintf("pypi-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourcePypiHostedName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourcePypiHostedName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourcePypiHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourcePypiHostedName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourcePypiHostedName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourcePypiHostedName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
					resource.TestCheckNoResourceAttr(resourcePypiHostedName, RES_ATTR_CLEANUP),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourcePypiProxyName, RES_ATTR_NAME, fmt.Sprintf("pypi-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourcePypiProxyName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourcePypiProxyName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourcePypiProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourcePypiProxyName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
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
					resource.TestCheckResourceAttr(resourcePypiProxyName, "repository_firewall.enabled", "false"),
					// resource.TestCheckResourceAttr(resourcePypiProxyName, "repository_firewall.quarantine", "true"),
					// resource.TestCheckResourceAttr(resourcePypiProxyName, "repository_firewall.pccs_enabled", "true"),

					// Verify Group
					resource.TestCheckResourceAttr(resourcePypiGroupName, RES_ATTR_NAME, fmt.Sprintf("pypi-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourcePypiGroupName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourcePypiGroupName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourcePypiGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourcePypiGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
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
					resource.TestCheckResourceAttr(resourcePypiGroupName, RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourcePypiGroupName, RES_ATTR_ONLINE, "true"),
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
