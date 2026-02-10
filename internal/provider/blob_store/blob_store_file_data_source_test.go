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
	"regexp"
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	dataSourceBlobStoreFile = "data.sonatyperepo_blob_store_file.b"
)

func TestAccBlobStoreFileDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test 1: Missing required argument
			{
				Config:      utils_test.ProviderConfig + `data "sonatyperepo_blob_store_file" "b" {}`,
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
			// Test 2: Happy path - default blob store exists and is readable
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_blob_store_file" "b" {
					name = "default"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceBlobStoreFile, "name", "default"),
					resource.TestCheckResourceAttr(dataSourceBlobStoreFile, "path", "default"),
					// Soft quota is absent in default config
					resource.TestCheckNoResourceAttr(dataSourceBlobStoreFile, "soft_quota"),
				),
			},
			// Test 3: Non-existent blob store
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_blob_store_file" "b" {
					name = "this-will-not-exist"
				}`,
				ExpectError: regexp.MustCompile("Error: " + common.ERROR_UNABLE_TO_READ_BLOB_STORE_FILE),
			},
			// Test 4: Invalid name values
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_blob_store_file" "b" {
					name = ""
				}`,
				ExpectError: regexp.MustCompile("Error:"),
			},
		},
	})
}
