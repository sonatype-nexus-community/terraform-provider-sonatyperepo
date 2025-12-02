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
	"os"
	"regexp"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceNameS3BlobStore = "sonatyperepo_blob_store_s3.test"

// TestAccBlobStoreS3ResourceValidation tests S3 resource validation without API calls
func TestAccBlobStoreS3ResourceValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test schema validation - will fail at API call without AWS credentials
			{
				Config:      buildS3ResourceConfig("test-validation"),
				ExpectError: regexp.MustCompile("Error creating S3 Blob Store|InvalidAccessKeyId|NoSuchBucket"),
			},
		},
	})
}

// TestAccBlobStoreS3ResourceWithCredentials tests full S3 resource CRUD when AWS credentials are available
func TestAccBlobStoreS3ResourceWithCredentials(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			if os.Getenv("TF_ACC_S3_BLOB_STORE") != "1" {
				t.Skip("S3 blob store resource tests require AWS credentials - set TF_ACC_S3_BLOB_STORE=1 to enable")
			}
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: buildS3ResourceCompleteConfig("test-crud"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameS3BlobStore, "name"),
					resource.TestCheckResourceAttrSet(resourceNameS3BlobStore, "bucket_configuration.0.bucket.0.region"),
				),
			},
			// Update testing
			{
				Config: buildS3ResourceCompleteConfig("test-crud-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameS3BlobStore, "name"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:      resourceNameS3BlobStore,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Configuration builder functions

func buildS3ResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_blob_store_s3" "test" {
  name = "test-s3-%s"
  bucket_configuration = {
	bucket = {
		region = "eu-west-2"
		name = "nexus-bucket-%s"
		prefix = "prefix-%s"
		expiration = 99
	}
  }
}
`, randomString, randomString, randomString)
}

func buildS3ResourceCompleteConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_blob_store_s3" "test" {
  name = "test-s3-complete-%s"
  bucket_configuration = {
	bucket = {
		region = "eu-west-2"
		name = "nexus-bucket-complete-%s"
		prefix = "prefix-%s"
		expiration = 99
	}
	bucket_security = {
		access_key_id = "AKIAIOSFODNN7EXAMPLE"
		secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	}
  }
  soft_quota = {
	type = "spaceUsedQuota"
	limit = 1099511627776
  }
}
`, randomString, randomString, randomString)
}
