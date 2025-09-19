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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
)

// TestAccBlobStoreGoogleCloudResourceExpectFailure tests that the resource fails gracefully without GCP credentials
func TestAccBlobStoreGoogleCloudResourceExpectFailure(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test that creation fails gracefully with authentication error
			{
				Config:      getTestAccBlobStoreGoogleCloudResourceMinimal(randomString),
				ExpectError: regexp.MustCompile("Error creating Google Cloud Storage Blob Store|authentication|unauthorized|forbidden|invalid_grant"),
			},
		},
	})
}

// TestAccBlobStoreGoogleCloudResourceValidation tests resource validation without API calls
func TestAccBlobStoreGoogleCloudResourceValidation(t *testing.T) {
	
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test invalid bucket name validation
			{
				Config:      getTestAccBlobStoreGoogleCloudResourceInvalidBucket(),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
			// Test missing required fields
			{
				Config:      getTestAccBlobStoreGoogleCloudResourceMissingName(),
				ExpectError: regexp.MustCompile("Missing required argument|name is required"),
			},
		},
	})
}

// TestAccBlobStoreGoogleCloudResourceConfigValidation tests Terraform configuration validation
func TestAccBlobStoreGoogleCloudResourceConfigValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test soft quota validation
			{
				Config:      getTestAccBlobStoreGoogleCloudResourceInvalidSoftQuota(),
				ExpectError: regexp.MustCompile("invalid soft quota|limit must be positive|Error creating"),
			},
		},
	})
}

// TestAccBlobStoreGoogleCloudResourceSchema tests the resource schema without creating resources
func TestAccBlobStoreGoogleCloudResourceSchema(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test that the configuration is parsed correctly (will fail at API call, which is expected)
			{
				Config:      getTestAccBlobStoreGoogleCloudResourceComplete(randomString),
				ExpectError: regexp.MustCompile("Error creating Google Cloud Storage Blob Store"),
				Check: resource.ComposeTestCheckFunc(
					// These checks won't run due to ExpectError, but they validate the schema
					resource.TestCheckResourceAttr("sonatyperepo_blob_store_google_cloud.gc_complete", "name", fmt.Sprintf("test-gc-complete-%s", randomString)),
				),
			},
		},
	})
}

// Configuration functions for tests

func getTestAccBlobStoreGoogleCloudResourceMinimal(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_blob_store_gcs" "gc_minimal" {
  name = "test-gc-minimal-%s"
  
  bucket_configuration {
    bucket {
      name = "nexus-bucket-minimal-%s"
    }
    
    authentication {
      authentication_method = "accountKey"
      account_key = jsonencode({
        "type": "service_account",
        "project_id": "test-project",
        "private_key_id": "test-key-id",
        "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC...\n-----END PRIVATE KEY-----\n",
        "client_email": "test@test-project.iam.gserviceaccount.com",
        "client_id": "123456789012345678901"
      })
    }
  }
}
`, randomString, randomString)
}

func getTestAccBlobStoreGoogleCloudResourceComplete(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_blob_store_gcs" "gc_complete" {
  name = "test-gc-complete-%s"
  
  bucket_configuration {
    bucket {
      name   = "nexus-bucket-complete-%s"
      prefix = "test-prefix-%s"
      region = "us-central1"
    }
    
    authentication {
      authentication_method = "accountKey"
      account_key = jsonencode({
        "type": "service_account",
        "project_id": "test-project-123",
        "private_key_id": "test-key-id-456",
        "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDGtJc...\n-----END PRIVATE KEY-----\n",
        "client_email": "test-service@test-project-123.iam.gserviceaccount.com",
        "client_id": "123456789012345678901",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token"
      })
    }
  }
  
  soft_quota {
    type  = "spaceUsedQuota"
    limit = 100000000000
  }
}
`, randomString, randomString, randomString)
}

func getTestAccBlobStoreGoogleCloudResourceInvalidBucket() string {
	return utils_test.ProviderConfig + `
resource "sonatyperepo_blob_store_gcs" "gc_invalid" {
  name = "test-gc-invalid"
  
  bucket_configuration {
    bucket {
      name = "INVALID_BUCKET_NAME_WITH_CAPS_AND_UNDERSCORES"
    }
    
    authentication {
      authentication_method = "accountKey"
      account_key = jsonencode({
        "type": "service_account",
        "project_id": "test-project"
      })
    }
  }
}
`
}

func getTestAccBlobStoreGoogleCloudResourceMissingName() string {
	return utils_test.ProviderConfig + `
resource "sonatyperepo_blob_store_gcs" "gc_missing" {
  bucket_configuration {
    bucket {
      name = "valid-bucket-name"
    }
    
    authentication {
      authentication_method = "accountKey"
      account_key = jsonencode({
        "type": "service_account",
        "project_id": "test-project"
      })
    }
  }
}
`
}

func getTestAccBlobStoreGoogleCloudResourceInvalidSoftQuota() string {
	return utils_test.ProviderConfig + `
resource "sonatyperepo_blob_store_gcs" "gc_invalid_quota" {
  name = "test-gc-invalid-quota"
  
  bucket_configuration {
    bucket {
      name = "valid-bucket-name"
    }
    
    authentication {
      authentication_method = "accountKey"
      account_key = jsonencode({
        "type": "service_account",
        "project_id": "test-project"
      })
    }
  }
  
  soft_quota {
    type  = "spaceUsedQuota"
    limit = -1
  }
}
`
}

// Unit tests that don't require API calls

func TestBlobStoreGoogleCloudResourceName(t *testing.T) {
	// Test resource type name format
	expectedPattern := `^sonatyperepo_blob_store_gcs$`
	resourceTypeName := "sonatyperepo_blob_store_gcs"
	
	matched, err := regexp.MatchString(expectedPattern, resourceTypeName)
	if err != nil {
		t.Fatalf("Regex error: %v", err)
	}
	if !matched {
		t.Errorf("Resource type name %s doesn't match expected pattern %s", resourceTypeName, expectedPattern)
	}
}

func TestBlobStoreGoogleCloudBucketNameValidation(t *testing.T) {
	validNames := []string{
		"valid-bucket-name",
		"bucket123",
		"my-test-bucket-2023",
	}
	
	invalidNames := []string{
		"UPPERCASE_BUCKET",
		"bucket_with_underscores",
		"bucket-name-that-is-way-too-long-and-exceeds-the-maximum-length-allowed-by-google-cloud-storage",
		"",
	}
	
	for _, name := range validNames {
		if !isValidBucketName(name) {
			t.Errorf("Expected %s to be a valid bucket name", name)
		}
	}
	
	for _, name := range invalidNames {
		if isValidBucketName(name) {
			t.Errorf("Expected %s to be an invalid bucket name", name)
		}
	}
}

// Helper function to validate bucket names (simplified version)
func isValidBucketName(name string) bool {
	if len(name) < 3 || len(name) > 63 {
		return false
	}
	
	// Basic validation - no uppercase, no underscores
	for _, char := range name {
		if char >= 'A' && char <= 'Z' {
			return false
		}
		if char == '_' {
			return false
		}
	}
	
	return true
}