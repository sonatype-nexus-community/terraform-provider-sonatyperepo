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
	resourceTypeRawGroup  = "sonatyperepo_repository_raw_group"
	resourceTypeRawHosted = "sonatyperepo_repository_raw_hosted"
	resourceTypeRawProxy  = "sonatyperepo_repository_raw_proxy"
)

var (
	resourceRawGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRawGroup)
	resourceRawHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRawHosted)
	resourceRawProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeRawProxy)
)

func TestAccRepositoryRawResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
  raw = {
	content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  raw = {
	content_disposition = "ATTACHMENT"
  }
}

resource "%s" "repo" {
  name = "raw-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://nodejs.org/dist/"
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
  raw = {
	content_disposition = "ATTACHMENT"
  }
}

resource "%s" "repo" {
  name = "raw-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["raw-proxy-repo-%s"]
  }
  raw = {
	content_disposition = "ATTACHMENT"
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeRawHosted, randomString, resourceTypeRawProxy, randomString, resourceTypeRawGroup, randomString, randomString, resourceTypeRawProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceRawHostedName, RES_ATTR_NAME, fmt.Sprintf("raw-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRawHostedName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceRawHostedName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceRawHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRawHostedName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceRawHostedName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceRawHostedName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
					resource.TestCheckNoResourceAttr(resourceRawHostedName, RES_ATTR_CLEANUP),
					resource.TestCheckResourceAttr(resourceRawHostedName, RES_ATTR_RAW_CONTENT_DISPOSITION, common.CONTENT_DISPOSITION_ATTACHMENT),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceRawProxyName, RES_ATTR_NAME, fmt.Sprintf("raw-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRawProxyName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceRawProxyName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceRawProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRawProxyName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "proxy.remote_url", "https://nodejs.org/dist/"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "proxy.content_max_age", "1442"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "proxy.metadata_max_age", "1400"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceRawProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceRawProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceRawProxyName, RES_ATTR_RAW_CONTENT_DISPOSITION, common.CONTENT_DISPOSITION_ATTACHMENT),

					// Verify Group
					resource.TestCheckResourceAttr(resourceRawGroupName, RES_ATTR_NAME, fmt.Sprintf("raw-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceRawGroupName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceRawGroupName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceRawGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceRawGroupName, "group.member_names.#", "1"),
					resource.TestCheckResourceAttr(resourceRawGroupName, RES_ATTR_RAW_CONTENT_DISPOSITION, common.CONTENT_DISPOSITION_ATTACHMENT),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryRawGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceType := "sonatyperepo_repository_raw_group"
	resourceTypeHosted := "sonatyperepo_repository_raw_hosted"
	resourceName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceType)
	repoName := fmt.Sprintf("raw-group-import-%s", randomString)
	memberName := fmt.Sprintf("raw-hosted-member-%s", randomString)

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
  raw = {
    content_disposition = "ATTACHMENT"
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
  raw = {
    content_disposition = "ATTACHMENT"
  }
  depends_on = [%s.member]
}
`, resourceTypeHosted, memberName, resourceType, repoName, memberName, resourceTypeHosted),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
