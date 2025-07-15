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
	resourceNamePyPiHosted = "sonatyperepo_repository_pypi_hosted.repo"
	resourceTypePyPiHosted = "sonatyperepo_repository_pypi_hosted"
)

func TestAccRepositoryPyPiHostedResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getRepositoryPyPiHostedResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNamePyPiHosted, "name", fmt.Sprintf("pypi-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNamePyPiHosted, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNamePyPiHosted, "url"),
					resource.TestCheckResourceAttr(resourceNamePyPiHosted, "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr(resourceNamePyPiHosted, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNamePyPiHosted, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceNamePyPiHosted, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceNamePyPiHosted, "cleanup"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryPyPiHostedResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pypi-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}
`, resourceTypePyPiHosted, randomString)
}
