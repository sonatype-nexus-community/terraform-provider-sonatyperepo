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
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeCargoGroup  = "sonatyperepo_repository_cargo_group"
	resourceTypeCargoHosted = "sonatyperepo_repository_cargo_hosted"
	resourceTypeCargoProxy  = "sonatyperepo_repository_cargo_proxy"
)

var (
	resourceCargoGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeCargoGroup)
	resourceCargoHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeCargoHosted)
	resourceCargoProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeCargoProxy)
)

func TestAccRepositoryCargoResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Know regression in NXRM 3.82.0 - skip these tests as they will fail - see https://sonatype.atlassian.net/browse/NEXUS-48088 - fix coming 3.88.x
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 82,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 87,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cargo-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
  cargo = {
    require_authentication = false
  }
}
`, resourceTypeCargoGroup, randomString),
				ExpectError: regexp.MustCompile(errorMessageGroupMemberNamesEmpty),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cargo-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}

resource "%s" "repo" {
  name = "cargo-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://index.crates.io/"
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
  cargo = {
    require_authentication = true
  }
}

resource "%s" "repo" {
  name = "cargo-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["cargo-proxy-repo-%s"]
  }
  cargo = {
    require_authentication = false
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeCargoHosted, randomString, resourceTypeCargoProxy, randomString, resourceTypeCargoGroup, randomString, randomString, resourceTypeCargoProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceCargoHostedName, RES_ATTR_NAME, fmt.Sprintf("cargo-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceCargoHostedName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceCargoHostedName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceCargoHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceCargoHostedName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceCargoHostedName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceCargoHostedName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
					resource.TestCheckNoResourceAttr(resourceCargoHostedName, RES_ATTR_CLEANUP),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceCargoProxyName, RES_ATTR_NAME, fmt.Sprintf("cargo-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceCargoProxyName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceCargoProxyName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceCargoProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, "default"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "proxy.remote_url", "https://index.crates.io/"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceCargoProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceCargoProxyName, "replication.asset_path_regex"),
					resource.TestCheckResourceAttr(resourceCargoProxyName, "cargo.require_authentication", "true"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceCargoGroupName, RES_ATTR_NAME, fmt.Sprintf("cargo-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceCargoGroupName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceCargoGroupName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceCargoGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, "default"),
					resource.TestCheckResourceAttr(resourceCargoGroupName, "group.member_names.#", "1"),
					resource.TestCheckResourceAttr(resourceCargoGroupName, "cargo.require_authentication", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
