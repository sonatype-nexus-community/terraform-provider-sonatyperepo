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
	resourceNamePrivilegeWildcard = fmt.Sprintf("%s.p", resourceTypePrivilegeWildcard)
)

func TestAccPrivilegeWildcardResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "p" {
	name = "test-priv-wildcard-%s"
	description = "some description"
	pattern = "test-pattern"
}`, resourceTypePrivilegeWildcard, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNamePrivilegeWildcard, "name", fmt.Sprintf("test-priv-wildcard-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNamePrivilegeWildcard, "description", "some description"),
					resource.TestCheckResourceAttr(resourceNamePrivilegeWildcard, "read_only", "false"),
					resource.TestCheckResourceAttr(resourceNamePrivilegeWildcard, "type", privilege_type.TypeWildcard.String()),
					resource.TestCheckResourceAttr(resourceNamePrivilegeWildcard, "pattern", "test-pattern"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
