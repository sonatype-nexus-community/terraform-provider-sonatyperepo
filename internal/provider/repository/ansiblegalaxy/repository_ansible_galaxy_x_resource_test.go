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

package ansiblegalaxy_test

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
	resourceTypeAnsibleGalaxyGroup  = "sonatyperepo_repository_ansiblegalaxy_group"
	resourceTypeAnsibleGalaxyHosted = "sonatyperepo_repository_ansiblegalaxy_hosted"
	resourceTypeAnsibleGalaxyProxy  = "sonatyperepo_repository_ansiblegalaxy_proxy"
)

var (
	resourceAnsibleGalaxyGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeAnsibleGalaxyGroup)
	resourceAnsibleGalaxyHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeAnsibleGalaxyHosted)
	resourceAnsibleGalaxyProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeAnsibleGalaxyProxy)
)

func TestAccRepositoryAnsibleGalaxyResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Requires NXRM 3.93.0+
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 92,
				Patch: 99,
			})
		},
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "ansiblegalaxy-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
}
`, resourceTypeAnsibleGalaxyGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "ansiblegalaxy-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}

resource "%s" "repo" {
  name = "ansiblegalaxy-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://galaxy.ansible.com"
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
  name = "ansiblegalaxy-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["ansiblegalaxy-proxy-repo-%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeAnsibleGalaxyHosted, randomString, resourceTypeAnsibleGalaxyProxy, randomString, resourceTypeAnsibleGalaxyGroup, randomString, randomString, resourceTypeAnsibleGalaxyProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyHostedName, repotest.RES_ATTR_NAME, fmt.Sprintf("ansiblegalaxy-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyHostedName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceAnsibleGalaxyHostedName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyHostedName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyHostedName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyHostedName, repotest.RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyHostedName, repotest.RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
					resource.TestCheckNoResourceAttr(resourceAnsibleGalaxyHostedName, repotest.RES_ATTR_CLEANUP),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_NAME, fmt.Sprintf("ansiblegalaxy-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_PROXY_REMOTE_URL, "https://galaxy.ansible.com"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_PROXY_CONTENT_MAX_AGE, "1441"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_PROXY_METADATA_MAX_AGE, "1440"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_NEGATIVE_CACHE_ENABLED, "true"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, "1440"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_BLOCKED, "false"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, "true"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, "false"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, "true"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, "true"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, "9"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, "999"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX, "terraform"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_USERNAME, "user"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PASSWORD, "pass"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PREMPTIVE, "true"),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_TYPE, "username"),
					resource.TestCheckNoResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_ROUTING_RULE_NAME),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED, "false"),
					resource.TestCheckNoResourceAttr(resourceAnsibleGalaxyProxyName, repotest.RES_ATTR_REPLICATION_ASSET_PATH_REGEX),

					// Verify Group
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyGroupName, repotest.RES_ATTR_NAME, fmt.Sprintf("ansiblegalaxy-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyGroupName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceAnsibleGalaxyGroupName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyGroupName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyGroupName, repotest.RES_ATTR_GROUP_MEMBER_NAMES, "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryAnsibleGalaxyGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("ansiblegalaxy-group-import-%s", randomString)
	memberName := fmt.Sprintf("ansiblegalaxy-hosted-member-%s", randomString)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Requires NXRM 3.93.0+
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 92,
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
`, resourceTypeAnsibleGalaxyHosted, memberName, resourceTypeAnsibleGalaxyGroup, repoName, memberName, resourceTypeAnsibleGalaxyHosted),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyGroupName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceAnsibleGalaxyGroupName, repotest.RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceAnsibleGalaxyGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
