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

	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
)

func TestAccRepositoriesDataSource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rawHostedRepoName := fmt.Sprintf("test-raw-hosted-ds-%s", randomString)
	rawHostedRepoName2 := fmt.Sprintf("test-raw-hosted-ds2-%s", randomString)
	dataSourceName := "data.sonatyperepo_repositories.repos"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test 1: Create test repositories and verify they appear in list
			{
				Config: testAccRepositoriesDataSourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that list is populated
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.#"),
					// Verify first created raw hosted repository appears in list
					resource.TestCheckTypeSetElemNestedAttrs(
						dataSourceName,
						"repositories.*",
						map[string]string{
							"name": rawHostedRepoName,
						},
					),
					// Verify second created raw hosted repository appears in list
					resource.TestCheckTypeSetElemNestedAttrs(
						dataSourceName,
						"repositories.*",
						map[string]string{
							"name": rawHostedRepoName2,
						},
					),
				),
			},
		},
	})
}

func testAccRepositoriesDataSourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_repository_raw_hosted" "test_hosted" {
  name   = "test-raw-hosted-ds-%s"
  online = true
  storage = {
    blob_store_name              = "default"
    strict_content_type_validation = true
    write_policy                 = "ALLOW_ONCE"
  }
  raw = {
    content_disposition = "ATTACHMENT"
  }
}

resource "sonatyperepo_repository_raw_hosted" "test_hosted2" {
  name   = "test-raw-hosted-ds2-%s"
  online = true
  storage = {
    blob_store_name              = "default"
    strict_content_type_validation = true
    write_policy                 = "ALLOW_ONCE"
  }
  raw = {
    content_disposition = "INLINE"
  }
  depends_on = [
    sonatyperepo_repository_raw_hosted.test_hosted
  ]
}

data "sonatyperepo_repositories" "repos" {
  depends_on = [
    sonatyperepo_repository_raw_hosted.test_hosted,
    sonatyperepo_repository_raw_hosted.test_hosted2
  ]
}
`, randomString, randomString)
}
