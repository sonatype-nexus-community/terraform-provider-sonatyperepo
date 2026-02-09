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
	"terraform-provider-sonatyperepo/internal/provider"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	TEST_SKIPPED_S3_BLOBSTORE string = "S3 blob store resource tests require AWS credentials - set TF_ACC_S3_BLOB_STORE=1 to enable"
)

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
				t.Skip(TEST_SKIPPED_S3_BLOBSTORE)
			}
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: buildS3ResourceCompleteConfig("test-crud"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, "test-s3-complete-test-crud"),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_S3, "bucket_configuration.0.bucket.0.region"),
				),
			},
			// Update testing
			{
				Config: buildS3ResourceCompleteConfig("test-crud-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, "test-s3-complete-test-crud-updated"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:      RES_NAME_BLOB_STORE_S3,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Configuration builder functions

func buildS3ResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
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
`, RES_TYPE_BLOB_STORE_S3, randomString, randomString, randomString)
}

func buildS3ResourceCompleteConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-s3-complete-%s"
  bucket_configuration = {
	bucket = {
		region = "eu-west-2"
		name = "nexus-bucket-complete-%s"
		prefix = "prefix-%s"
		expiration = 99
	}
	bucket_security = {
		access_key_id = "not-a-valid-aws-key"
		secret_access_key = "not-a-valid-aws-secret-key"
	}
	pre_signed_url_enabled = false
  }
  soft_quota = {
	type = "spaceUsedQuota"
	limit = 1099511627776
  }
}
`, RES_TYPE_BLOB_STORE_S3, randomString, randomString, randomString)
}

// TestAccBlobStoreS3ResourcePreSignedUrlValidation tests pre-signed URL validation without API calls
func TestAccBlobStoreS3ResourcePreSignedUrlValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test with pre_signed_url_enabled = true
			{
				Config:      buildS3ResourcePreSignedUrlConfig("test-presigned-true", true),
				ExpectError: regexp.MustCompile("Error creating S3 Blob Store|InvalidAccessKeyId|NoSuchBucket"),
			},
			// Test with pre_signed_url_enabled = false
			{
				Config:      buildS3ResourcePreSignedUrlConfig("test-presigned-false", false),
				ExpectError: regexp.MustCompile("Error creating S3 Blob Store|InvalidAccessKeyId|NoSuchBucket"),
			},
			// Test with pre_signed_url_enabled omitted (should default to false)
			{
				Config:      buildS3ResourcePreSignedUrlConfigOmitted("test-presigned-default"),
				ExpectError: regexp.MustCompile("Error creating S3 Blob Store|InvalidAccessKeyId|NoSuchBucket"),
			},
		},
	})
}

// TestAccBlobStoreS3ResourcePreSignedUrlWithCredentials tests pre-signed URL CRUD when AWS credentials are available
func TestAccBlobStoreS3ResourcePreSignedUrlWithCredentials(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			if os.Getenv("TF_ACC_S3_BLOB_STORE") != "1" {
				t.Skip(TEST_SKIPPED_S3_BLOBSTORE)
			}
		},
		Steps: []resource.TestStep{
			// Create with pre_signed_url_enabled = false
			{
				Config: buildS3ResourcePreSignedUrlCompleteConfig("test-presigned-crud", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, "test-s3-presigned-test-presigned-crud"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED, "false"),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_REGION),
				),
			},
			// Update to pre_signed_url_enabled = true
			{
				Config: buildS3ResourcePreSignedUrlCompleteConfig("test-presigned-crud", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, "test-s3-presigned-test-presigned-crud"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED, "true"),
				),
			},
			// Update back to pre_signed_url_enabled = false
			{
				Config: buildS3ResourcePreSignedUrlCompleteConfig("test-presigned-crud", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, "test-s3-presigned-test-presigned-crud"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED, "false"),
				),
			},
			// Import and verify
			{
				ResourceName:      RES_NAME_BLOB_STORE_S3,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccBlobStoreS3ResourcePreSignedUrlDefault tests that pre_signed_url_enabled defaults to false
func TestAccBlobStoreS3ResourcePreSignedUrlDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			if os.Getenv("TF_ACC_S3_BLOB_STORE") != "1" {
				t.Skip(TEST_SKIPPED_S3_BLOBSTORE)
			}
		},
		Steps: []resource.TestStep{
			// Create without specifying pre_signed_url_enabled
			{
				Config: buildS3ResourcePreSignedUrlCompleteConfigOmitted("test-presigned-default"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, "test-s3-presigned-test-presigned-default"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED, "false"),
				),
			},
		},
	})
}

func TestAccBlobStoreS3ResourceStateUpgradeV0ToV1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("TF_ACC_S3_BLOB_STORE") != "1" {
				t.Skip(TEST_SKIPPED_S3_BLOBSTORE)
			}
		},
		Steps: []resource.TestStep{
			{
				// Step 1: Create resource with old provider version (v0.17.0)
				ExternalProviders: map[string]resource.ExternalProvider{
					"sonatyperepo": {
						Source:            "sonatype-nexus-community/sonatyperepo",
						VersionConstraint: "0.17.0", // Last version without pre_signed_url_enabled
					},
				},
				Config: buildS3ResourcePreSignedUrlConfigOmitted("upgrade-v0-v1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_S3), RES_ATTR_NAME, "test-s3-presigned-upgrade-v0-v1"),
					resource.TestCheckResourceAttr(fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_S3), RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_REGION, "eu-west-2"),
				),
			},
			{
				// Step 2: Upgrade to current provider version
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"sonatyperepo": providerserver.NewProtocol6WithError(provider.New("test")()),
				},
				Config: buildS3ResourcePreSignedUrlConfigOmitted("upgrade-v0-v1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify pre_signed_url_enabled was added by state upgrade
					resource.TestCheckResourceAttr(
						fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_S3),
						RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED,
						"false",
					),
					// Verify other fields unchanged
					resource.TestCheckResourceAttr(fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_S3), RES_ATTR_NAME, "test-s3-presigned-upgrade-v0-v1"),
					resource.TestCheckResourceAttr(fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_S3), RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_REGION, "eu-west-2"),
				),
			},
		},
	})
}

// Configuration builder functions for pre-signed URL tests

func buildS3ResourcePreSignedUrlConfig(randomString string, preSignedUrlEnabled bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-s3-presigned-%s"
  bucket_configuration = {
	bucket = {
		region = "eu-west-2"
		name = "nexus-bucket-presigned-%s"
		prefix = "prefix-%s"
		expiration = 99
	}
	pre_signed_url_enabled = %t
  }
}
`, RES_TYPE_BLOB_STORE_S3, randomString, randomString, randomString, preSignedUrlEnabled)
}

func buildS3ResourcePreSignedUrlConfigOmitted(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-s3-presigned-%s"
  bucket_configuration = {
	bucket = {
		region = "eu-west-2"
		name = "nexus-bucket-presigned-%s"
		prefix = "prefix-%s"
		expiration = 99
	}
  }
}
`, RES_TYPE_BLOB_STORE_S3, randomString, randomString, randomString)
}

func buildS3ResourcePreSignedUrlCompleteConfig(randomString string, preSignedUrlEnabled bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-s3-presigned-%s"
  bucket_configuration = {
	bucket = {
		region = "eu-west-2"
		name = "nexus-bucket-presigned-%s"
		prefix = "prefix-%s"
		expiration = 99
	}
	bucket_security = {
		access_key_id = "not-a-valid-aws-key"
		secret_access_key = "not-a-valid-aws-secret-key"
	}
	pre_signed_url_enabled = %t
  }
  soft_quota = {
	type = "spaceUsedQuota"
	limit = 1099511627776
  }
}
`, RES_TYPE_BLOB_STORE_S3, randomString, randomString, randomString, preSignedUrlEnabled)
}

func buildS3ResourcePreSignedUrlCompleteConfigOmitted(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-s3-presigned-%s"
  bucket_configuration = {
	bucket = {
		region = "eu-west-2"
		name = "nexus-bucket-presigned-%s"
		prefix = "prefix-%s"
		expiration = 99
	}
	bucket_security = {
		access_key_id = "not-a-valid-aws-key"
		secret_access_key = "not-a-valid-aws-secret-key"
	}
  }
  soft_quota = {
	type = "spaceUsedQuota"
	limit = 1099511627776
  }
}
`, RES_TYPE_BLOB_STORE_S3, randomString, randomString, randomString)
}
