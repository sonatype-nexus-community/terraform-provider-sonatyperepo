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

package helm_test

import (
	"fmt"
	"regexp"
	"terraform-provider-sonatyperepo/internal/provider/common"
	repotest "terraform-provider-sonatyperepo/internal/provider/repository/repotest"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	"testing"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeHelmGroup  = "sonatyperepo_repository_helm_group"
	resourceTypeHelmHosted = "sonatyperepo_repository_helm_hosted"
	resourceTypeHelmProxy  = "sonatyperepo_repository_helm_proxy"
)

var (
	resourceHelmGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeHelmGroup)
	resourceHelmHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeHelmHosted)
)

func helmGroupPreCheck(t *testing.T) {
	// Helm Group repositories are available from NXRM 3.92.0
	testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
		Major: 3,
		Minor: 0,
		Patch: 0,
	}, &common.SystemVersion{
		Major: 3,
		Minor: 91,
		Patch: 99,
	})
}

func TestAccRepositoryHelmGroupResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { helmGroupPreCheck(t) },
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Validation: empty member_names should fail
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "helm-group-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = []
  }
}
`, resourceTypeHelmGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			// Create hosted + group
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "helm-hosted-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
    write_policy = "ALLOW_ONCE"
  }
}

resource "%s" "repo" {
  name = "helm-group-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = ["helm-hosted-repo-%s"]
  }

  depends_on = [
    %s.repo
  ]
}
`, resourceTypeHelmHosted, randomString, resourceTypeHelmGroup, randomString, randomString, resourceTypeHelmHosted),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceHelmHostedName, repotest.RES_ATTR_NAME, fmt.Sprintf("helm-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceHelmHostedName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceHelmHostedName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceHelmHostedName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceHelmHostedName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceHelmHostedName, repotest.RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceHelmHostedName, repotest.RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
					resource.TestCheckNoResourceAttr(resourceHelmHostedName, repotest.RES_ATTR_CLEANUP),

					// Verify Group
					resource.TestCheckResourceAttr(resourceHelmGroupName, repotest.RES_ATTR_NAME, fmt.Sprintf("helm-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceHelmGroupName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceHelmGroupName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceHelmGroupName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceHelmGroupName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceHelmGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryHelmGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("helm-group-import-%s", randomString)
	memberName := fmt.Sprintf("helm-hosted-member-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { helmGroupPreCheck(t) },
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "member" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
    write_policy = "ALLOW"
  }
}

resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = ["%s"]
  }
  depends_on = [%s.member]
}
`, resourceTypeHelmHosted, memberName, resourceTypeHelmGroup, repoName, memberName, resourceTypeHelmHosted),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceHelmGroupName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceHelmGroupName, repotest.RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceHelmGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
