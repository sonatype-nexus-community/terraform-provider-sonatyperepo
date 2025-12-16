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

package blob_store_test

import "fmt"

const (
	errMessageBlobStoreGroupNoMembers        string = "Blob Store '.*' cannot be empty"
	errMessageBlobStoreGroupIneligibleMember string = "Blob Store '.*' is set as storage for .* repositories and is not eligible to be a group member"

	RES_ATTR_NAME          string = "name"
	RES_ATTR_PATH          string = "path"
	RES_ATTR_SOFT_QUOTA    string = "soft_quota"
	RES_ATTR_FILL_POLICY   string = "fill_policy"
	RES_ATTR_MEMBERS_COUNT string = "members.#"
	RES_ATTR_LAST_UPDATED  string = "last_updated"

	RES_TYPE_BLOB_STORE_FILE  string = "sonatyperepo_blob_store_file"
	RES_TYPE_BLOB_STORE_GROUP string = "sonatyperepo_blob_store_group"
)

var (
	RES_NAME_BLOB_STORE_GROUP string = fmt.Sprintf("%s.test", RES_TYPE_BLOB_STORE_GROUP)
)
