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

package repository_test

import (
	"fmt"
	"regexp"
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceNameCocoaPodsProxy = "sonatyperepo_repository_cocoapods_proxy.repo"
	resourceTypeCocoaPodsProxy = "sonatyperepo_repository_cocoapods_proxy"
)

func TestAccRepositoryCocoaPodsProxyResourceNoReplication(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Group validation - empty member_names
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_repository_cocoapods_group" "repo" {
  name = "cocoapods-group-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = []
  }
}
`, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			// Create and Read testing
			{
				Config: getRepositoryCocoaPodsProxyResourceConfig(randomString, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "name", fmt.Sprintf("cocoapods-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceNameCocoaPodsProxy, "url"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "proxy.remote_url", "https://cdn.cocoapods.org/"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceNameCocoaPodsProxy, "routing_rule"),
					resource.TestCheckResourceAttr(resourceNameCocoaPodsProxy, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceNameCocoaPodsProxy, "replication.asset_path_regex"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryCocoaPodsProxyResourceConfig(randomString string, includeReplication bool) string {
	var replicationConfig = ""
	if includeReplication {
		replicationConfig = `
	replication = {
		preemptive_pull_enabled = true
		asset_path_regex = "some-value"
	}	
`
	}
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://cdn.cocoapods.org/"
    content_max_age = 1441
    metadata_max_age = 1440
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
	connection = {
		enable_cookies = true
		retries = 9
		timeout = 999
		use_trust_store = true
		user_agent_suffix = "terraform"
	}
	authentication = {
		username = "user"
		password = "pass"
		preemptive = true
		type = "username"
	}
  }
  %s
}
`, resourceTypeCocoaPodsProxy, randomString, replicationConfig)
}

func TestAccRepositoryCocoaPodsProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-url-without-protocol"
    content_max_age = 1440
    metadata_max_age = 1440
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
  }
}
`, resourceTypeCocoaPodsProxy, randomString),
				ExpectError: regexp.MustCompile("must be a valid URL|must be a valid HTTP URL"),
			},
		},
	})
}

func TestAccRepositoryCocoapodsProxyInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  cocoapods = {}
}
`, "sonatyperepo_repository_cocoapods_proxy", randomString),
				ExpectError: regexp.MustCompile("Blob store.*not found|Blob store.*does not exist"),
			},
		},
	})
}

func TestAccRepositoryCocoapodsProxyMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "cocoapods-proxy-repo-%s"
  online = true
  # Missing storage block
}
`, "sonatyperepo_repository_cocoapods_proxy", randomString),
				ExpectError: regexp.MustCompile("Attribute storage is required"),
			},
		},
	})
}



