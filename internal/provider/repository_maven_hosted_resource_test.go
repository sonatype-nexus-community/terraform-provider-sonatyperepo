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

func TestAccRepositoryMavenHostedResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getRepositoryMavenHostedResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_hosted.repo", "name", fmt.Sprintf("maven-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_hosted.repo", "format", REPOSITORY_FORMAT_MAVEN),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_hosted.repo", "type", REPOSITORY_TYPE_HOSTED),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_hosted.repo", "online", "true"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_hosted.repo", "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr("sonatyperepo_repository_maven_hosted.repo", "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr("sonatyperepo_repository_maven_hosted.repo", "cleanup"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryMavenHostedResourceConfig(randomString string) string {
	return fmt.Sprintf(providerConfig+`
resource "sonatyperepo_repository_maven_hosted" "repo" {
  name = "maven-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  maven = {
	content_disposition = "ATTACHMENT"
	layout_policy = "STRICT"
	version_policy = "RELEASE"
  }
}
`, randomString)
}
