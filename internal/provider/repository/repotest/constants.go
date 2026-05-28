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

// Package repotest provides shared test constants for repository acceptance tests.
package repotest

const (
	RES_ATTR_NAME                                             string = "name"
	RES_ATTR_ONLINE                                           string = "online"
	RES_ATTR_CLEANUP                                          string = "cleanup"
	RES_ATTR_CLEANUP_POLICY_COUNT                             string = "cleanup.policy_names.#"
	RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS                 string = "component.proprietary_components"
	RES_ATTR_DOCKER_FORCE_BASIC_AUTH                          string = "docker.force_basic_auth"
	RES_ATTR_DOCKER_PATH_ENABLED                              string = "docker.path_enabled"
	RES_ATTR_DOCKER_V1_ENABLED                                string = "docker.v1_enabled"
	RES_ATTR_DOCKER_PROXY_CACHE_FOREIGN_LAYERS                string = "docker_proxy.cache_foreign_layers"
	RES_ATTR_DOCKER_PROXY_INDEX_TYPE                          string = "docker_proxy.index_type"
	RES_ATTR_RAW_CONTENT_DISPOSITION                          string = "raw.content_disposition"
	RES_ATTR_STORAGE_BLOB_STORE_NAME                          string = "storage.blob_store_name"
	RES_ATTR_STORAGE_LATEST_POLICY                            string = "storage.latest_policy"
	RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION           string = "storage.strict_content_type_validation"
	RES_ATTR_STORAGE_WRITE_POLICY                             string = "storage.write_policy"
	RES_ATTR_URL                                              string = "url"
	RES_ATTR_PROXY_REMOTE_URL                                 string = "proxy.remote_url"
	RES_ATTR_PROXY_CONTENT_MAX_AGE                            string = "proxy.content_max_age"
	RES_ATTR_PROXY_METADATA_MAX_AGE                           string = "proxy.metadata_max_age"
	RES_ATTR_NEGATIVE_CACHE_ENABLED                           string = "negative_cache.enabled"
	RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE                      string = "negative_cache.time_to_live"
	RES_ATTR_HTTP_CLIENT_BLOCKED                              string = "http_client.blocked"
	RES_ATTR_HTTP_CLIENT_AUTO_BLOCK                           string = "http_client.auto_block"
	RES_ATTR_HTTP_CLIENT_AUTHENTICATION                       string = "http_client.authentication"
	RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PASSWORD              string = "http_client.authentication.password"
	RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PREMPTIVE             string = "http_client.authentication.preemptive"
	RES_ATTR_HTTP_CLIENT_AUTHENTICATION_TYPE                  string = "http_client.authentication.type"
	RES_ATTR_HTTP_CLIENT_AUTHENTICATION_USERNAME              string = "http_client.authentication.username"
	RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS string = "http_client.connection.enable_circular_redirects"
	RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES            string = "http_client.connection.enable_cookies"
	RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE           string = "http_client.connection.use_trust_store"
	RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES                   string = "http_client.connection.retries"
	RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT                   string = "http_client.connection.timeout"
	RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX         string = "http_client.connection.user_agent_suffix"
	RES_ATTR_GROUP_MEMBER_NAMES                               string = "group.member_names.#"
	RES_ATTR_REPLICATION                                      string = "replication"
	RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED             string = "replication.preemptive_pull_enabled"
	RES_ATTR_REPLICATION_ASSET_PATH_REGEX                     string = "replication.asset_path_regex"
	RES_ATTR_REPOSITORY_FIREWALL                              string = "repository_firewall"
	RES_ATTR_ROUTING_RULE_NAME                                string = "routing_rule"
	RES_ATTR_REPOSITORY_FIREWALL_ENABLED                      string = "repository_firewall.enabled"
	RES_ATTR_REPOSITORY_FIREWALL_QUARANTINE                   string = "repository_firewall.quarantine"
	RES_ATTR_APT_DISTRIBUTION                                 string = "apt.distribution"
	RES_ATTR_CARGO_REQUIRE_AUTHENTICATION                     string = "cargo.require_authentication"
	RES_ATTR_CONAN_PROXY_CONAN_VERSION                        string = "conan.conan_version"
	RES_ATTR_MAVEN_CONTENT_DISPOSITION                        string = "maven.content_disposition"
	RES_ATTR_MAVEN_LAYOUT_POLICY                              string = "maven.layout_policy"
	RES_ATTR_MAVEN_VERSION_POLICY                             string = "maven.version_policy"
	RES_ATTR_NUGET_PROXY_NUGET_VERSION                        string = "nuget_proxy.nuget_version"
	RES_ATTR_TERRAFORM_REQUIRE_AUTH                           string = "terraform.require_authentication"
	RES_ATTR_TERRAFORM_SIGNING_SIGNING_KEY                    string = "terraform_signing.signing_key"

	ConfigBlockHostedDefaultTerraform string = `terraform_signing = { signing_key = <<-EOT
-----BEGIN PGP PRIVATE KEY BLOCK-----

lQHYBGmLEQEBBADBLrTiM/XmBoTSBTdGRSMFgqM12vVi0+3K2vMk9Zd+HUN3O0zY
ho0Q16SQU9hY7eWRXp/XiyL59u7HQhtrBq36dthvZTCPh23G3ldCtlruPhQWtHI/
xO3phio8skST4MDRfS3csoyRc/rnY7Rc00P7J8HP7dx+sRqv+SnIBeyOLQARAQAB
AAP9E3Q4Z4IrjGlSJVM8pIEwXGzyMil1Ziko7HF9pFZuFddtFJv+alysZoqMyjMD
WbtFT80bZCmhEVKWa68C01WWHfK2CqPOsEFiWG/fxUbnUG7RlehMKrI6KF+2wWBv
o452loV/Bzua64uR1kP+l43BH69LzJE6uWHl5KNJyX1uoskCAMvp9kkzc1Pe2/hT
Vc72s/CkMlw6GMSI4Lk6+YuvajGlr/HxsFhBjM9ADLkWIDoywxCQ1kKSxtF/FG4a
zZG3GxkCAPKHAh6ByWSc+dfg1acQx1/LHaGdmLACaJYK7OAy4+ra+VrX6c3th6ye
T8SzJG2sq3aBztDBwdtjdWf+8BazwjUCALcVncOFj5a/N1vZZ6chuo27wEVw6Bpb
iN2rb+SxuT1iTaCE3/RfSywlqf0aVMkh3Dygz5/CwOwEmffNVGhV9gOpM7RpU29u
YXR5cGUgQ29tbXVuaXR5IFRlc3QgRGF0YSBLZXkgKFRFU1QgREFUQSBPTkxZIC0g
RE8gTk9UIFVTRSBGT1IgU0lHTklORykgPGNvbW11bml0eS1ncm91cEBzb25hdHlw
ZS5jb20+iM4EEwEIADgWIQTHqk34+snIQVhkIMl2y863mbo3rgUCaYsRAQIbAwUL
CQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRB2y863mbo3rqk8A/97anQKQ5ZUc/Oz
FSUpRI7sKCwot2C9dP7wAAifUtfX7vn32H4fz3T9BDB8CJMtVurMYVNskLvy5rKe
n//joo0cp+XX+KDVmuEPtxbZD5+Py+JGUwOuOkK8bO5N/xGe0N/CWfW2DXvpMvut
1x/uN9aslm9GBhezNr1V5totT7Tx5p0B2ARpixEBAQQAzmA6ajRv7U3tQV9aJav1
/y/+byUrm4hta2pGuo0qTeP8i2SEao3S7DkoENcA3MhRtxJk4fzMGmH68saQyKK5
66se0mZLTQjkPKGXje0pjAT4hTaffi3PDycR5MBEe80rbu08ouqSkKQ5xUrTL3FE
MC9BBbFMacl8EAbeIiIyvLEAEQEAAQAD+QGysjOuIRA2jpMrGj2NEXlHMJXYZuLu
PkoRR09TrU9pDB/skf1DXm2OmkCFDVsJBqjt2hYaLPF9YNnF0JqHAhENSJqShPgd
tyRprrObKuTYTH+cv4QI1cr52Oxr6BkAqP3VqPJxqrqXWLnscveryoxEPMlNvXbk
5ATexXyThwRBAgDlorsxh4YywRJdrQCSSnQiiYlqt/L2cKliTedbEN3ffFOD3OH/
zZbaoXruev75FrIZtgtfSgprLELw/fwsTxHBAgDmEd01R3S8R3tpkzMGvdCwZFL2
6uGaVmgTZ415XpXiNrDWSO0QeD15F6mnMwM8PsEEarRwrnc23pPZSLw2eIbxAf9D
o/bDNpno4GBpCd6P8sUhFRCw8UweU4EHVfz7OfnBkid8tvn4y85U2HJUi9jXj4/v
+yDRM+uhsch4VBac5xhVpwOItgQYAQgAIBYhBMeqTfj6ychBWGQgyXbLzreZujeu
BQJpixEBAhsMAAoJEHbLzreZujeu0o0D/iHzDEXpkHE/sbd82JwPR48YR8cmBzq+
CMnhPvAykyWgXvRoEQmXj+rzH9nlTsD9TNIVrnReTT5PBDGWVTk5DXkpb9ZfaajU
USlNrkrzgatlosKJskQSSrKSbmqju1/R885DZMtTb4ryjzHVqvwALzCFTXyyEpMl
OVIAhWNMezYE
=RbgK
-----END PGP PRIVATE KEY BLOCK-----
EOT
}`
)
