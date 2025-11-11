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
	AUTH_SCHEME_NONE                                                       string = "NONE"
	AUTH_SCHEME_SIMPLE                                                     string = "SIMPLE"
	AUTH_SCHEME_DIGEST_MD5                                                 string = "DIGEST_MD5"
	AUTH_SCHEME_CRAM_MD5                                                   string = "CRAM_MD5"
	CAPABILITY_FIREWALL_AUDIT_QUARANTINE_DEFAULT_QUARANTINE                bool   = false
	CAPABILITY_HEALTHCHECK_DEFAULT_CONFIGURED_FOR_ALL                      bool   = true
	CAPABILITY_HEALTHCHECK_DEFAULT_USE_NEXUS_TRUSTSTORE                    bool   = false
	CAPABILITY_OUTREACH_DEFAULT_ALWAYS_REMOTE                              bool   = false
	CAPABILITY_STORAGE_SETTINGS_DEFAULT_LAST_DOWNLOADED_INTERVAL           int32  = 12
	CAPABILITY_UI_BRANDING_DEFAULT_FOOTER_ENABLED                          bool   = false
	CAPABILITY_UI_BRANDING_DEFAULT_HEADER_ENABLED                          bool   = false
	CAPABILITY_UI_BRANDING_DEFAULT_HEADER_HTML                             string = ""
	CAPABILITY_UI_BRANDING_DEFAULT_FOOTER_HTML                             string = ""
	CAPABILITY_UI_SETTINGS_DEFAULT_DEBUG_ALLOWED                           bool   = false
	CAPABILITY_UI_SETTINGS_DEFAULT_LONG_REQUEST_TIMEOUT                    int32  = 180
	CAPABILITY_UI_SETTINGS_DEFAULT_REQUEST_TIMEOUT                         int32  = 60
	CAPABILITY_UI_SETTINGS_DEFAULT_SESSION_TIMEOUT                         int32  = 30
	CAPABILITY_UI_SETTINGS_DEFAULT_STATUS_INTERVAL_ANONYMOUS               int32  = 60
	CAPABILITY_UI_SETTINGS_DEFAULT_STATUS_INTERVAL_AUTHENTICATED           int32  = 5
	CAPABILITY_UI_SETTINGS_DEFAULT_TITLE                                   string = "Sonatype Nexus Repository"
	CONAN_PROTOCOL_V1                                                      string = "V1"
	CONAN_PROTOCOL_V2                                                      string = "V2"
	CONTENT_DISPOSITION_ATTACHMENT                                         string = "ATTACHMENT"
	CONTENT_DISPOSITION_INLINE                                             string = "INLINE"
	DEFAULT_ANONYMOUS_USERNAME                                             string = "anonymous"
	DEFAULT_BLOB_STORE_NAME                                                string = "default"
	DEFAULT_HTTP_CONNECTION_RETRIES                                        int64  = 0
	DEFAULT_HTTP_CONNECTION_TIMEOUT                                        int64  = 60
	DEFAULT_REALM_NAME                                                     string = "NexusAuthorizingRealm"
	DEFAULT_USER_SOURCE                                                    string = "default"
	DEPLOY_POLICY_PERMISSIVE                                               string = "PERMISSIVE"
	DEPLOY_POLICY_STRICT                                                   string = "STRICT"
	DOCKER_PROXY_INDEX_TYPE_HUB                                            string = "HUB"
	DOCKER_PROXY_INDEX_TYPE_REGISTRY                                       string = "REGISTRY"
	DOCKER_PROXY_INDEX_TYPE_CUSTOM                                         string = "CUSTOM"
	FREQUENCY_SCHEDULE_MANUAL                                              string = "manual"
	FREQUENCY_SCHEDULE_ONCE                                                string = "once"
	FREQUENCY_SCHEDULE_HOURLY                                              string = "hourly"
	FREQUENCY_SCHEDULE_DAILY                                               string = "daily"
	FREQUENCY_SCHEDULE_WEEKLY                                              string = "weekly"
	FREQUENCY_SCHEDULE_MONTHLY                                             string = "monthly"
	FREQUENCY_SCHEDULE_CRON                                                string = "cron"
	HTTP_AUTH_TYPE_BEARER_TOKEN                                            string = "bearerToken"
	HTTP_AUTH_TYPE_NTLM                                                    string = "ntlm"
	HTTP_AUTH_TYPE_USERNAME                                                string = "username"
	IQ_AUTHENTICATON_TYPE_PKI                                              string = "PKI"
	IQ_AUTHENTICATON_TYPE_USER                                             string = "USER"
	LDAP_GROUP_MAPPING_DYNAMIC                                             string = "DYNAMIC"
	LDAP_GROUP_MAPPING_STATIC                                              string = "STATIC"
	MAVEN_CONTENT_DISPOSITION_ATTACHMENT                                   string = "ATTACHMENT"
	MAVEN_CONTENT_DISPOSITION_INLINE                                       string = "INLINE"
	MAVEN_LAYOUT_STRICT                                                    string = "STRICT"
	MAVEN_LAYOUT_PERMISSIVE                                                string = "PERMISSIVE"
	MAVEN_VERSION_POLICY_RELEASE                                           string = "RELEASE"
	MAVEN_VERSION_POLICY_SNAPSHOT                                          string = "SNAPSHOT"
	MAVEN_VERSION_POLICY_MIXED                                             string = "MIXED"
	NOTIFICATION_CONDITION_FAILURE                                         string = "FAILURE"
	NOTIFICATION_CONDITION_SUCCESS_OR_FAILURE                              string = "SUCCESS_FAILURE"
	NUGET_PROTOCOL_V2                                                      string = "V2"
	NUGET_PROTOCOL_V3                                                      string = "V3"
	PLACEHOLDER_PASSWORD                                                   string = "#~NXRM~PLACEHOLDER~PASSWORD~#"
	PROTOCOL_LDAP                                                          string = "LDAP"
	PROTOCOL_LDAPS                                                         string = "LDAPS"
	REPO_FORMAT_APT                                                        string = "APT"
	REPO_FORMAT_CARGO                                                      string = "CARGO"
	REPO_FORMAT_COCOAPODS                                                  string = "COCOAPODS"
	REPO_FORMAT_COMPOSER                                                   string = "COMPOSER"
	REPO_FORMAT_CONAN                                                      string = "CONAN"
	REPO_FORMAT_CONDA                                                      string = "CONDA"
	REPO_FORMAT_DOCKER                                                     string = "DOCKER"
	REPO_FORMAT_GIT_LFS                                                    string = "GITLFS"
	REPO_FORMAT_GO                                                         string = "GO"
	REPO_FORMAT_HELM                                                       string = "HELM"
	REPO_FORMAT_HUGGING_FACE                                               string = "HUGGINGFACE"
	REPO_FORMAT_MAVEN                                                      string = "MAVEN"
	REPO_FORMAT_NPM                                                        string = "NPM"
	REPO_FORMAT_NUGET                                                      string = "NUGET"
	REPO_FORMAT_P2                                                         string = "P2"
	REPO_FORMAT_PYPI                                                       string = "PYPI"
	REPO_FORMAT_RAW                                                        string = "RAW"
	REPO_FORMAT_R                                                          string = "R"
	REPO_FORMAT_RUBY_GEMS                                                  string = "RUBY_GEMS"
	REPO_FORMAT_YUM                                                        string = "YUM"
	TASK_REPOSITORY_DOCKER_GC_DEFAULT_DEPLOY_OFFSET                        int32  = 24
	TASK_REPOSITORY_DOCKER_UPLOAD_PURGE_DEFAULT_AGE                        int32  = 24
	TASK_REPOSITORY_MAVEN_REMOVE_SNAPSHOTS_DEFAULT_MINIMUM_RETAINED        int32  = 1
	TASK_REPOSITORY_MAVEN_REMOVE_SNAPSHOTS_DEFAULT_REMOVE_IF_RELEASED      bool   = false
	TASK_REPOSITORY_MAVEN_REMOVE_SNAPSHOTS_DEFAULT_SNAPSHOT_RETENTION_DAYS int32  = 30
	USER_STATUS_ACTIVE                                                     string = "active"
	USER_STATUS_LOCKED                                                     string = "locked"
	USER_STATUS_DISABLED                                                   string = "disabled"
	USER_STATUS_CHANGE_PASSWORD                                            string = "changepassword"
	WRITE_POLICY_ALLOW                                                     string = "ALLOW"
	WRITE_POLICY_ALLOW_ONCE                                                string = "ALLOW_ONCE"
	WRITE_POLICY_DENY                                                      string = "DENY"
	SECURITY_USER_TOKEN_DEFAULT_ENABLED                                    bool   = false
	SECURITY_USER_TOKEN_DEFAULT_EXPIRATION_DAYS                            int32  = 1
	SECURITY_USER_TOKEN_DEFAULT_EXPIRATION_ENABLED                         bool   = false
	SECURITY_USER_TOKEN_DEFAULT_PROTECT_CONTENT                            bool   = false
	ERROR_MESSAGE_UNAUTHORIZED                                             string = "Your user is unauthorized to access this resource or feature."
)

func NewFalse() *bool {
	b := false
	return &b
}

func NewTrue() *bool {
	b := true
	return &b
}

func StringPointer(s string) *string { return &s }
