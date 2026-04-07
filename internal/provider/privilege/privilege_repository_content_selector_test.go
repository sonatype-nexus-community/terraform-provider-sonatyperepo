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

package privilege_test

import (
	"fmt"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceTypeRepositoryContentSelector = "sonatyperepo_privilege_repository_content_selector"
	resourceNameRepositoryContentSelector = fmt.Sprintf("%s.p", resourceTypeRepositoryContentSelector)
)

func TestAccRepositoryContentSelectorResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: buildRepositortyContentSelectorResourceMinimal(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify minimal configuration
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "name", fmt.Sprintf("test-repo-cs-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "description", ""),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "content_selector", fmt.Sprintf("test-content-selector-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "format", "*"),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "repository", "maven-central"),
				),
			},
			// Update to full configuration
			{
				Config: buildRepositoryContentSelectorResourceComplete(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify full configuration
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "name", fmt.Sprintf("test-repo-cs-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "description", "a description"),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "content_selector", fmt.Sprintf("test-content-selector-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "format", "*"),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "repository", "maven-central"),
				),
			},
			// Delete testing automatically occurs in TestCase
			// ImportState testing
			{
				ResourceName:                         resourceNameRepositoryContentSelector,
				ImportStateId:                        fmt.Sprintf("test-repo-cs-%s", randomString),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func buildRepositortyContentSelectorResourceMinimal(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_content_selector" "cs" {
	name = "test-content-selector-%s"
	description = ""
	expression = "format == \"maven2\" and path =^ \"/org/sonatype/%s\""
}
resource "%s" "p" {
	name = "test-repo-cs-%s"
	description = ""
	actions = ["BROWSE"]
	content_selector = "test-content-selector-%s"
	format = "*"
  	repository = "maven-central"

	depends_on = [ sonatyperepo_content_selector.cs ]
}`, randomString, randomString, resourceTypeRepositoryContentSelector, randomString, randomString)
}

func buildRepositoryContentSelectorResourceComplete(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_content_selector" "cs" {
	name = "test-content-selector-%s"
	description = ""
	expression = "format == \"maven2\" and path =^ \"/org/sonatype/%s\""
}
resource "%s" "p" {
	name = "test-repo-cs-%s"
	description = "a description"
	actions = ["BROWSE"]
	content_selector = "test-content-selector-%s"
	format = "*"
  	repository = "maven-central"

	depends_on = [ sonatyperepo_content_selector.cs ]
}`, randomString, randomString, resourceTypeRepositoryContentSelector, randomString, randomString)
}
