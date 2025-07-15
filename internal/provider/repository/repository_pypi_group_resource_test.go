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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
)

func TestAccRepositoryPyPiGroupResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := "sonatyperepo_repository_pypi_group.repo"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      getRepositoryPyPiGroupResourceConfigNoMembers(randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: getRepositoryPyPiGroupResourceConfigWithMembers(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(

					// Verify
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("pypi-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "group.member_names.#", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "group.writable_member"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryPyPiGroupResourceConfigNoMembers(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_repository_pypi_group" "repo" {
  name = "pypi-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
}
`, randomString)
}

func getRepositoryPyPiGroupResourceConfigWithMembers(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_repository_pypi_group" "repo" {
  name = "pypi-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["pupy-proxy"]
  }
}
`, randomString)
}
