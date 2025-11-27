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

package user_test

import (
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	dataSourceUsers = "data.sonatyperepo_users.us"
)

func TestAccUsersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test 1: Verify users can be listed
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_users" "us" {
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceUsers, "users.#"),
				),
			},
			// Test 2: Verify response structure and user attributes
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_users" "us" {
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify count is greater than 0
					resource.TestCheckResourceAttrSet(dataSourceUsers, "users.#"),
					// Verify at least one user exists with expected attributes
					resource.TestCheckResourceAttrSet(dataSourceUsers, "users.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSourceUsers, "users.0.first_name"),
					resource.TestCheckResourceAttrSet(dataSourceUsers, "users.0.last_name"),
					resource.TestCheckResourceAttrSet(dataSourceUsers, "users.0.email_address"),
					resource.TestCheckResourceAttrSet(dataSourceUsers, "users.0.read_only"),
					resource.TestCheckResourceAttrSet(dataSourceUsers, "users.0.source"),
					resource.TestCheckResourceAttrSet(dataSourceUsers, "users.0.status"),
					// Verify user has roles assigned
					resource.TestCheckResourceAttrSet(dataSourceUsers, "users.0.roles.#"),
					// Verify user has external roles assigned
					resource.TestCheckResourceAttrSet(dataSourceUsers, "users.0.external_roles.#"),
				),
			},
		},
	})
}
