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
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/types"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	"github.com/sonatype-nexus-community/terraform-provider-shared/util"
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
	if api == nil {
		api = v3.NewBlobStoreApiSoftQuotaWithDefaults()
	}
	api.Type = m.Type.ValueStringPointer()
	api.Limit = m.Limit.ValueInt64Pointer()
}

// BlobStoreFileModel
// ------------------------------------
type BlobStoreFileModelDS struct {
	Name      types.String        `tfsdk:"name"`
	Path      types.String        `tfsdk:"path"`
	SoftQuota *BlobStoreSoftQuota `tfsdk:"soft_quota"`
}

type BlobStoreFileModel struct {
	Name        types.String        `tfsdk:"name"`
	Path        types.String        `tfsdk:"path"`
	SoftQuota   *BlobStoreSoftQuota `tfsdk:"soft_quota"`
	LastUpdated types.String        `tfsdk:"last_updated"`
}

// BlobStoreGroupModel
// ------------------------------------
type BlobStoreGroupModelDS struct {
	Name       types.String        `tfsdk:"name"`
	SoftQuota  *BlobStoreSoftQuota `tfsdk:"soft_quota"`
	Members    []types.String      `tfsdk:"members"`
	FillPolicy types.String        `tfsdk:"fill_policy"`
}

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
type BlobStoreS3ModelDS struct {
	Name                types.String                         `tfsdk:"name"`
	Type                types.String                         `tfsdk:"type"`
	SoftQuota           *BlobStoreSoftQuota                  `tfsdk:"soft_quota"`
	BucketConfiguration *BlobStoreS3BucketConfigurationModel `tfsdk:"bucket_configuration"`
}

type BlobStoreS3ModelV0 struct {
	Name                types.String                           `tfsdk:"name"`
	Type                types.String                           `tfsdk:"type"`
	SoftQuota           *BlobStoreSoftQuota                    `tfsdk:"soft_quota"`
	BucketConfiguration *BlobStoreS3BucketConfigurationModelV0 `tfsdk:"bucket_configuration"`
	LastUpdated         types.String                           `tfsdk:"last_updated"`
}
type BlobStoreS3ModelV1 struct {
	Name                types.String                         `tfsdk:"name"`
	Type                types.String                         `tfsdk:"type"`
	SoftQuota           *BlobStoreSoftQuota                  `tfsdk:"soft_quota"`
	BucketConfiguration *BlobStoreS3BucketConfigurationModel `tfsdk:"bucket_configuration"`
	LastUpdated         types.String                         `tfsdk:"last_updated"`
}
type BlobStoreS3Model = BlobStoreS3ModelV1

func (m *BlobStoreS3Model) MapFromApi(api *v3.S3BlobStoreApiModel) {
	m.Name = types.StringValue(api.Name)
	m.Type = types.StringValue(common.BLOB_STORE_TYPE_S3)
	if api.SoftQuota != nil {
		m.SoftQuota.MapFromApi(api.SoftQuota)
	}
	m.BucketConfiguration.MapFromApi(&api.BucketConfiguration)
}

func (m *BlobStoreS3Model) MapToApi(api *v3.S3BlobStoreApiModel) {
	api.Name = m.Name.ValueString()
	// api.Type = m.Type.ValueStringPointer()
	if m.SoftQuota != nil {
		m.SoftQuota.MapToApi(api.SoftQuota)
	}
	m.BucketConfiguration.MapToApi(&api.BucketConfiguration)
}

// BlobStoreS3BucketConfigurationModel
// ------------------------------------
type BlobStoreS3BucketConfigurationModelV0 struct {
	Bucket                   BlobStoreS3BucketModel                    `tfsdk:"bucket"`
	Encryption               *BlobStoreS3Encryption                    `tfsdk:"encryption"`
	BucketSecurity           *BlobStoreS3BucketSecurityModel           `tfsdk:"bucket_security"`
	AdvancedBucketConnection *BlobStoreS3AdvancedBucketConnectionModel `tfsdk:"advanced_bucket_connection"`
}
type BlobStoreS3BucketConfigurationModelV1 struct {
	Bucket                   BlobStoreS3BucketModel                    `tfsdk:"bucket"`
	Encryption               *BlobStoreS3Encryption                    `tfsdk:"encryption"`
	BucketSecurity           *BlobStoreS3BucketSecurityModel           `tfsdk:"bucket_security"`
	AdvancedBucketConnection *BlobStoreS3AdvancedBucketConnectionModel `tfsdk:"advanced_bucket_connection"`
	PreSignedUrlEnabled      types.Bool                                `tfsdk:"pre_signed_url_enabled"`
}
type BlobStoreS3BucketConfigurationModel = BlobStoreS3BucketConfigurationModelV1

func (m *BlobStoreS3BucketConfigurationModel) MapFromApi(api *v3.S3BlobStoreApiBucketConfiguration) {
	m.Bucket.MapFromApi(&api.Bucket)
	if api.Encryption != nil {
		m.Encryption.MapFromApi(api.Encryption)
	}
	if api.BucketSecurity != nil {
		m.BucketSecurity.MapFromApi(api.BucketSecurity)
	}
	if api.AdvancedBucketConnection != nil {
		m.AdvancedBucketConnection.MapFromApi(api.AdvancedBucketConnection)
	}
	if api.PreSignedUrlEnabled == nil {
		m.PreSignedUrlEnabled = types.BoolValue(false)
	} else {
		m.PreSignedUrlEnabled = types.BoolPointerValue(api.PreSignedUrlEnabled)
	}
}

func (m *BlobStoreS3BucketConfigurationModel) MapToApi(api *v3.S3BlobStoreApiBucketConfiguration) {
	m.Bucket.MapToApi(&api.Bucket)
	if m.Encryption != nil {
		m.Encryption.MapToApi(api.Encryption)
	}
	if m.BucketSecurity != nil {
		if api.BucketSecurity == nil {
			api.BucketSecurity = v3.NewS3BlobStoreApiBucketSecurityWithDefaults()
		}
		m.BucketSecurity.MapToApi(api.BucketSecurity)
	}
	if m.AdvancedBucketConnection != nil {
		if api.AdvancedBucketConnection == nil {
			api.AdvancedBucketConnection = v3.NewS3BlobStoreApiAdvancedBucketConnectionWithDefaults()
		}
		m.AdvancedBucketConnection.MapToApi(api.AdvancedBucketConnection)
	}
	api.PreSignedUrlEnabled = util.BoolToPtr(m.PreSignedUrlEnabled.ValueBool())
}

// BlobStoreS3BucketModel
// ------------------------------------
type BlobStoreS3BucketModel struct {
	Region types.String `tfsdk:"region"`
	Name   types.String `tfsdk:"name"`
	Prefix types.String `tfsdk:"prefix"`
}

func (m *BlobStoreS3BucketModel) MapFromApi(api *v3.S3BlobStoreApiBucket) {
	m.Region = types.StringValue(api.Region)
	m.Name = types.StringValue(api.Name)
	m.Prefix = types.StringPointerValue(api.Prefix)
}

func (m *BlobStoreS3BucketModel) MapToApi(api *v3.S3BlobStoreApiBucket) {
	api.Region = m.Region.ValueString()
	api.Name = m.Name.ValueString()
	api.Prefix = m.Prefix.ValueStringPointer()
}

// BlobStoreS3Encryption
// ------------------------------------
type BlobStoreS3Encryption struct {
	EncryptionType types.String `tfsdk:"encryption_type"`
	EncryptionKey  types.String `tfsdk:"encryption_key"`
}

func (m *BlobStoreS3Encryption) MapFromApi(api *v3.S3BlobStoreApiEncryption) {
	if api == nil {
		api = v3.NewS3BlobStoreApiEncryptionWithDefaults()
	}
	m.EncryptionType = types.StringPointerValue(api.EncryptionType)
	m.EncryptionKey = types.StringPointerValue(api.EncryptionKey)
}

func (m *BlobStoreS3Encryption) MapToApi(api *v3.S3BlobStoreApiEncryption) {
	api.EncryptionType = m.EncryptionType.ValueStringPointer()
	api.EncryptionKey = m.EncryptionKey.ValueStringPointer()
}

// BlobStoreS3BucketSecurityModel
// ------------------------------------
type BlobStoreS3BucketSecurityModel struct {
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	Role            types.String `tfsdk:"role"`
	SessionToken    types.String `tfsdk:"session_token"`
}

func (m *BlobStoreS3BucketSecurityModel) MapFromApi(api *v3.S3BlobStoreApiBucketSecurity) {
	m.AccessKeyId = types.StringPointerValue(api.AccessKeyId)
	// m.SecretAccessKey = types.StringPointerValue(api.SecretAccessKey)
	m.Role = types.StringPointerValue(api.Role)
	m.SessionToken = types.StringPointerValue(api.SessionToken)
}

func (m *BlobStoreS3BucketSecurityModel) MapToApi(api *v3.S3BlobStoreApiBucketSecurity) {
	api.AccessKeyId = m.AccessKeyId.ValueStringPointer()
	api.SecretAccessKey = m.SecretAccessKey.ValueStringPointer()
	api.Role = m.Role.ValueStringPointer()
	api.SessionToken = m.SessionToken.ValueStringPointer()
}

// BlobStoreS3AdvancedBucketConnectionModel
// ------------------------------------
type BlobStoreS3AdvancedBucketConnectionModel struct {
	Endpoint              types.String `tfsdk:"endpoint"`
	SignerType            types.String `tfsdk:"signer_type"`
	ForcePathStyle        types.Bool   `tfsdk:"force_path_style"`
	MaxConnectionPoolSize types.Int64  `tfsdk:"max_connection_pool_size"`
}

func (m *BlobStoreS3AdvancedBucketConnectionModel) MapFromApi(api *v3.S3BlobStoreApiAdvancedBucketConnection) {
	m.Endpoint = types.StringPointerValue(api.Endpoint)
	m.SignerType = types.StringPointerValue(api.SignerType)
	m.ForcePathStyle = types.BoolPointerValue(api.ForcePathStyle)
	m.MaxConnectionPoolSize = util.Int32PtrToValue(api.MaxConnectionPoolSize)
}

func (m *BlobStoreS3AdvancedBucketConnectionModel) MapToApi(api *v3.S3BlobStoreApiAdvancedBucketConnection) {
	api.Endpoint = m.Endpoint.ValueStringPointer()
	api.SignerType = m.SignerType.ValueStringPointer()
	api.ForcePathStyle = m.ForcePathStyle.ValueBoolPointer()
	api.MaxConnectionPoolSize = util.Int32ToPtr(int32(m.MaxConnectionPoolSize.ValueInt64()))
}

// BlobStoreGoogleCloudModel
// ------------------------------------
type BlobStoreGoogleCloudModel struct {
	Name                types.String                             `tfsdk:"name"`
	Type                types.String                             `tfsdk:"type"`
	BucketConfiguration *BlobStoreGoogleCloudBucketConfiguration `tfsdk:"bucket_configuration"`
	SoftQuota           *BlobStoreSoftQuota                      `tfsdk:"soft_quota"`
	LastUpdated         types.String                             `tfsdk:"last_updated"`
}

// BlobStoreGoogleCloudBucketConfiguration
// ------------------------------------
type BlobStoreGoogleCloudBucketConfiguration struct {
	Bucket         BlobStoreGoogleCloudBucket          `tfsdk:"bucket"`
	Authentication *BlobStoreGoogleCloudAuthentication `tfsdk:"authentication"`
	Encryption     *BlobStoreGoogleCloudEncryption     `tfsdk:"encryption"`
}

// BlobStoreGoogleCloudBucket
// ------------------------------------
type BlobStoreGoogleCloudBucket struct {
	Name      types.String `tfsdk:"name"`
	Prefix    types.String `tfsdk:"prefix"`
	Region    types.String `tfsdk:"region"`
	ProjectId types.String `tfsdk:"project_id"`
}

// BlobStoreGoogleCloudAuthentication
// ------------------------------------
type BlobStoreGoogleCloudAuthentication struct {
	AuthenticationMethod types.String `tfsdk:"authentication_method"`
	AccountKey           types.String `tfsdk:"account_key"`
}

// BlobStoreGoogleCloudEncryption
// ------------------------------------
type BlobStoreGoogleCloudEncryption struct {
	EncryptionType types.String `tfsdk:"encryption_type"`
	EncryptionKey  types.String `tfsdk:"encryption_key"`
}
