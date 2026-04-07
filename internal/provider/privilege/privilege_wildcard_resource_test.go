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
	"terraform-provider-sonatyperepo/internal/provider/privilege/privilege_type"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceTypePrivilegeWildcard = "sonatyperepo_privilege_wildcard"
	resourceNamePrivilegeWildcard = fmt.Sprintf("%s.p", resourceTypeRepositoryContentSelector)
)

func TestAccPrivilegeWildcardResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: buildRepositortyContentSelectorResourceMinimal(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify minimal configuration
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "name", fmt.Sprintf("test-priv-wildcard-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "description", "some description"),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "read_only", "false"),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "type", privilege_type.TypeWildcard.String()),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "pattern", "test-pattern"),
				),
			},
			// Update to full configuration
			{
				Config: buildRepositoryContentSelectorResourceComplete(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify full configuration
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "name", fmt.Sprintf("test-priv-wildcard-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "description", "updated description"),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "read_only", "false"),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "type", privilege_type.TypeWildcard.String()),
					resource.TestCheckResourceAttr(resourceNameRepositoryContentSelector, "pattern", "updated-pattern-*"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func buildPrivilegeWildcardResourceMinimal(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "p" {
	name = "test-priv-wildcard-%s"
	description = "some description"
	pattern = "test-pattern"
}`, resourceTypeRepositoryContentSelector, randomString)
}

func buildPrivilegeWildcardResourceComplete(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "p" {
	name = "test-priv-wildcard-%s"
	description = "updated description"
	pattern = "updated-pattern-*"
}`, resourceTypeRepositoryContentSelector, randomString)
}
