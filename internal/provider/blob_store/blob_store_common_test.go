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
	errMessageBlobStoreAcsErrorCreating      string = "Error creating Azure Cloud Storage Blob Store"
	errMessageBlobStoreGroupNoMembers        string = "cannot be empty"
	errMessageBlobStoreGroupIneligibleMember string = "is not eligible to be a group"
	errMessageBlobStoreS3ErrorCreating       string = "Error creating S3 Blob Store|InvalidAccessKeyId|NoSuchBucket"

	awsRegionEuWest2 string = "eu-west-2"

	testS3NamePresigned       string = "test-s3-presigned-%s"
	testS3BucketPrefix        string = "prefix-%s"
	testS3BucketNamePresigned string = "nexus-bucket-presigned-%s"

	RES_ATTR_NAME             string = "name"
	RES_ATTR_PATH             string = "path"
	RES_ATTR_SOFT_QUOTA       string = "soft_quota"
	RES_ATTR_SOFT_QUOTA_LIMIT string = "soft_quota.limit"
	RES_ATTR_SOFT_QUOTA_TYPE  string = "soft_quota.type"
	RES_ATTR_FILL_POLICY      string = "fill_policy"
	RES_ATTR_MEMBERS_COUNT    string = "members.#"
	RES_ATTR_LAST_UPDATED     string = "last_updated"

	RES_ATTR_ACS_BUCKET_CONFIGURATION_ACCOUNT_NAME                     string = "bucket_configuration.account_name"
	RES_ATTR_ACS_BUCKET_CONFIGURATION_CONTAINER_NAME                   string = "bucket_configuration.container_name"
	RES_ATTR_ACS_BUCKET_CONFIGURATION_AUTH_AUTHENTICATION_METHOD       string = "bucket_configuration.authentication.authentication_method"
	RES_ATTR_ACS_BUCKET_CONFIGURATION_AUTH_ACCOUNT_KEY                 string = "bucket_configuration.authentication.account_key"
	RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_REGION                     string = "bucket_configuration.bucket.region"
	RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_NAME                       string = "bucket_configuration.bucket.name"
	RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_PREFIX                     string = "bucket_configuration.bucket.prefix"
	RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_EXPIRATION                 string = "bucket_configuration.bucket.expiration"
	RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_KEY_ID     string = "bucket_configuration.bucket_security.access_key_id"
	RES_ATTR_S3_BUCKET_CONFIGURATION_BUCKET_SECURITY_ACCESS_SECRET_KEY string = "bucket_configuration.bucket_security.secret_access_key"
	RES_ATTR_S3_BUCKET_CONFIGURATION_PRE_SIGNED_ENABLED                string = "bucket_configuration.pre_signed_url_enabled"

	RES_NAME_FMT string = "%s.test"

	RES_TYPE_BLOB_STORE_ACS   string = "sonatyperepo_blob_store_acs"
	RES_TYPE_BLOB_STORE_FILE  string = "sonatyperepo_blob_store_file"
	RES_TYPE_BLOB_STORE_GCS   string = "sonatyperepo_blob_store_gcs"
	RES_TYPE_BLOB_STORE_GROUP string = "sonatyperepo_blob_store_group"
	RES_TYPE_BLOB_STORE_S3    string = "sonatyperepo_blob_store_s3"
)

var (
	RES_NAME_BLOB_STORE_ACS   string = fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_ACS)
	RES_NAME_BLOB_STORE_FILE  string = fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_FILE)
	RES_NAME_BLOB_STORE_GCS   string = fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_GCS)
	RES_NAME_BLOB_STORE_GROUP string = fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_GROUP)
	RES_NAME_BLOB_STORE_S3    string = fmt.Sprintf(RES_NAME_FMT, RES_TYPE_BLOB_STORE_S3)
)
