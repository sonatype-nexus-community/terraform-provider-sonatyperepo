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
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeMavenGroup  = "sonatyperepo_repository_maven2_group"
	resourceTypeMavenHosted = "sonatyperepo_repository_maven2_hosted"
	resourceTypeMavenProxy  = "sonatyperepo_repository_maven2_proxy"
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
		PreCheck: func() {
			// Skip these tests for NXRM 3.88.x due to https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/268
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 88,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 88,
				Patch: 99,
			})
		},
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
					resource.TestCheckResourceAttr(resourceMavenHostedName, RES_ATTR_NAME, fmt.Sprintf("maven-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceMavenHostedName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceMavenHostedName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceMavenHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceMavenHostedName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceMavenHostedName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceMavenHostedName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
					resource.TestCheckNoResourceAttr(resourceMavenHostedName, RES_ATTR_CLEANUP),
					resource.TestCheckResourceAttr(resourceMavenHostedName, RES_ATTR_MAVEN_CONTENT_DISPOSITION, common.MAVEN_CONTENT_DISPOSITION_ATTACHMENT),
					resource.TestCheckResourceAttr(resourceMavenHostedName, RES_ATTR_MAVEN_LAYOUT_POLICY, common.MAVEN_LAYOUT_STRICT),
					resource.TestCheckResourceAttr(resourceMavenHostedName, RES_ATTR_MAVEN_VERSION_POLICY, common.MAVEN_VERSION_POLICY_RELEASE),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_NAME, fmt.Sprintf("maven-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceMavenProxyName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_PROXY_REMOTE_URL, "https://repo1.maven.org/maven2/"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_PROXY_CONTENT_MAX_AGE, "1442"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_PROXY_METADATA_MAX_AGE, "1400"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_NEGATIVE_CACHE_ENABLED, "true"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, "1440"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_HTTP_CLIENT_BLOCKED, "false"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, "true"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, "false"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, "true"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, "true"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, "9"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, "999"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX, "terraform"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION_USERNAME, "user"),
					resource.TestCheckResourceAttr(resourceMavenProxyName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PASSWORD, "pass"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceMavenGroupName, RES_ATTR_NAME, fmt.Sprintf("maven-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceMavenGroupName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceMavenGroupName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceMavenGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceMavenGroupName, RES_ATTR_GROUP_MEMBER_NAMES, "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryMavenGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("maven-group-import-%s", randomString)
	memberName := fmt.Sprintf("maven-hosted-member-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Skip these tests for NXRM 3.88.x due to https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/268
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 88,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 88,
				Patch: 99,
			})
		},
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
					resource.TestCheckResourceAttr(resourceMavenGroupName, RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceMavenGroupName, RES_ATTR_ONLINE, "true"),
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
