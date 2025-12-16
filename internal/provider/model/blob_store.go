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

package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type BlobStoreModel struct {
	Name                  types.String       `tfsdk:"name"`
	Type                  types.String       `tfsdk:"type"`
	Unavailable           types.Bool         `tfsdk:"unavailable"`
	BlobCount             types.Int64        `tfsdk:"blob_count"`
	TotalSizeInBytes      types.Int64        `tfsdk:"total_size_in_bytes"`
	AvailableSpaceInBytes types.Int64        `tfsdk:"available_space_in_bytes"`
	SoftQuota             BlobStoreSoftQuota `tfsdk:"soft_quota"`
}

type BlobStoresModel struct {
	BlobStores []BlobStoreModel `tfsdk:"blob_stores"`
}

// BlobStoreSoftQuota
// ------------------------------------
type BlobStoreSoftQuota struct {
	Type  types.String `tfsdk:"type"`
	Limit types.Int64  `tfsdk:"limit"`
}

func (m *BlobStoreSoftQuota) MapFromApi(api *v3.BlobStoreApiSoftQuota) {
	m.Type = types.StringPointerValue(api.Type)
	m.Limit = types.Int64PointerValue(api.Limit)
}

func (m *BlobStoreSoftQuota) MapToApi(api *v3.BlobStoreApiSoftQuota) {
	api.Type = m.Type.ValueStringPointer()
	api.Limit = m.Limit.ValueInt64Pointer()
}

// BlobStoreFileModel
// ------------------------------------
type BlobStoreFileModel struct {
	Name        types.String        `tfsdk:"name"`
	Path        types.String        `tfsdk:"path"`
	SoftQuota   *BlobStoreSoftQuota `tfsdk:"soft_quota"`
	LastUpdated types.String        `tfsdk:"last_updated"`
}

// BlobStoreGroupModel
// ------------------------------------
type BlobStoreGroupModel struct {
	Name        types.String        `tfsdk:"name"`
	SoftQuota   *BlobStoreSoftQuota `tfsdk:"soft_quota"`
	Members     []types.String      `tfsdk:"members"`
	FillPolicy  types.String        `tfsdk:"fill_policy"`
	LastUpdated types.String        `tfsdk:"last_updated"`
}

func (m *BlobStoreGroupModel) MapFromApi(api *v3.GroupBlobStoreApiModel) {
	// Name is not in API response
	m.FillPolicy = types.StringPointerValue(api.FillPolicy)
	if api.SoftQuota != nil {
		m.SoftQuota = &BlobStoreSoftQuota{}
		m.SoftQuota.MapFromApi(api.SoftQuota)
	}
	m.Members = make([]types.String, 0)
	for _, member := range api.Members {
		m.Members = append(m.Members, types.StringValue(member))
	}
}

func (m *BlobStoreGroupModel) MapToApiCreate(api *v3.GroupBlobStoreApiCreateRequest) {
	api.Name = m.Name.ValueStringPointer()
	m.mapCommonGroupFields(api)
}

func (m *BlobStoreGroupModel) MapToApiUpdate(api *v3.GroupBlobStoreApiUpdateRequest) {
	m.mapCommonGroupFields(api)
}

func (m *BlobStoreGroupModel) mapCommonGroupFields(api interface {
	SetSoftQuota(v3.BlobStoreApiSoftQuota)
	SetFillPolicy(string)
	SetMembers([]string)
}) {
	if m.SoftQuota != nil {
		softQuota := v3.NewBlobStoreApiSoftQuotaWithDefaults()
		m.SoftQuota.MapToApi(softQuota)
		api.SetSoftQuota(*softQuota)
	}
	api.SetFillPolicy(m.FillPolicy.ValueString())

	members := make([]string, len(m.Members))
	for i, member := range m.Members {
		members[i] = member.ValueString()
	}
	api.SetMembers(members)
}

// BlobStoreS3Model
// ------------------------------------
type BlobStoreS3Model struct {
	Name                types.String                         `tfsdk:"name"`
	Type                types.String                         `tfsdk:"type"`
	SoftQuota           *BlobStoreSoftQuota                  `tfsdk:"soft_quota"`
	BucketConfiguration *BlobStoreS3BucketConfigurationModel `tfsdk:"bucket_configuration"`
	LastUpdated         types.String                         `tfsdk:"last_updated"`
}

type BlobStoreS3BucketConfigurationModel struct {
	Bucket                   BlobStoreS3BucketModel                    `tfsdk:"bucket"`
	Encryption               *BlobStoreS3Encryption                    `tfsdk:"encryption"`
	BucketSecurity           *BlobStoreS3BucketSecurityModel           `tfsdk:"bucket_security"`
	AdvancedBucketConnection *BlobStoreS3AdvancedBucketConnectionModel `tfsdk:"advanced_bucket_connection"`
}

type BlobStoreS3BucketModel struct {
	Region types.String `tfsdk:"region"`
	Name   types.String `tfsdk:"name"`
	Prefix types.String `tfsdk:"prefix"`
}

type BlobStoreS3Encryption struct {
	EncryptionType types.String `tfsdk:"encryption_type"`
	EncryptionKey  types.String `tfsdk:"encryption_key"`
}

type BlobStoreS3BucketSecurityModel struct {
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	Role            types.String `tfsdk:"role"`
	SessionToken    types.String `tfsdk:"session_token"`
}

type BlobStoreS3AdvancedBucketConnectionModel struct {
	Endpoint              types.String `tfsdk:"endpoint"`
	SignerType            types.String `tfsdk:"signer_type"`
	ForcePathStyle        types.Bool   `tfsdk:"force_path_style"`
	MaxConnectionPoolSize types.Int64  `tfsdk:"max_connection_pool_size"`
}

type BlobStoreGoogleCloudModel struct {
	Name                types.String                             `tfsdk:"name"`
	Type                types.String                             `tfsdk:"type"`
	BucketConfiguration *BlobStoreGoogleCloudBucketConfiguration `tfsdk:"bucket_configuration"`
	SoftQuota           *BlobStoreSoftQuota                      `tfsdk:"soft_quota"`
	LastUpdated         types.String                             `tfsdk:"last_updated"`
}

type BlobStoreGoogleCloudBucketConfiguration struct {
	Bucket         BlobStoreGoogleCloudBucket          `tfsdk:"bucket"`
	Authentication *BlobStoreGoogleCloudAuthentication `tfsdk:"authentication"`
	Encryption     *BlobStoreGoogleCloudEncryption     `tfsdk:"encryption"`
}

type BlobStoreGoogleCloudBucket struct {
	Name      types.String `tfsdk:"name"`
	Prefix    types.String `tfsdk:"prefix"`
	Region    types.String `tfsdk:"region"`
	ProjectId types.String `tfsdk:"project_id"`
}

type BlobStoreGoogleCloudAuthentication struct {
	AuthenticationMethod types.String `tfsdk:"authentication_method"`
	AccountKey           types.String `tfsdk:"account_key"`
}

type BlobStoreGoogleCloudEncryption struct {
	EncryptionType types.String `tfsdk:"encryption_type"`
	EncryptionKey  types.String `tfsdk:"encryption_key"`
}
