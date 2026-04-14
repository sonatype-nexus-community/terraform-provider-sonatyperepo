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
				ExpectError: regexp.MustCompile(errMessageBlobStoreS3ErrorCreating),
			},
		},
	})
}

// TestAccBlobStoreS3ResourceWithCredentials tests full S3 resource CRUD when AWS credentials are available
func TestAccBlobStoreS3ResourceWithCredentials(t *testing.T) {
	awsAccessKeyId := os.Getenv("TF_ACC_AWS_ACCESS_KEY_ID")
	awsAccessSecretKey := os.Getenv("TF_ACC_AWS_ACCESS_SECRET_KEY")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			if os.Getenv("TF_ACC_S3_BLOB_STORE") != "1" || awsAccessKeyId == "" || awsAccessSecretKey == "" {
				t.Skip(TEST_SKIPPED_S3_BLOBSTORE)
			}
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: buildS3ResourceCompleteConfig("test-crud", awsAccessKeyId, awsAccessSecretKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, "test-s3-complete-test-crud"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_REGION, awsRegionEuWest2),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_NAME, "nexus-bucket-complete-test-crud"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_PREFIX, "prefix-test-crud"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_KEY_ID, awsAccessKeyId),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_SECRET_KEY, awsAccessSecretKey),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED, "false"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_LIMIT, "1099511627776"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_TYPE, "spaceUsedQuota"),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_S3, RES_ATTR_LAST_UPDATED),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         RES_NAME_BLOB_STORE_S3,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "test-s3-complete-test-crud",
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"bucket_configuration.bucket_security.secret_access_key", "last_updated"},
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
		region = "%s"
		name = "nexus-bucket-%s"
		prefix = "prefix-%s"
	}
  }
}
`, RES_TYPE_BLOB_STORE_S3, randomString, awsRegionEuWest2, randomString, randomString)
}

func buildS3ResourceCompleteConfig(randomString, awsAccessKeyId, awsAccessSecretKey string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-s3-complete-%s"
  bucket_configuration = {
	bucket = {
		region = "%s"
		name = "nexus-bucket-complete-%s"
		prefix = "prefix-%s"
	}
	bucket_security = {
		access_key_id = "%s"
		secret_access_key = "%s"
	}
	pre_signed_url_enabled = false
  }
  soft_quota = {
	type = "spaceUsedQuota"
	limit = 1099511627776
  }
}
`, RES_TYPE_BLOB_STORE_S3, randomString, awsRegionEuWest2, randomString, randomString, awsAccessKeyId, awsAccessSecretKey)
}

// TestAccBlobStoreS3ResourcePreSignedUrlValidation tests pre-signed URL validation without API calls
func TestAccBlobStoreS3ResourcePreSignedUrlValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test with pre_signed_url_enabled = true
			{
				Config:      buildS3ResourcePreSignedUrlConfig("test-presigned-true", true),
				ExpectError: regexp.MustCompile(errMessageBlobStoreS3ErrorCreating),
			},
			// Test with pre_signed_url_enabled = false
			{
				Config:      buildS3ResourcePreSignedUrlConfig("test-presigned-false", false),
				ExpectError: regexp.MustCompile(errMessageBlobStoreS3ErrorCreating),
			},
			// Test with pre_signed_url_enabled omitted (should default to false)
			{
				Config:      buildS3ResourcePreSignedUrlConfigOmitted("test-presigned-default", "rubbish", "invalid"),
				ExpectError: regexp.MustCompile(errMessageBlobStoreS3ErrorCreating),
			},
		},
	})
}

// TestAccBlobStoreS3ResourcePreSignedUrlWithCredentials tests pre-signed URL CRUD when AWS credentials are available
func TestAccBlobStoreS3ResourcePreSignedUrlWithCredentials(t *testing.T) {
	awsAccessKeyId := os.Getenv("TF_ACC_AWS_ACCESS_KEY_ID")
	awsAccessSecretKey := os.Getenv("TF_ACC_AWS_ACCESS_SECRET_KEY")
	randomString := "test-presigned-crud"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			if os.Getenv("TF_ACC_S3_BLOB_STORE") != "1" || awsAccessKeyId == "" || awsAccessSecretKey == "" {
				t.Skip(TEST_SKIPPED_S3_BLOBSTORE)
			}
		},
		Steps: []resource.TestStep{
			// Create with pre_signed_url_enabled = false
			{
				Config: buildS3ResourcePreSignedUrlCompleteConfig(randomString, awsAccessKeyId, awsAccessSecretKey, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, fmt.Sprintf(testS3NamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_REGION, awsRegionEuWest2),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_NAME, fmt.Sprintf(testS3BucketNamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_PREFIX, fmt.Sprintf(testS3BucketPrefix, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_KEY_ID, awsAccessKeyId),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_SECRET_KEY, awsAccessSecretKey),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED, "false"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_LIMIT, "1099511627776"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_TYPE, "spaceUsedQuota"),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_S3, RES_ATTR_LAST_UPDATED),
				),
			},
			// Update to pre_signed_url_enabled = true
			{
				Config: buildS3ResourcePreSignedUrlCompleteConfig(randomString, awsAccessKeyId, awsAccessSecretKey, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, fmt.Sprintf(testS3NamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_REGION, awsRegionEuWest2),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_NAME, fmt.Sprintf(testS3BucketNamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_PREFIX, fmt.Sprintf(testS3BucketPrefix, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_KEY_ID, awsAccessKeyId),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_SECRET_KEY, awsAccessSecretKey),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED, "true"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_LIMIT, "1099511627776"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_TYPE, "spaceUsedQuota"),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_S3, RES_ATTR_LAST_UPDATED),
				),
			},
			// Update back to pre_signed_url_enabled = false
			{
				Config: buildS3ResourcePreSignedUrlCompleteConfig(randomString, awsAccessKeyId, awsAccessSecretKey, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, fmt.Sprintf(testS3NamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_REGION, awsRegionEuWest2),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_NAME, fmt.Sprintf(testS3BucketNamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_PREFIX, fmt.Sprintf(testS3BucketPrefix, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_KEY_ID, awsAccessKeyId),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_SECRET_KEY, awsAccessSecretKey),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED, "false"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_LIMIT, "1099511627776"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_TYPE, "spaceUsedQuota"),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_S3, RES_ATTR_LAST_UPDATED),
				),
			},
			// Import and verify
			{
				ResourceName:                         RES_NAME_BLOB_STORE_S3,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        fmt.Sprintf("test-s3-presigned-%s", randomString),
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"bucket_configuration.bucket_security.secret_access_key", "last_updated"},
			},
		},
	})
}

// TestAccBlobStoreS3ResourcePreSignedUrlDefault tests that pre_signed_url_enabled defaults to false
func TestAccBlobStoreS3ResourcePreSignedUrlDefault(t *testing.T) {
	awsAccessKeyId := os.Getenv("TF_ACC_AWS_ACCESS_KEY_ID")
	awsAccessSecretKey := os.Getenv("TF_ACC_AWS_ACCESS_SECRET_KEY")
	randomString := "test-presigned-default"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			if os.Getenv("TF_ACC_S3_BLOB_STORE") != "1" || awsAccessKeyId == "" || awsAccessSecretKey == "" {
				t.Skip(TEST_SKIPPED_S3_BLOBSTORE)
			}
		},
		Steps: []resource.TestStep{
			// Create without specifying pre_signed_url_enabled
			{
				Config: buildS3ResourcePreSignedUrlCompleteConfigOmitted(randomString, awsAccessKeyId, awsAccessSecretKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, fmt.Sprintf(testS3NamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_REGION, awsRegionEuWest2),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_NAME, fmt.Sprintf(testS3BucketNamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_PREFIX, fmt.Sprintf(testS3BucketPrefix, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_KEY_ID, awsAccessKeyId),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_SECRET_KEY, awsAccessSecretKey),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED, "false"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_LIMIT, "1099511627776"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_TYPE, "spaceUsedQuota"),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_S3, RES_ATTR_LAST_UPDATED),
				),
			},
		},
	})
}

func TestAccBlobStoreS3ResourceStateUpgradeV0ToV1(t *testing.T) {
	awsAccessKeyId := os.Getenv("TF_ACC_AWS_ACCESS_KEY_ID")
	awsAccessSecretKey := os.Getenv("TF_ACC_AWS_ACCESS_SECRET_KEY")
	randomString := "upgrade-v0-v1"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("TF_ACC_S3_BLOB_STORE") != "1" || awsAccessKeyId == "" || awsAccessSecretKey == "" {
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
				Config: buildS3ResourcePreSignedUrlConfigOmitted(randomString, awsAccessKeyId, awsAccessSecretKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, fmt.Sprintf(testS3NamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_REGION, awsRegionEuWest2),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_NAME, fmt.Sprintf(testS3BucketNamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_PREFIX, fmt.Sprintf(testS3BucketPrefix, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_KEY_ID, awsAccessKeyId),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_SECRET_KEY, awsAccessSecretKey),
					resource.TestCheckNoResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED),
					resource.TestCheckNoResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_LIMIT),
					resource.TestCheckNoResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_TYPE),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_S3, RES_ATTR_LAST_UPDATED),
				),
			},
			{
				// Step 2: Upgrade to current provider version
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"sonatyperepo": providerserver.NewProtocol6WithError(provider.New("test")()),
				},
				Config: buildS3ResourcePreSignedUrlConfigOmitted(randomString, awsAccessKeyId, awsAccessSecretKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_NAME, fmt.Sprintf(testS3NamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_REGION, awsRegionEuWest2),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_NAME, fmt.Sprintf(testS3BucketNamePresigned, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_PREFIX, fmt.Sprintf(testS3BucketPrefix, randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_KEY_ID, awsAccessKeyId),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_SECRET_KEY, awsAccessSecretKey),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED, "false"),
					resource.TestCheckNoResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_LIMIT),
					resource.TestCheckNoResourceAttr(RES_NAME_BLOB_STORE_S3, RES_ATTR_SOFT_QUOTA_TYPE),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_S3, RES_ATTR_LAST_UPDATED),
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
	}
	pre_signed_url_enabled = %t
  }
}
`, RES_TYPE_BLOB_STORE_S3, randomString, randomString, randomString, preSignedUrlEnabled)
}

func buildS3ResourcePreSignedUrlConfigOmitted(randomString, awsAccessKeyId, awsAccessSecretKey string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-s3-presigned-%s"
  bucket_configuration = {
	bucket = {
		region = "eu-west-2"
		name = "nexus-bucket-presigned-%s"
		prefix = "prefix-%s"
	}
	bucket_security = {
		access_key_id = "%s"
		secret_access_key = "%s"
	}
  }
}
`, RES_TYPE_BLOB_STORE_S3, randomString, randomString, randomString, awsAccessKeyId, awsAccessSecretKey)
}

func buildS3ResourcePreSignedUrlCompleteConfig(randomString, awsAccessKeyId, awsAccessSecretKey string, preSignedUrlEnabled bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-s3-presigned-%s"
  bucket_configuration = {
	bucket = {
		region = "eu-west-2"
		name = "nexus-bucket-presigned-%s"
		prefix = "prefix-%s"
	}
	bucket_security = {
		access_key_id = "%s"
		secret_access_key = "%s"
	}
	pre_signed_url_enabled = %t
  }
  soft_quota = {
	type = "spaceUsedQuota"
	limit = 1099511627776
  }
}
`, RES_TYPE_BLOB_STORE_S3, randomString, randomString, randomString, awsAccessKeyId, awsAccessSecretKey, preSignedUrlEnabled)
}

func buildS3ResourcePreSignedUrlCompleteConfigOmitted(randomString, awsAccessKeyId, awsAccessSecretKey string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-s3-presigned-%s"
  bucket_configuration = {
	bucket = {
		region = "eu-west-2"
		name = "nexus-bucket-presigned-%s"
		prefix = "prefix-%s"
	}
	bucket_security = {
		access_key_id = "%s"
		secret_access_key = "%s"
	}
  }
  soft_quota = {
	type = "spaceUsedQuota"
	limit = 1099511627776
  }
}
`, RES_TYPE_BLOB_STORE_S3, randomString, randomString, randomString, awsAccessKeyId, awsAccessSecretKey)
}
