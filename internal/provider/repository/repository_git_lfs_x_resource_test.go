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
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceTypeGitLfsHosted = "sonatyperepo_repository_gitlfs_hosted"

var resourceGitLfsHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeGitLfsHosted)

func TestAccRepositoryGitLfsResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git-lfs-hosted-repo-%s"
   online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}
`, resourceTypeGitLfsHosted, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Proxy
					resource.TestCheckResourceAttr(resourceGitLfsHostedName, "name", fmt.Sprintf("git-lfs-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceGitLfsHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceGitLfsHostedName, "url"),
					resource.TestCheckResourceAttr(resourceGitLfsHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceGitLfsHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceGitLfsHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceGitLfsHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceGitLfsHostedName, "cleanup"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryGitLfsHostedInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git_lfs-hosted-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  git_lfs = {}
}
`, resourceTypeGitLfsHosted, randomString),
				ExpectError: regexp.MustCompile("Blob store.*not found|Blob store.*does not exist"),
			},
		},
	})
}

func TestAccRepositoryGitLfsHostedMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git_lfs-hosted-repo-%s"
  online = true
  # Missing storage block
}
`, resourceTypeGitLfsHosted, randomString),
				ExpectError: regexp.MustCompile("Attribute storage is required"),
			},
		},
	})
}
