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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-sonatyperepo/internal/provider/utils"
)

const (
	resourceNameMavenProxy = "sonatyperepo_repository_maven_proxy.repo"
	resourceTypeMavenProxy = "sonatyperepo_repository_maven_proxy"
)

func TestAccRepositoryMavenProxyResourceNoReplication(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	// resourceName := "sonatyperepo_repository_maven_proxy.repo"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getRepositoryMavenProxyResourceConfig(randomString, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "name", fmt.Sprintf("maven-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNameMavenProxy, "url"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "proxy.remote_url", "https://repo1.maven.org/maven2/"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceNameMavenProxy, "routing_rule"),
					resource.TestCheckResourceAttr(resourceNameMavenProxy, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceNameMavenProxy, "replication.asset_path_regex"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// Replication can only be from an NXRM instance and we have no instance to Acceptance Test this against
//
// func TestAccRepositoryMavenProxyResourceWithReplication(t *testing.T) {

// 	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
// 	// resourceName := "sonatyperepo_repository_maven_proxy.repo"

// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: utils.TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Create and Read testing
// 			{
// 				Config: getRepositoryMavenProxyResourceConfig(randomString, true),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					// Verify
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "name", fmt.Sprintf("maven-proxy-repo-%s", randomString)),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "online", "true"),
// 					resource.TestCheckResourceAttrSet(resourceNameMavenProxy, "url"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "storage.blob_store_name", "default"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "storage.strict_content_type_validation", "true"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "storage.write_policy", common.WRITE_POLICY_ALLOW),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "proxy.remote_url", "https://repo1.maven.org/maven2/"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "proxy.content_max_age", "1441"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "proxy.metadata_max_age", "1440"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "negative_cache.enabled", "true"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "negative_cache.time_to_live", "1440"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.blocked", "false"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.auto_block", "true"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.enable_circular_redirects", "false"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.enable_cookies", "true"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.use_trust_store", "true"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.retries", "9"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.timeout", "999"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.connection.user_agent_suffix", "terraform"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.authentication.username", "user"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.authentication.password", "pass"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.authentication.preemptive", "true"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "http_client.authentication.type", "username"),
// 					resource.TestCheckNoResourceAttr(resourceNameMavenProxy, "routing_rule"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "replication.preemptive_pull_enabled", "true"),
// 					resource.TestCheckResourceAttr(resourceNameMavenProxy, "replication.asset_path_regex", "some-value"),
// 				),
// 			},
// 			// Delete testing automatically occurs in TestCase
// 		},
// 	})
// }

func getRepositoryMavenProxyResourceConfig(randomString string, includeReplication bool) string {
	var replicationConfig = ""
	if includeReplication {
		replicationConfig = `
	replication = {
		preemptive_pull_enabled = true
		asset_path_regex = "some-value"
	}	
`
	}
	return fmt.Sprintf(utils.ProviderConfig+`
resource "%s" "repo" {
  name = "maven-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo1.maven.org/maven2/"
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
  maven = {
	content_disposition = "ATTACHMENT"
	layout_policy = "STRICT"
	version_policy = "RELEASE"
  }
  %s
}
`, resourceTypeMavenProxy, randomString, replicationConfig)
}
