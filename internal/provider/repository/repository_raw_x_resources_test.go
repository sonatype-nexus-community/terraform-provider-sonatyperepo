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

func TestAccRepositoryRawResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceTypeGroup := "sonatyperepo_repository_raw_group"
	resourceTypeHosted := "sonatyperepo_repository_raw_hosted"
	resourceTypeProxy := "sonatyperepo_repository_raw_proxy"
	resourceGroupName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeGroup)
	resourceHostedName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeHosted)
	resourceProxyName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeProxy)

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
`, resourceTypeGroup, randomString),
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
`, resourceTypeHosted, randomString, resourceTypeProxy, randomString, resourceTypeGroup, randomString, randomString, resourceTypeProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceHostedName, "name", fmt.Sprintf("raw-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceHostedName, "url"),
					resource.TestCheckResourceAttr(resourceHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceHostedName, "cleanup"),
					resource.TestCheckResourceAttr(resourceHostedName, RES_ATTR_RAW_CONTENT_DISPOSITION, common.CONTENT_DISPOSITION_ATTACHMENT),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceProxyName, "name", fmt.Sprintf("raw-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceProxyName, "url"),
					resource.TestCheckResourceAttr(resourceProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceProxyName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "proxy.remote_url", "https://nodejs.org/dist/"),
					resource.TestCheckResourceAttr(resourceProxyName, "proxy.content_max_age", "1442"),
					resource.TestCheckResourceAttr(resourceProxyName, "proxy.metadata_max_age", "1400"),
					resource.TestCheckResourceAttr(resourceProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceProxyName, RES_ATTR_RAW_CONTENT_DISPOSITION, common.CONTENT_DISPOSITION_ATTACHMENT),

					// Verify Group
					resource.TestCheckResourceAttr(resourceGroupName, "name", fmt.Sprintf("raw-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceGroupName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceGroupName, "url"),
					resource.TestCheckResourceAttr(resourceGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceGroupName, "group.member_names.#", "1"),
					resource.TestCheckResourceAttr(resourceGroupName, RES_ATTR_RAW_CONTENT_DISPOSITION, common.CONTENT_DISPOSITION_ATTACHMENT),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
