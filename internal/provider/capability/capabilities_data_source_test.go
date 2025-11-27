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

package capability_test

import (
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCapabilitiesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Capabilities data source not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Test 1: Verify capabilities can be listed
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_capabilities" "caps" {
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonatyperepo_capabilities.caps", "capabilities.#"),
				),
			},
			// Test 2: Verify response structure and capability attributes
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_capabilities" "caps" {
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify count is greater than 0
					resource.TestCheckResourceAttrSet("data.sonatyperepo_capabilities.caps", "capabilities.#"),
					// Verify at least one capability exists with expected attributes
					resource.TestCheckResourceAttrSet("data.sonatyperepo_capabilities.caps", "capabilities.0.id"),
					resource.TestCheckResourceAttrSet("data.sonatyperepo_capabilities.caps", "capabilities.0.type"),
					resource.TestCheckResourceAttrSet("data.sonatyperepo_capabilities.caps", "capabilities.0.enabled"),
				),
			},
		},
	})
}
