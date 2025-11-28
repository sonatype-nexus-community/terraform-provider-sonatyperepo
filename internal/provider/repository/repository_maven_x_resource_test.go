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
	resourceTypeMavenGroup  = "sonatyperepo_repository_maven_group"
	resourceTypeMavenHosted = "sonatyperepo_repository_maven_hosted"
	resourceTypeMavenProxy  = "sonatyperepo_repository_maven_proxy"
)

var (
	resourceMavenGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeMavenGroup)
	resourceMavenHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeMavenHosted)
	resourceMavenProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeMavenProxy)
)

func TestAccRepositoryMavenResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "maven-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
}
`, resourceTypeMavenGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "maven-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  maven = {
	content_disposition = "ATTACHMENT"
	layout_policy = "STRICT"
	version_policy = "RELEASE"
  }
}

resource "%s" "repo" {
  name = "maven-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo1.maven.org/maven2/"
    content_max_age = 1442
    metadata_max_age = 1400
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
  maven = {
	content_disposition = "ATTACHMENT"
	layout_policy = "STRICT"
	version_policy = "RELEASE"
  }
}

resource "%s" "repo" {
  name = "maven-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["maven-proxy-repo-%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeMavenHosted, randomString, resourceTypeMavenProxy, randomString, resourceTypeMavenGroup, randomString, randomString, resourceTypeMavenProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceMavenHostedName, "name", fmt.Sprintf("maven-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceMavenHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceMavenHostedName, "url"),
					resource.TestCheckResourceAttr(resourceMavenHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceMavenHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceMavenHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceMavenHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceMavenHostedName, "cleanup"),
					resource.TestCheckResourceAttr(resourceMavenHostedName, "maven.content_disposition", common.MAVEN_CONTENT_DISPOSITION_ATTACHMENT),
					resource.TestCheckResourceAttr(resourceMavenHostedName, "maven.layout_policy", common.MAVEN_LAYOUT_STRICT),
					resource.TestCheckResourceAttr(resourceMavenHostedName, "maven.version_policy", common.MAVEN_VERSION_POLICY_RELEASE),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceMavenProxyName, "name", fmt.Sprintf("maven-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceMavenProxyName, "url"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "proxy.remote_url", "https://repo1.maven.org/maven2/"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "proxy.content_max_age", "1442"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "proxy.metadata_max_age", "1400"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "http_client.authentication.password", "pass"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceMavenGroupName, "name", fmt.Sprintf("maven-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceMavenGroupName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceMavenGroupName, "url"),
					resource.TestCheckResourceAttr(resourceMavenGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceMavenGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryMavenHostedImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("maven-hosted-import-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
    write_policy = "ALLOW_ONCE"
  }
  maven = {
    content_disposition = "ATTACHMENT"
    layout_policy = "STRICT"
    version_policy = "RELEASE"
  }
}
`, resourceTypeMavenHosted, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceMavenHostedName, "name", repoName),
					resource.TestCheckResourceAttr(resourceMavenHostedName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceMavenHostedName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryMavenProxyImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("maven-proxy-import-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo1.maven.org/maven2/"
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
  maven = {
    content_disposition = "ATTACHMENT"
    layout_policy = "STRICT"
    version_policy = "RELEASE"
  }
}
`, resourceTypeMavenProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceMavenProxyName, "name", repoName),
					resource.TestCheckResourceAttr(resourceMavenProxyName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceMavenProxyName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryMavenGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("maven-group-import-%s", randomString)
	memberName := fmt.Sprintf("maven-hosted-member-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "member" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
    write_policy = "ALLOW"
  }
  maven = {
    content_disposition = "ATTACHMENT"
    layout_policy = "STRICT"
    version_policy = "RELEASE"
  }
}

resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = ["%s"]
  }
  depends_on = [%s.member]
}
`, resourceTypeMavenHosted, memberName, resourceTypeMavenGroup, repoName, memberName, resourceTypeMavenHosted),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceMavenGroupName, "name", repoName),
					resource.TestCheckResourceAttr(resourceMavenGroupName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceMavenGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
func TestAccRepositoryMavenProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "maven-proxy-repo-%s"
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
  maven = {
    version_policy = "RELEASE"
    layout_policy = "STRICT"
  }
}
`, resourceTypeMavenProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageInvalidRemoteUrl),
			},
		},
	})
}

func TestAccRepositoryMavenHostedInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "maven-hosted-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  maven = {
    version_policy = "release"
    layout_policy = "permissive"
  }
}
`, resourceTypeMavenHosted, randomString),
				ExpectError: regexp.MustCompile(errorMessageBlobStoreNotFound),
			},
		},
	})
}

func TestAccRepositoryMavenHostedMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "maven-hosted-repo-%s"
  online = true
  # Missing storage block
}
`, resourceTypeMavenHosted, randomString),
				ExpectError: regexp.MustCompile(errorMessageStorageRequired),
			},
		},
	})
}

func TestAccRepositoryMavenProxyInvalidTimeoutTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too large, max is 3600)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "maven-proxy-repo-timeout-%s"
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
`, resourceTypeMavenProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryMavenProxyInvalidTimeoutTooSmall(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too small, min is 1)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "maven-proxy-repo-timeout-small-%s"
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
`, resourceTypeMavenProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
			},
		},
	})
}

func TestAccRepositoryMavenProxyInvalidRetriesTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (too large, max is 10)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "maven-proxy-repo-retries-%s"
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
`, resourceTypeMavenProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryMavenProxyInvalidRetriesNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "maven-proxy-repo-retries-neg-%s"
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
`, resourceTypeMavenProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
			},
		},
	})
}

func TestAccRepositoryMavenProxyInvalidTimeToLiveNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid time_to_live (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "maven-proxy-repo-ttl-%s"
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
`, resourceTypeMavenProxy, randomString),
				ExpectError: regexp.MustCompile(errorMessageNegativeCacheTimeoutValue),
			},
		},
	})
}
