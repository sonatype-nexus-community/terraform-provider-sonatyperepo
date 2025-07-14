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
	AUTH_SCHEME_NONE                            = "NONE"
	AUTH_SCHEME_SIMPLE                          = "SIMPLE"
	AUTH_SCHEME_DIGEST_MD5                      = "DIGEST_MD5"
	AUTH_SCHEME_CRAM_MD5                        = "CRAM_MD5"
	DEFAULT_ANONYMOUS_USERNAME           string = "anonymous"
	DEFAULT_REALM_NAME                   string = "NexusAuthorizingRealm"
	DEFAULT_USER_SOURCE                  string = "default"
	HTTP_AUTH_TYPE_BEARER_TOKEN          string = "bearerToken"
	HTTP_AUTH_TYPE_NTLM                  string = "ntlm"
	HTTP_AUTH_TYPE_USERNAME              string = "username"
	IQ_AUTHENTICATON_TYPE_PKI            string = "PKI"
	IQ_AUTHENTICATON_TYPE_USER           string = "USER"
	LDAP_GROUP_MAPPING_DYNAMIC           string = "DYNAMIC"
	LDAP_GROUP_MAPPING_STATIC            string = "STATIC"
	MAVEN_CONTENT_DISPOSITION_ATTACHMENT string = "ATTACHMENT"
	MAVEN_CONTENT_DISPOSITION_INLINE     string = "INLINE"
	MAVEN_LAYOUT_STRICT                  string = "STRICT"
	MAVEN_LAYOUT_PERMISSIVE              string = "PERMISSIVE"
	MAVEN_VERSION_POLICY_RELEASE         string = "RELEASE"
	MAVEN_VERSION_POLICY_SNAPSHOT        string = "SNAPSHOT"
	MAVEN_VERSION_POLICY_MIXED           string = "MIXED"
	PLACEHOLDER_PASSWORD                 string = "#~NXRM~PLACEHOLDER~PASSWORD~#"
	PROTOCOL_LDAP                        string = "LDAP"
	PROTOCOL_LDAPS                       string = "LDAPS"
	REPO_FORMAT_APT                      string = "APT"
	REPO_FORMAT_DOCKER                   string = "DOCKER"
	REPO_FORMAT_MAVEN                    string = "MAVEN"
	REPO_FORMAT_NPM                      string = "NPM"
	USER_STATUS_ACTIVE                   string = "active"
	USER_STATUS_LOCKED                   string = "locked"
	USER_STATUS_DISABLED                 string = "disabled"
	USER_STATUS_CHANGE_PASSWORD          string = "changepassword"
	WRITE_POLICY_ALLOW                   string = "ALLOW"
	WRITE_POLICY_ALLOW_ONCE              string = "ALLOW_ONCE"
	WRITE_POLICY_DENY                    string = "DENY"
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
