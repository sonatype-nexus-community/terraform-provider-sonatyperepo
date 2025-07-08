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
	"terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceNameMavenHosted = "sonatyperepo_repository_maven_hosted.repo"
	resourceTypeMavenHosted = "sonatyperepo_repository_maven_hosted"
)

func TestAccRepositoryMavenHostedResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getRepositoryMavenHostedResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameMavenHosted, "name", fmt.Sprintf("maven-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameMavenHosted, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNameMavenHosted, "url"),
					resource.TestCheckResourceAttr(resourceNameMavenHosted, "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr(resourceNameMavenHosted, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNameMavenHosted, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceNameMavenHosted, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceNameMavenHosted, "cleanup"),
					resource.TestCheckResourceAttr(resourceNameMavenHosted, "maven.content_disposition", common.MAVEN_CONTENT_DISPOSITION_ATTACHMENT),
					resource.TestCheckResourceAttr(resourceNameMavenHosted, "maven.layout_policy", common.MAVEN_LAYOUT_STRICT),
					resource.TestCheckResourceAttr(resourceNameMavenHosted, "maven.version_policy", common.MAVEN_VERSION_POLICY_RELEASE),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryMavenHostedResourceConfig(randomString string) string {
	return fmt.Sprintf(utils.ProviderConfig+`
resource "%s" "repo" {
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
`, resourceTypeMavenHosted, randomString)
}
