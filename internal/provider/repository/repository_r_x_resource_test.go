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
					resource.TestCheckResourceAttr(resourceRHostedName, RES_ATTR_NAME, fmt.Sprintf("r-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRHostedName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceRHostedName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceRHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRHostedName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceRHostedName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceRHostedName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
					resource.TestCheckNoResourceAttr(resourceRHostedName, RES_ATTR_CLEANUP),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceRProxyName, RES_ATTR_NAME, fmt.Sprintf("r-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRProxyName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceRProxyName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceRProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRProxyName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
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
					resource.TestCheckResourceAttr(resourceRGroupName, RES_ATTR_NAME, fmt.Sprintf("r-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRGroupName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceRGroupName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceRGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
