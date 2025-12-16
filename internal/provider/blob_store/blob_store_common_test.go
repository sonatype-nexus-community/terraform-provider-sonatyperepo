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
	errMessageBlobStoreGroupNoMembers        string = "cannot be empty"
	errMessageBlobStoreGroupIneligibleMember string = "is not eligible to be a group member"

	RES_ATTR_NAME             string = "name"
	RES_ATTR_PATH             string = "path"
	RES_ATTR_SOFT_QUOTA       string = "soft_quota"
	RES_ATTR_SOFT_QUOTA_LIMIT string = "soft_quota.limit"
	RES_ATTR_SOFT_QUOTA_TYPE  string = "soft_quota.type"
	RES_ATTR_FILL_POLICY      string = "fill_policy"
	RES_ATTR_MEMBERS_COUNT    string = "members.#"
	RES_ATTR_LAST_UPDATED     string = "last_updated"

	RES_NAME_FMT string = "%s.test"

	RES_TYPE_BLOB_STORE_FILE  string = "sonatyperepo_blob_store_file"
	RES_TYPE_BLOB_STORE_GCS   string = "sonatyperepo_blob_store_gcs"
	RES_TYPE_BLOB_STORE_GROUP string = "sonatyperepo_blob_store_group"
	RES_TYPE_BLOB_STORE_S3    string = "sonatyperepo_blob_store_s3"
)

var (
	RES_NAME_BLOB_STORE_FILE  string = fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_FILE)
	RES_NAME_BLOB_STORE_GCS   string = fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_GCS)
	RES_NAME_BLOB_STORE_GROUP string = fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_GROUP)
	RES_NAME_BLOB_STORE_S3    string = fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_S3)
)
