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

func TestAccRepositoryGitLfsResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceTypeHosted := "sonatyperepo_repository_gitlfs_hosted"
	resourceTypeGroup := "sonatyperepo_repository_gitlfs_group"
	resourceHostedName := fmt.Sprintf("%s.repo", resourceTypeHosted)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Group validation - empty member_names
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git-lfs-group-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = []
  }
}
`, resourceTypeGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git-lfs-hosted-repo-%s"
   online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
}
`, resourceTypeHosted, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Proxy
					resource.TestCheckResourceAttr(resourceHostedName, "name", fmt.Sprintf("git-lfs-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceHostedName, "url"),
					resource.TestCheckResourceAttr(resourceHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceHostedName, "cleanup"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccRepositoryGitLfsProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git_lfs-proxy-repo-%s"
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
`, "sonatyperepo_repository_git_lfs_proxy", randomString),
				ExpectError: regexp.MustCompile("must be a valid URL|must be a valid HTTP URL"),
			},
		},
	})
}

func TestAccRepositoryGitLfsHostedInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git_lfs-hosted-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  git_lfs = {}
}
`, "sonatyperepo_repository_git_lfs_hosted", randomString),
				ExpectError: regexp.MustCompile("Blob store.*not found|Blob store.*does not exist"),
			},
		},
	})
}

func TestAccRepositoryGitLfsHostedMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git_lfs-hosted-repo-%s"
  online = true
  # Missing storage block
}
`, "sonatyperepo_repository_git_lfs_hosted", randomString),
				ExpectError: regexp.MustCompile("Attribute storage is required"),
			},
		},
	})
}



func TestAccRepositoryGitLfsProxyInvalidTimeoutTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too large, max is 3600)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git_lfs-proxy-repo-timeout-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      timeout = 3601
    }
  }
}
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be less than or equal to 3600"),
			},
		},
	})
}

func TestAccRepositoryGitLfsProxyInvalidTimeoutTooSmall(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too small, min is 1)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git_lfs-proxy-repo-timeout-small-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      timeout = 0
    }
  }
}
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be greater than or equal to 1"),
			},
		},
	})
}

func TestAccRepositoryGitLfsProxyInvalidRetriesTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (too large, max is 10)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git_lfs-proxy-repo-retries-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      retries = 11
    }
  }
}
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be less than or equal to 10"),
			},
		},
	})
}

func TestAccRepositoryGitLfsProxyInvalidRetriesNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git_lfs-proxy-repo-retries-neg-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      retries = -1
    }
  }
}
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be greater than or equal to 0"),
			},
		},
	})
}

func TestAccRepositoryGitLfsProxyInvalidMaxAgeNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid content_max_age (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git_lfs-proxy-repo-maxage-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
    content_max_age = -1
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
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be greater than or equal to|cannot be negative"),
			},
		},
	})
}

func TestAccRepositoryGitLfsProxyInvalidTimeToLiveNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid time_to_live (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "git_lfs-proxy-repo-ttl-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
    content_max_age = 1440
    metadata_max_age = 1440
  }
  negative_cache = {
    enabled = true
    time_to_live = -1
  }
  http_client = {
    blocked = false
    auto_block = true
  }
}
`, resourceTypeProxy, randomString),
				ExpectError: regexp.MustCompile("must be greater than or equal to|cannot be negative"),
			},
		},
	})
}

