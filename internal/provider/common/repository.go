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

package common

const (
	CONAN_PROTOCOL_V1                                        string = "V1"
	CONAN_PROTOCOL_V2                                        string = "V2"
	CONTENT_DISPOSITION_ATTACHMENT                           string = "ATTACHMENT"
	CONTENT_DISPOSITION_INLINE                               string = "INLINE"
	DEFAULT_PROXY_CONTENT_MAX_AGE                            int64  = 1440
	DEFAULT_PROXY_METADATA_MAX_AGE                           int64  = 1440
	DEFAULT_PROXY_NEGATIVE_CACHE_ENABLED                     bool   = true
	DEFAULT_PROXY_NEGATIVE_CACHE_TTL                         int64  = 1440
	DEFAULT_PROXY_PREEMPTIVE_PULL                            bool   = false
	DEFAULT_HTTP_CLIENT_AUTO_BLOCK                           bool   = true
	DEFAULT_HTTP_CLIENT_BLOCKED                              bool   = false
	DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS bool   = false
	DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES            bool   = false
	DEFAULT_HTTP_CLIENT_CONNECTION_RETRIES                   int64  = 0
	DEFAULT_HTTP_CLIENT_CONNECTION_TIMEOUT                   int64  = 60
	DEFAULT_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE           bool   = false
	DEPLOY_POLICY_PERMISSIVE                                 string = "PERMISSIVE"
	DEPLOY_POLICY_STRICT                                     string = "STRICT"
	DOCKER_PROXY_INDEX_TYPE_HUB                              string = "HUB"
	DOCKER_PROXY_INDEX_TYPE_REGISTRY                         string = "REGISTRY"
	DOCKER_PROXY_INDEX_TYPE_CUSTOM                           string = "CUSTOM"
	MAVEN_CONTENT_DISPOSITION_ATTACHMENT                     string = "ATTACHMENT"
	MAVEN_CONTENT_DISPOSITION_INLINE                         string = "INLINE"
	MAVEN_LAYOUT_STRICT                                      string = "STRICT"
	MAVEN_LAYOUT_PERMISSIVE                                  string = "PERMISSIVE"
	MAVEN_VERSION_POLICY_RELEASE                             string = "RELEASE"
	MAVEN_VERSION_POLICY_SNAPSHOT                            string = "SNAPSHOT"
	MAVEN_VERSION_POLICY_MIXED                               string = "MIXED"
	NUGET_PROTOCOL_V2                                        string = "V2"
	NUGET_PROTOCOL_V3                                        string = "V3"
	REPO_FORMAT_APT                                          string = "APT"
	REPO_FORMAT_CARGO                                        string = "CARGO"
	REPO_FORMAT_COCOAPODS                                    string = "COCOAPODS"
	REPO_FORMAT_COMPOSER                                     string = "COMPOSER"
	REPO_FORMAT_CONAN                                        string = "CONAN"
	REPO_FORMAT_CONDA                                        string = "CONDA"
	REPO_FORMAT_DOCKER                                       string = "DOCKER"
	REPO_FORMAT_GIT_LFS                                      string = "GITLFS"
	REPO_FORMAT_GO                                           string = "GO"
	REPO_FORMAT_HELM                                         string = "HELM"
	REPO_FORMAT_HUGGING_FACE                                 string = "HUGGINGFACE"
	REPO_FORMAT_MAVEN                                        string = "MAVEN2"
	REPO_FORMAT_NPM                                          string = "NPM"
	REPO_FORMAT_NUGET                                        string = "NUGET"
	REPO_FORMAT_P2                                           string = "P2"
	REPO_FORMAT_PYPI                                         string = "PYPI"
	REPO_FORMAT_RAW                                          string = "RAW"
	REPO_FORMAT_R                                            string = "R"
	REPO_FORMAT_RUBY_GEMS                                    string = "RUBYGEMS"
	REPO_FORMAT_TERRAFORM                                    string = "TERRAFORM"
	REPO_FORMAT_YUM                                          string = "YUM"
	REPOSITORY_HTTP_CLIENT_CONNECTION_RETRIES_MIN            int64  = 1
	REPOSITORY_HTTP_CLIENT_CONNECTION_RETRIES_MAX            int64  = 10
	REPOSITORY_HTTP_CLIENT_CONNECTION_TIMEOUT_MIN            int64  = 1
	REPOSITORY_HTTP_CLIENT_CONNECTION_TIMEOUT_MAX            int64  = 3600
	WRITE_POLICY_ALLOW                                       string = "ALLOW"
	WRITE_POLICY_ALLOW_ONCE                                  string = "ALLOW_ONCE"
	WRITE_POLICY_DENY                                        string = "DENY"
)

func AllHostedFormats() []string {
	return []string{
		REPO_FORMAT_APT,
		REPO_FORMAT_CARGO,
		REPO_FORMAT_CONAN,
		REPO_FORMAT_DOCKER,
		REPO_FORMAT_GIT_LFS,
		REPO_FORMAT_HELM,
		REPO_FORMAT_MAVEN,
		REPO_FORMAT_NPM,
		REPO_FORMAT_NUGET,
		REPO_FORMAT_PYPI,
		REPO_FORMAT_R,
		REPO_FORMAT_RAW,
		REPO_FORMAT_RUBY_GEMS,
		REPO_FORMAT_YUM,
	}
}

func AllProxyFormats() []string {
	return []string{
		REPO_FORMAT_APT,
		REPO_FORMAT_CARGO,
		REPO_FORMAT_COCOAPODS,
		REPO_FORMAT_COMPOSER,
		REPO_FORMAT_CONAN,
		REPO_FORMAT_CONDA,
		REPO_FORMAT_DOCKER,
		REPO_FORMAT_GO,
		REPO_FORMAT_HELM,
		REPO_FORMAT_HUGGING_FACE,
		REPO_FORMAT_MAVEN,
		REPO_FORMAT_NPM,
		REPO_FORMAT_NUGET,
		REPO_FORMAT_P2,
		REPO_FORMAT_PYPI,
		REPO_FORMAT_R,
		REPO_FORMAT_RAW,
		REPO_FORMAT_RUBY_GEMS,
		REPO_FORMAT_YUM,
	}
}
