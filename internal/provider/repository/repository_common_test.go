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
	"terraform-provider-sonatyperepo/internal/provider/common"
)

const (
	RES_ATTR_NAME                                   string = "name"
	RES_ATTR_ONLINE                                 string = "online"
	RES_ATTR_CLEANUP                                string = "cleanup"
	RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS       string = "component.proprietary_components"
	RES_ATTR_DOCKER_FORCE_BASIC_AUTH                string = "docker.force_basic_auth"
	RES_ATTR_DOCKER_PATH_ENABLED                    string = "docker.path_enabled"
	RES_ATTR_DOCKER_V1_ENABLED                      string = "docker.v1_enabled"
	RES_ATTR_RAW_CONTENT_DISPOSITION                string = "raw.content_disposition"
	RES_ATTR_STORAGE_BLOB_STORE_NAME                string = "storage.blob_store_name"
	RES_ATTR_STORAGE_LATEST_POLICY                  string = "storage.latest_policy"
	RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION string = "storage.strict_content_type_validation"
	RES_ATTR_STORAGE_WRITE_POLICY                   string = "storage.write_policy"
	RES_ATTR_URL                                    string = "url"
	RES_ATTR_PROXY_REMOTE_URL                       string = "proxy.remote_url"
	RES_ATTR_PROXY_CONTENT_MAX_AGE                  string = "proxy.content_max_age"
	RES_ATTR_PROXY_METADATA_MAX_AGE                 string = "proxy.metadata_max_age"
	RES_ATTR_NEGATIVE_CACHE_ENABLED                 string = "negative_cache.enabled"
	RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE            string = "negative_cache.time_to_live"
	RES_ATTR_HTTP_CLIENT_BLOCKED                    string = "http_client.blocked"
	RES_ATTR_HTTP_CLIENT_AUTO_BLOCK                 string = "http_client.auto_block"
	RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS string = "http_client.connection.enable_circular_redirects"
	RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES  string = "http_client.connection.enable_cookies"
	RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE string = "http_client.connection.use_trust_store"
	RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES         string = "http_client.connection.retries"
	RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT         string = "http_client.connection.timeout"
	RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX string = "http_client.connection.user_agent_suffix"
	RES_ATTR_GROUP_MEMBER_NAMES                     string = "group.member_names.#"
)

var (
	errorMessageBlobStoreNotFound                = "Blob store.*not found"
	errorMessageGroupMemberNamesEmpty            = "Attribute group.member_names list must contain at least 1 elements"
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
