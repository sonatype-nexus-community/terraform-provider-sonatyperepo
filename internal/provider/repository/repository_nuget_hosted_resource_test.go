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
	resourceNameNugetHosted = "sonatyperepo_repository_nuget_hosted.repo"
	resourceTypeNugetHosted = "sonatyperepo_repository_nuget_hosted"
)

func TestAccRepositoryNugetHostedResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getRepositoryNugetHostedResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameNugetHosted, "name", fmt.Sprintf("nuget-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameNugetHosted, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNameNugetHosted, "url"),
					resource.TestCheckResourceAttr(resourceNameNugetHosted, "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr(resourceNameNugetHosted, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNameNugetHosted, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceNameNugetHosted, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceNameNugetHosted, "cleanup"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryNugetHostedResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "nuget-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}
`, resourceTypeNugetHosted, randomString)
}
