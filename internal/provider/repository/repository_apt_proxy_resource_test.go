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
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceNameAptProxy = "sonatyperepo_repository_apt_proxy.repo"
	resourceTypeAptProxy = "sonatyperepo_repository_apt_proxy"
)

func TestAccRepositoryAptProxyResourceNoReplication(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getRepositoryAptProxyResourceConfig(randomString, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameAptProxy, "name", fmt.Sprintf("apt-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNameAptProxy, "url"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "proxy.remote_url", "https://archive.ubuntu.com/ubuntu/"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "proxy.content_max_age", "1442"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "proxy.metadata_max_age", "1400"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceNameAptProxy, "routing_rule"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceNameAptProxy, "replication.asset_path_regex"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "apt.distribution", "bionic"),
					resource.TestCheckResourceAttr(resourceNameAptProxy, "apt.flat", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryAptProxyResourceConfig(randomString string, includeReplication bool) string {
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
  name = "apt-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://archive.ubuntu.com/ubuntu/"
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
  apt = {
	distribution = "bionic"
  }
  %s
}
`, resourceTypeAptProxy, randomString, replicationConfig)
}
