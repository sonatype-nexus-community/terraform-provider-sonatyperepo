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
	resourceTypeYumGroup  = "sonatyperepo_repository_yum_group"
	resourceTypeYumHosted = "sonatyperepo_repository_yum_hosted"
	resourceTypeYumProxy  = "sonatyperepo_repository_yum_proxy"
)

var (
	resourceYumGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeYumGroup)
	resourceYumHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeYumHosted)
	resourceYumProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeYumProxy)
)

func TestAccRepositoryYumResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "yum-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
}
`, resourceTypeYumGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "yum-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  yum = {
    repo_data_depth = 0 
  }
}

resource "%s" "repo" {
  name = "yum-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://mirror.centos.org/centos/"
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
  name = "yum-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["yum-proxy-repo-%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeYumHosted, randomString, resourceTypeYumProxy, randomString, resourceTypeYumGroup, randomString, randomString, resourceTypeYumProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceYumHostedName, RES_ATTR_NAME, fmt.Sprintf("yum-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceYumHostedName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceYumHostedName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceYumHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceYumHostedName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceYumHostedName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceYumHostedName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
					resource.TestCheckNoResourceAttr(resourceYumHostedName, RES_ATTR_CLEANUP),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceYumProxyName, RES_ATTR_NAME, fmt.Sprintf("yum-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceYumProxyName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceYumProxyName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceYumProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceYumProxyName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "proxy.remote_url", "https://mirror.centos.org/centos/"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceYumProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceYumProxyName, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceYumProxyName, "replication.asset_path_regex"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceYumGroupName, RES_ATTR_NAME, fmt.Sprintf("yum-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceYumGroupName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceYumGroupName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceYumGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceYumGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryYumGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("yum-group-import-%s", randomString)
	memberName := fmt.Sprintf("yum-hosted-member-%s", randomString)

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
  yum = {
    repo_data_depth = 0
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
`, resourceTypeYumHosted, memberName, resourceTypeYumGroup, repoName, memberName, resourceTypeYumHosted),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceYumGroupName, RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceYumGroupName, RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceYumGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
