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
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBlobStoreGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			// {
			// 	Config: utils.ProviderConfig + `data "sonatyperepo_blob_store_file" "b" {
			// 		name = "default"
			// 	}`,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("data.sonatyperepo_blob_store_file.b", "path", "default"),
			// 		resource.TestCheckResourceAttrSet("data.sonatyperepo_blob_store_file.b", "soft_quota.%"),
			// 	),
			// },
			{
				Config: utils_test.ProviderConfig + `data "sonatyperepo_blob_store_group" "b" {
					name = "this-will-not-exist"
				}`,
				ExpectError: regexp.MustCompile("Unable to read group blob store"),
			},
		},
	})
}
