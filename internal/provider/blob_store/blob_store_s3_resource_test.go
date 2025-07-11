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
	"fmt"
	"regexp"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TODO: Acceptance Tests do not have access to an environment that has appropriate AWS connectivity.
func TestAccBlobStoreS3Resource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	// resourceName := "sonatyperepo_blob_store_s3.b"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      getTestAccBlobStoreS3ResourceWillFail(randomString),
				ExpectError: regexp.MustCompile("Error creating S3 Blob Store"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getTestAccBlobStoreS3ResourceWillFail(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_blob_store_s3" "b" {
  name = "test-%s"
  bucket_configuration = {
	bucket = {
		region = "eu-west-2"
		name = "bucket-name-%s"
		prefix = "prefix-%s"
		expiration = 99
	}
  }
}
`, randomString, randomString, randomString)
}
