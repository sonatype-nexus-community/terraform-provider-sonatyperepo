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

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepositoryMavenProxyResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getRepositoryMavenProxyResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "name", fmt.Sprintf("maven-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "format", REPOSITORY_FORMAT_MAVEN),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "type", REPOSITORY_TYPE_PROXY),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "online", "true"),
					resource.TestCheckResourceAttrSet("sonatyperepo_repository_maven_proxy.repo", "url"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "proxy.remote_url", "https://repo1.maven.org/maven2/"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "http_client.blocked", "false"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "http_client.connection.enable_cookies", "false"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_proxy.repo", "http_client.connection.use_trust_store", "false"),
					resource.TestCheckNoResourceAttr("sonatyperepo_repository_maven_proxy.repo", "http_client.connection.authentication"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryMavenProxyResourceConfig(randomString string) string {
	return fmt.Sprintf(providerConfig+`
resource "sonatyperepo_repository_maven_proxy" "repo" {
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
  }
  maven = {
	content_disposition = "ATTACHMENT"
	layout_policy = "STRICT"
	version_policy = "RELEASE"
  }
}
`, randomString)
}
