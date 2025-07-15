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
	resourceNameRawHosted = "sonatyperepo_repository_raw_hosted.repo"
	resourceTypeRawHosted = "sonatyperepo_repository_raw_hosted"
)

func TestAccRepositoryRawHostedResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getRepositoryRawHostedResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameRawHosted, "name", fmt.Sprintf("raw-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRawHosted, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNameRawHosted, "url"),
					resource.TestCheckResourceAttr(resourceNameRawHosted, "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr(resourceNameRawHosted, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNameRawHosted, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceNameRawHosted, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceNameRawHosted, "cleanup"),
					resource.TestCheckResourceAttr(resourceNameRawHosted, "raw.content_disposition", "ATTACHMENT"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryRawHostedResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "raw-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  raw = {
	content_disposition = "ATTACHMENT"
  }
}
`, resourceTypeRawHosted, randomString)
}
