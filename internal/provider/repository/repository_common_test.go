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
	"strings"
	"testing"

	"terraform-provider-sonatyperepo/internal/provider/common"
	repotest "terraform-provider-sonatyperepo/internal/provider/repository/repotest"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	repoNameFString string = "%s.repo"

	hostedNameFString         string = "test-%s-hosted-repo-%s"
	proxyNameFString          string = "test-%s-proxy-repo-%s"
	resourceTypeHostedFString string = "sonatyperepo_repository_%s_hosted"
	resourceTypeProxyFString  string = "sonatyperepo_repository_%s_proxy"

	TEST_DATA_APT_PROXY_REMOTE_URL          string = "https://archive.ubuntu.com/ubuntu/"
	TEST_DATA_CARGO_PROXY_REMOTE_URL        string = "https://index.crates.io/"
	TEST_DATA_COCOAPODS_PROXY_REMOTE_URL    string = "https://cdn.cocoapods.org/"
	TEST_DATA_COMPOSER_PROXY_REMOTE_URL     string = "https://packagist.org/"
	TEST_DATA_CONAN_PROXY_REMOTE_URL        string = "https://center2.conan.io"
	TEST_DATA_CONDA_PROXY_REMOTE_URL        string = "https://repo.anaconda.com/pkgs/"
	TEST_DATA_DOCKER_PROXY_REMOTE_URL       string = "https://registry-1.docker.io"
	TEST_DATA_GO_PROXY_REMOTE_URL           string = "https://proxy.golang.org"
	TEST_DATA_HELM_PROXY_REMOTE_URL         string = "https://charts.helm.sh/stable"
	TEST_DATA_HUGGING_FACE_PROXY_REMOTE_URL string = "https://huggingface.co"
	TEST_DATA_MAVEN_PROXY_REMOTE_URL        string = "https://repo.maven.apache.org/maven2"
	TEST_DATA_NPM_PROXY_REMOTE_URL          string = "https://registry.npmjs.org"
	TEST_DATA_NUGET_PROXY_REMOTE_URL        string = "https://api.nuget.org/v3/index.json"
	TEST_DATA_P2_PROXY_REMOTE_URL           string = "https://download.eclipse.org/releases/2025-06"
	TEST_DATA_PYPI_PROXY_REMOTE_URL         string = "https://pypi.org/simple"
	TEST_DATA_R_PROXY_REMOTE_URL            string = "https://cran.r-project.org/"
	TEST_DATA_RAW_PROXY_REMOTE_URL          string = "https://nodejs.org/dist/"
	TEST_DATA_RUBY_GEMS_PROXY_REMOTE_URL    string = "https://rubygems.org"
	TEST_DATA_SWIFT_PROXY_REMOTE_URL        string = "https://github.com/"
	TEST_DATA_TERRAFORM_PROXY_REMOTE_URL    string = "https://registry.terraform.io"
	TEST_DATA_YUM_PROXY_REMOTE_URL          string = "https://mirror.centos.org/centos/"
	TEST_DATA_TIMEOUT                       string = "1439"
)

var (
	errorMessageBlobStoreNotFound                = "Blob store.*not found"
errorMessageInvalidRemoteUrl                 = "Attribute proxy.remote_url must be a valid HTTP URL"
	errorMessageHttpClientConnectionRetriesValue = fmt.Sprintf(
		"Attribute http_client.connection.retries value must be between %d and %d",
		common.REPOSITORY_HTTP_CLIENT_CONNECTION_RETRIES_MIN,
		common.REPOSITORY_HTTP_CLIENT_CONNECTION_RETRIES_MAX,
	)
	errorMessageHttpClientConnectionTimeoutValue = fmt.Sprintf(
		"Attribute http_client.connection.timeout value must be between %d and %d",
		common.REPOSITORY_HTTP_CLIENT_CONNECTION_TIMEOUT_MIN,
		common.REPOSITORY_HTTP_CLIENT_CONNECTION_TIMEOUT_MAX,
	)
	errorMessageNegativeCacheTimeoutValue = "Attribute negative_cache.time_to_live value must be at least 0"
	errorMessageStorageRequired           = "The argument \"storage\" is required, but no definition was found."
)

type repositoryHostedTestData struct {
	CheckFunc                     func(resourceName string) []resource.TestCheckFunc
	RepoFormat                    string
	SchemaFunc                    func(resourceType, repoName, repoFormat, randomString string, completeData, supportsProprietaryComponents bool) string
	SupportsProprietaryComponents bool
	TestImport                    bool
	TestPreCheck                  func(t *testing.T) func()
}

type repositoryProxyTestData struct {
	CheckFunc            func(resourceName string) []resource.TestCheckFunc
	FormatSpecificConfig string
	RemoteUrl            string
	RepoFormat           string
	SchemaFunc           func(resourceType, repoName, repoFormat, remoteUrl, randomString, formatSpecificConfig string, completeData bool) string
	TestImport           bool
	TestPreCheck         func(t *testing.T) func()
}

func allRepositoryFormatsGroupGeneric() []string {
	return []string{
		strings.ToLower(common.REPO_FORMAT_CONAN),
		// strings.ToLower(common.REPO_FORMAT_GO),		// No hosted repo
		// strings.ToLower(common.REPO_FORMAT_MAVEN),	// Hosted repo requires specific config
		strings.ToLower(common.REPO_FORMAT_NPM),
		strings.ToLower(common.REPO_FORMAT_NUGET),
		strings.ToLower(common.REPO_FORMAT_PYPI),
		strings.ToLower(common.REPO_FORMAT_R),
		strings.ToLower(common.REPO_FORMAT_RUBY_GEMS),
		// strings.ToLower(common.REPO_FORMAT_YUM),		// Hosted repo requires specific config
	}
}

// ------------------------------------------------------------
// GROUP REPO TESTING (GENERIC)
// ------------------------------------------------------------
func TestAccRepositoryGenericGroupResources(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	for _, repoFormat := range allRepositoryFormatsGroupGeneric() {
		resourceTypeGroup := fmt.Sprintf("sonatyperepo_repository_%s_group", repoFormat)
		resourceTypeHosted := fmt.Sprintf("sonatyperepo_repository_%s_hosted", repoFormat)
		reosurceNameGroup := fmt.Sprintf(repoNameFString, resourceTypeGroup)
		reosurceNameHosted := fmt.Sprintf(repoNameFString, resourceTypeHosted)
		repoNameGroup := fmt.Sprintf("%s-group-repo-%s", repoFormat, randomString)
		repoNameHosted := fmt.Sprintf("%s-hosted-repo-%s", repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: repositoryGroupResourceNoCleanupNoProprietaryConfig(
						resourceTypeGroup,
						resourceTypeHosted,
						repoNameGroup,
						repoNameHosted,
						common.WRITE_POLICY_ALLOW_ONCE,
						false,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify Hosted
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_NAME, repoNameHosted),
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_ONLINE, "false"),
						resource.TestCheckResourceAttrSet(reosurceNameHosted, repotest.RES_ATTR_URL),
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
						resource.TestCheckNoResourceAttr(reosurceNameHosted, repotest.RES_ATTR_CLEANUP),

						// Verify Group
						resource.TestCheckResourceAttr(reosurceNameGroup, repotest.RES_ATTR_NAME, repoNameGroup),
						resource.TestCheckResourceAttr(reosurceNameGroup, repotest.RES_ATTR_ONLINE, "false"),
						resource.TestCheckResourceAttrSet(reosurceNameGroup, repotest.RES_ATTR_URL),
						resource.TestCheckResourceAttr(reosurceNameGroup, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(reosurceNameGroup, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						resource.TestCheckResourceAttr(reosurceNameGroup, repotest.RES_ATTR_GROUP_MEMBER_NAMES, "1"),
					),
				},
				// Update 1 - Set Online
				{
					Config: repositoryGroupResourceNoCleanupNoProprietaryConfig(
						resourceTypeGroup,
						resourceTypeHosted,
						repoNameGroup,
						repoNameHosted,
						common.WRITE_POLICY_ALLOW_ONCE,
						true,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify Hosted
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_NAME, repoNameHosted),
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_ONLINE, "true"),
						resource.TestCheckResourceAttrSet(reosurceNameHosted, repotest.RES_ATTR_URL),
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
						resource.TestCheckResourceAttr(reosurceNameHosted, repotest.RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
						resource.TestCheckNoResourceAttr(reosurceNameHosted, repotest.RES_ATTR_CLEANUP),

						// Verify Group
						resource.TestCheckResourceAttr(reosurceNameGroup, repotest.RES_ATTR_NAME, repoNameGroup),
						resource.TestCheckResourceAttr(reosurceNameGroup, repotest.RES_ATTR_ONLINE, "true"),
						resource.TestCheckResourceAttrSet(reosurceNameGroup, repotest.RES_ATTR_URL),
						resource.TestCheckResourceAttr(reosurceNameGroup, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(reosurceNameGroup, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						resource.TestCheckResourceAttr(reosurceNameGroup, repotest.RES_ATTR_GROUP_MEMBER_NAMES, "1"),
					),
				},
			},
		})
	}
}

func repositoryGroupResourceNoCleanupNoProprietaryConfig(resourceTypeGroup, resourceTypeHosted, repoNameGroup, repoNameHosted, writePolicy string, repoOnline bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
	name = "%s"
	online = %t
	storage = {
		blob_store_name = "default"
		strict_content_type_validation = true
		write_policy = "%s"
	}
}

resource "%s" "repo" {
  name = "%s"
  online = %t
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeHosted, repoNameHosted, repoOnline, writePolicy, resourceTypeGroup, repoNameGroup, repoOnline, repoNameHosted, resourceTypeHosted)
}

// ------------------------------------------------------------
// PROXY REPO TESTING (GENERIC)
// ------------------------------------------------------------
