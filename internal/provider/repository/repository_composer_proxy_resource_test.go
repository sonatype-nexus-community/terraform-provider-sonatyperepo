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
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceNameComposerProxy = "sonatyperepo_repository_composer_proxy.repo"
	resourceTypeComposerProxy = "sonatyperepo_repository_composer_proxy"
)

func TestAccRepositoryComposerProxyResourceNoReplication(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getRepositoryComposerProxyResourceConfig(randomString, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "name", fmt.Sprintf("composer-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNameComposerProxy, "url"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "proxy.remote_url", "https://repo.packagist.org/"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceNameComposerProxy, "routing_rule"),
					resource.TestCheckResourceAttr(resourceNameComposerProxy, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceNameComposerProxy, "replication.asset_path_regex"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryComposerProxyResourceConfig(randomString string, includeReplication bool) string {
	var replicationConfig = ""
	if includeReplication {
		replicationConfig = `
	replication = {
		preemptive_pull_enabled = true
		asset_path_regex = "some-value"
	}	
`
	}
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "composer-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.packagist.org/"
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
  %s
}
`, resourceTypeComposerProxy, randomString, replicationConfig)
}
