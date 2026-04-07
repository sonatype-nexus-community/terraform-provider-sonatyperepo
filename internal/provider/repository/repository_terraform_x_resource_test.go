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
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeTerraformGroup  = "sonatyperepo_repository_terraform_group"
	resourceTypeTerraformHosted = "sonatyperepo_repository_terraform_hosted"
	resourceTypeTerraformProxy  = "sonatyperepo_repository_terraform_proxy"
)

var (
	resourceTerraformGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeTerraformGroup)
	resourceTerraformHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeTerraformHosted)
	resourceTerraformProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeTerraformProxy)
)

func TestAccRepositoryTerraformResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Only works on NXRM 3.90.0 or later
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 89,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "terraform-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  %s
}

resource "%s" "repo" {
  name = "terraform-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://registry.terraform.io"
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
  terraform = { }
}

resource "%s" "repo" {
  name = "terraform-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["terraform-proxy-repo-%s"]
  }
  terraform = { }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeTerraformHosted, randomString, configBlockHostedDefaultTerraform, resourceTypeTerraformProxy, randomString, resourceTypeTerraformGroup, randomString, randomString, resourceTypeTerraformProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceTerraformHostedName, RES_ATTR_NAME, fmt.Sprintf("terraform-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceTerraformHostedName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceTerraformHostedName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceTerraformHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceTerraformHostedName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceTerraformHostedName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckNoResourceAttr(resourceTerraformHostedName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS),
					resource.TestCheckNoResourceAttr(resourceTerraformHostedName, RES_ATTR_CLEANUP),
					resource.TestCheckResourceAttrSet(resourceTerraformHostedName, RES_ATTR_TERRAFORM_SIGNING_SIGNING_KEY),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_NAME, fmt.Sprintf("terraform-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceTerraformProxyName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_PROXY_REMOTE_URL, "https://registry.terraform.io"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_PROXY_CONTENT_MAX_AGE, "1442"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_PROXY_METADATA_MAX_AGE, "1400"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_NEGATIVE_CACHE_ENABLED, "true"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, "1440"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_BLOCKED, "false"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, "true"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, "false"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, "true"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, "true"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, "9"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, "999"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX, "terraform"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION_USERNAME, "user"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PASSWORD, "pass"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PREMPTIVE, "true"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION_TYPE, "username"),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED, "false"),
					resource.TestCheckNoResourceAttr(resourceTerraformProxyName, RES_ATTR_REPLICATION_ASSET_PATH_REGEX),
					resource.TestCheckNoResourceAttr(resourceTerraformProxyName, RES_ATTR_ROUTING_RULE_NAME),
					resource.TestCheckNoResourceAttr(resourceTerraformProxyName, RES_ATTR_REPOSITORY_FIREWALL),
					resource.TestCheckResourceAttr(resourceTerraformProxyName, RES_ATTR_TERRAFORM_REQUIRE_AUTH, "false"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceTerraformGroupName, RES_ATTR_NAME, fmt.Sprintf("terraform-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceTerraformGroupName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceTerraformGroupName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceTerraformGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceTerraformGroupName, RES_ATTR_GROUP_MEMBER_NAMES, "1"),
					resource.TestCheckResourceAttr(resourceTerraformGroupName, RES_ATTR_TERRAFORM_REQUIRE_AUTH, "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryTerraformGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("terraform-group-import-%s", randomString)
	memberName := fmt.Sprintf("terraform-hosted-member-%s", randomString)

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
  %s
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
  terraform = {
    require_authentication = true
  }
  depends_on = [%s.member]
}
`, resourceTypeTerraformHosted, memberName, configBlockHostedDefaultTerraform, resourceTypeTerraformGroup, repoName, memberName, resourceTypeTerraformHosted),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Group
					resource.TestCheckResourceAttr(resourceTerraformGroupName, RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceTerraformGroupName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceTerraformGroupName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceTerraformGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceTerraformGroupName, RES_ATTR_GROUP_MEMBER_NAMES, "1"),
					resource.TestCheckResourceAttr(resourceTerraformGroupName, RES_ATTR_TERRAFORM_REQUIRE_AUTH, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceTerraformGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
