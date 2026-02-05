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
	resourceTypeConanGroup  = "sonatyperepo_repository_conan_group"
	resourceTypeConanHosted = "sonatyperepo_repository_conan_hosted"
	resourceTypeConanProxy  = "sonatyperepo_repository_conan_proxy"
)

var (
	resourceConanGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeConanGroup)
	resourceConanHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeConanHosted)
	resourceConanProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeConanProxy)
)

func TestAccRepositoryConanResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "conan-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
}
`, resourceTypeConanGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "conan-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}

resource "%s" "repo" {
  name = "conan-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://center2.conan.io"
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
  conan = {
    conan_version = "V2"
  }
}

resource "%s" "repo" {
  name = "conan-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["conan-proxy-repo-%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeConanHosted, randomString, resourceTypeConanProxy, randomString, resourceTypeConanGroup, randomString, randomString, resourceTypeConanProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceConanHostedName, RES_ATTR_NAME, fmt.Sprintf("conan-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceConanHostedName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceConanHostedName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceConanHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceConanHostedName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceConanHostedName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceConanHostedName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
					resource.TestCheckNoResourceAttr(resourceConanHostedName, RES_ATTR_CLEANUP),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceConanProxyName, RES_ATTR_NAME, fmt.Sprintf("conan-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceConanProxyName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceConanProxyName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceConanProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceConanProxyName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "proxy.remote_url", "https://center2.conan.io"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceConanProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceConanProxyName, "replication.asset_path_regex"),
					resource.TestCheckResourceAttr(resourceConanProxyName, "conan.conan_version", "V2"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceConanGroupName, RES_ATTR_NAME, fmt.Sprintf("conan-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceConanGroupName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceConanGroupName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceConanGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceConanGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
