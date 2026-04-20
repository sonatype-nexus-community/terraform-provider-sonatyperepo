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
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	TEST_SKIPPED_ACS_BLOBSTORE string = "Azure Cloud Storage blob store resource tests require Azure credentials - set TF_ACC_ACS_BLOB_STORE=1 to enable"
)

// TestAccBlobStoreAcsResourceValidation tests ACS resource validation without API calls
func TestAccBlobStoreAcsResourceValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test schema validation - will fail at API call without Azure credentials
			{
				Config:      buildAcsResourceConfig("test-validation"),
				ExpectError: regexp.MustCompile("Error creating Azure Cloud Storage Blob Store"),
			},
		},
	})
}

// TestAccBlobStoreAcsResourceWithCredentials tests full ACS resource CRUD when Azure credentials are available
func TestAccBlobStoreAcsResourceWithCredentials(t *testing.T) {
	azureAccountId := os.Getenv("TF_ACC_AZURE_ACCOUNT_KEY")
	azureStorageAccountName := os.Getenv("TF_ACC_AZURE_STORAGE_ACCOUNT_NAME")
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			if os.Getenv("TF_ACC_ACS_BLOB_STORE") != "1" || azureStorageAccountName == "" || azureAccountId == "" {
				t.Skip(TEST_SKIPPED_ACS_BLOBSTORE)
			}
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: buildAcsResourceCompleteConfig(randomString, azureAccountId, azureStorageAccountName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_ACS, RES_ATTR_NAME, fmt.Sprintf("test-acs-complete-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_ACS, RES_ATTR_ACS_BUCKET_CONFIGURATION_ACCOUNT_NAME, azureStorageAccountName),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_ACS, RES_ATTR_ACS_BUCKET_CONFIGURATION_CONTAINER_NAME, fmt.Sprintf("testacs%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_ACS, RES_ATTR_ACS_BUCKET_CONFIGURATION_AUTH_AUTHENTICATION_METHOD, common.BLOB_STORE_ACS_AUTH_METHOD_ACCOUNT_KEY),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_ACS, RES_ATTR_ACS_BUCKET_CONFIGURATION_AUTH_ACCOUNT_KEY, azureAccountId),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_ACS, RES_ATTR_SOFT_QUOTA_LIMIT, "1099511627776"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_ACS, RES_ATTR_SOFT_QUOTA_TYPE, "spaceUsedQuota"),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_ACS, RES_ATTR_LAST_UPDATED),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         RES_NAME_BLOB_STORE_ACS,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        fmt.Sprintf("test-acs-complete-%s", randomString),
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{RES_ATTR_ACS_BUCKET_CONFIGURATION_AUTH_ACCOUNT_KEY, RES_ATTR_LAST_UPDATED},
			},
		},
	})
}

// Configuration builder functions
func buildAcsResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-acs-%s"
  bucket_configuration = {
    account_name = "something"
	container_name = "anything"
	authentication = {
		authentication_method = "%s"
		account_key = "rubbish-%s-rubbish"
	}
  }
}
`, RES_TYPE_BLOB_STORE_ACS, randomString, common.BLOB_STORE_ACS_AUTH_METHOD_ACCOUNT_KEY, randomString)
}

func buildAcsResourceCompleteConfig(randomString, azureAccountId, azureStorageAccountName string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-acs-complete-%s"
  bucket_configuration = {
    account_name = "%s"
	container_name = "testacs%s"
	authentication = {
		authentication_method = "%s"
		account_key = "%s"
	}
  }
  soft_quota = {
	type = "spaceUsedQuota"
	limit = 1099511627776
  }
}
`, RES_TYPE_BLOB_STORE_ACS, randomString, azureStorageAccountName, randomString, common.BLOB_STORE_ACS_AUTH_METHOD_ACCOUNT_KEY, azureAccountId)
}
