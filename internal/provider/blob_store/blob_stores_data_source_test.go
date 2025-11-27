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

package blob_store_test

import (
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBlobStoresDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test 1: Verify blob stores can be listed
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_blob_stores" "blob_stores" {
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonatyperepo_blob_stores.blob_stores", "blob_stores.#"),
				),
			},
			// Test 2: Verify response structure and blob store attributes
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_blob_stores" "blob_stores" {
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify count is greater than 0
					resource.TestCheckResourceAttrSet("data.sonatyperepo_blob_stores.blob_stores", "blob_stores.#"),
					// Verify at least one blob store exists and has expected attributes
					resource.TestCheckResourceAttrSet("data.sonatyperepo_blob_stores.blob_stores", "blob_stores.0.name"),
					resource.TestCheckResourceAttrSet("data.sonatyperepo_blob_stores.blob_stores", "blob_stores.0.type"),
					// Verify the default 'file' type blob store exists
					resource.TestCheckTypeSetElemNestedAttrs("data.sonatyperepo_blob_stores.blob_stores", "blob_stores.*", map[string]string{
						"name": "default",
						"type": "File",
					}),
				),
			},
		},
	})
}
