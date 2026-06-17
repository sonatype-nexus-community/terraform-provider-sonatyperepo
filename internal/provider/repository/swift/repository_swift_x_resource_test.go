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

package swift_test

import (
	"fmt"
	"regexp"
	"terraform-provider-sonatyperepo/internal/provider/common"
	repotest "terraform-provider-sonatyperepo/internal/provider/repository/repotest"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeSwiftGroup = "sonatyperepo_repository_swift_group"
	resourceTypeSwiftProxy = "sonatyperepo_repository_swift_proxy"
)

var (
	resourceSwiftGroupName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeSwiftGroup)
	resourceSwiftProxyName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeSwiftProxy)
)

func TestAccRepositorySwiftResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Requires NXRM 3.91.0+
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 90,
				Patch: 99,
			})
		},
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "swift-group-repo-%s"
  online = true
  storage = {
	  blob_store_name = "default"
	  strict_content_type_validation = true
  }
  group = {
	  member_names = []
  }
}
`, resourceTypeSwiftGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "swift-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://github.com/"
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
  swift = {
    require_authentication = true
  }
}

resource "%s" "repo" {
  name = "swift-group-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
	  member_names = ["swift-proxy-repo-%s"]
  }

  depends_on = [
	  %s.repo
  ]
}
`, resourceTypeSwiftProxy, randomString, resourceTypeSwiftGroup, randomString, randomString, resourceTypeSwiftProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Proxy
					resource.TestCheckResourceAttr(resourceSwiftProxyName, repotest.RES_ATTR_NAME, fmt.Sprintf("swift-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceSwiftProxyName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "proxy.remote_url", "https://github.com/"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceSwiftProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceSwiftProxyName, "replication.asset_path_regex"),
					resource.TestCheckResourceAttr(resourceSwiftProxyName, "swift.require_authentication", "true"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceSwiftGroupName, repotest.RES_ATTR_NAME, fmt.Sprintf("swift-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceSwiftGroupName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceSwiftGroupName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceSwiftGroupName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceSwiftGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositorySwiftGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("swift-group-import-%s", randomString)
	memberName := fmt.Sprintf("swift-proxy-member-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Requires NXRM 3.91.0+
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 90,
				Patch: 99,
			})
		},
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
  }
  proxy = {
    remote_url = "https://github.com/"
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
  swift = {
    require_authentication = false
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
`, resourceTypeSwiftProxy, memberName, resourceTypeSwiftGroup, repoName, memberName, resourceTypeSwiftProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceSwiftGroupName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceSwiftGroupName, repotest.RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceSwiftGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
