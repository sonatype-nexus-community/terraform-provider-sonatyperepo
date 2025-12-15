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
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	dataSourcePrivileges = "data.sonatyperepo_privileges.ps"
)

func TestAccPrivilegesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test 1: Verify privileges can be listed
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_privileges" "ps" {
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourcePrivileges, "privileges.#"),
				),
			},
			// Test 2: Verify response structure and privilege attributes
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_privileges" "ps" {
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify count is greater than 0
					resource.TestCheckResourceAttrSet(dataSourcePrivileges, "privileges.#"),
					// Verify at least one privilege exists with expected attributes
					resource.TestCheckResourceAttrSet(dataSourcePrivileges, "privileges.0.name"),
					resource.TestCheckResourceAttrSet(dataSourcePrivileges, "privileges.0.description"),
					resource.TestCheckResourceAttrSet(dataSourcePrivileges, "privileges.0.read_only"),
					resource.TestCheckResourceAttrSet(dataSourcePrivileges, "privileges.0.type"),
				),
			},
		},
	})
}
