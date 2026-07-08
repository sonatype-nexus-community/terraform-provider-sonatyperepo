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

package alpine_test

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
	resourceTypeAlpineGroup  = "sonatyperepo_repository_alpine_group"
	resourceTypeAlpineHosted = "sonatyperepo_repository_alpine_hosted"
	resourceTypeAlpineProxy  = "sonatyperepo_repository_alpine_proxy"
)

var (
	resourceAlpineGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeAlpineGroup)
	resourceAlpineHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeAlpineHosted)
	resourceAlpineProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeAlpineProxy)
)

func alpineSkipPreCheck(t *testing.T) func() {
	return func() {
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
	}
}

func TestAccRepositoryAlpineResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 alpineSkipPreCheck(t),
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "alpine-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
}
`, resourceTypeAlpineGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "alpine-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  alpine = {
    key_pair = "123"
    passphrase = "123"
  }
}

resource "%s" "repo" {
  name = "alpine-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://dl-cdn.alpinelinux.org/alpine/"
    content_max_age = 1441
    metadata_max_age = 1440
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  alpine = {
    key_pair = "123"
    passphrase = "123"
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
  name = "alpine-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["alpine-proxy-repo-%s"]
  }
  alpine = {
    key_pair = "123"
    passphrase = "123"
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeAlpineHosted, randomString, resourceTypeAlpineProxy, randomString, resourceTypeAlpineGroup, randomString, randomString, resourceTypeAlpineProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceAlpineHostedName, repotest.RES_ATTR_NAME, fmt.Sprintf("alpine-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceAlpineHostedName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceAlpineHostedName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceAlpineHostedName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceAlpineHostedName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceAlpineHostedName, repotest.RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckNoResourceAttr(resourceAlpineHostedName, repotest.RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS),
					resource.TestCheckNoResourceAttr(resourceAlpineHostedName, repotest.RES_ATTR_CLEANUP),
					resource.TestCheckResourceAttr(resourceAlpineHostedName, "alpine.key_pair", "123"),
					resource.TestCheckResourceAttr(resourceAlpineHostedName, "alpine.passphrase", "123"),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_NAME, fmt.Sprintf("alpine-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceAlpineProxyName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_PROXY_REMOTE_URL, "https://dl-cdn.alpinelinux.org/alpine/"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_PROXY_CONTENT_MAX_AGE, "1441"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_PROXY_METADATA_MAX_AGE, "1440"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_NEGATIVE_CACHE_ENABLED, "true"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, "1440"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_BLOCKED, "false"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, "true"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, "false"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, "true"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, "true"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, "9"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, "999"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX, "terraform"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_USERNAME, "user"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PASSWORD, "pass"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PREMPTIVE, "true"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_TYPE, "username"),
					resource.TestCheckNoResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_ROUTING_RULE_NAME),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED, "false"),
					resource.TestCheckNoResourceAttr(resourceAlpineProxyName, repotest.RES_ATTR_REPLICATION_ASSET_PATH_REGEX),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, "alpine.key_pair", "123"),
					resource.TestCheckResourceAttr(resourceAlpineProxyName, "alpine.passphrase", "123"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceAlpineGroupName, repotest.RES_ATTR_NAME, fmt.Sprintf("alpine-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceAlpineGroupName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceAlpineGroupName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceAlpineGroupName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceAlpineGroupName, repotest.RES_ATTR_GROUP_MEMBER_NAMES, "1"),
					resource.TestCheckResourceAttr(resourceAlpineGroupName, "alpine.key_pair", "123"),
					resource.TestCheckResourceAttr(resourceAlpineGroupName, "alpine.passphrase", "123"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryAlpineGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("alpine-group-import-%s", randomString)
	memberName := fmt.Sprintf("alpine-hosted-member-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck:                 alpineSkipPreCheck(t),
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
  alpine = {
    key_pair = "123"
  }
  depends_on = [%s.member]
}
`, resourceTypeAlpineHosted, memberName, resourceTypeAlpineGroup, repoName, memberName, resourceTypeAlpineHosted),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAlpineGroupName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceAlpineGroupName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttr(resourceAlpineGroupName, "alpine.key_pair", "123"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceAlpineGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				// alpineSigning is write-only and never returned by the GET API, so it can't be verified after import
				ImportStateVerifyIgnore: []string{"last_updated", "alpine.%", "alpine.key_pair"},
			},
		},
	})
}
