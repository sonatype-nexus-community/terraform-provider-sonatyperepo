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
	RES_ATTR_DOCKER_FORCE_BASIC_AUTH string = "docker.force_basic_auth"
	RES_ATTR_DOCKER_PATH_ENABLED     string = "docker.path_enabled"
	RES_ATTR_DOCKER_V1_ENABLED       string = "docker.v1_enabled"
	RES_ATTR_RAW_CONTENT_DISPOSITION string = "raw.content_disposition"
	RES_ATTR_STORAGE_BLOB_STORE_NAME string = "storage.blob_store_name"
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
