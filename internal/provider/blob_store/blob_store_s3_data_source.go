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

package blob_store

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &s3BlobStoreDataSource{}
	_ datasource.DataSourceWithConfigure = &s3BlobStoreDataSource{}
)

// BlobStoreS3DataSource is a helper function to simplify the provider implementation.
func BlobStoreS3DataSource() datasource.DataSource {
	return &s3BlobStoreDataSource{}
}

// s3BlobStoreDataSource is the data source implementation.
type s3BlobStoreDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *s3BlobStoreDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blob_store_s3"
}

// Schema defines the schema for the data source.
func (d *s3BlobStoreDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get a specific S3 Blob Store by it's name",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the Blob Store",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: fmt.Sprintf("Type of this Blob Store - will always be '%s'", BLOB_STORE_TYPE_S3),
				Required:    false,
				Optional:    true,
				Computed:    true,
			},
			"soft_quota": schema.SingleNestedAttribute{
				Description: "Soft Quota for this Blob Store",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Soft Quota type",
						Required:    false,
						Optional:    true,
						Computed:    true,
					},
					"limit": schema.Int64Attribute{
						Description: "Quota limit",
						Required:    false,
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"bucket_configuration": schema.SingleNestedAttribute{
				Description: "Bucket Configuration for this Blob Store",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"bucket": schema.SingleNestedAttribute{
						Description: "Main Bucket Configuration for this Blob Store",
						Required:    false,
						Optional:    true,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"region": schema.StringAttribute{
								Description: "The AWS region to create a new S3 bucket in or an existing S3 bucket's region",
								Required:    false,
								Optional:    true,
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "The name of the S3 bucket",
								Required:    false,
								Optional:    true,
								Computed:    true,
							},
							"prefix": schema.StringAttribute{
								Description: "The S3 blob store (i.e S3 object) key prefix",
								Required:    false,
								Optional:    true,
								Computed:    true,
							},
							"expiration": schema.Int64Attribute{
								Description: "How many days until deleted blobs are finally removed from the S3 bucket (-1 to disable)",
								Required:    false,
								Optional:    true,
								Computed:    true,
							},
						},
					},
					"encryption": schema.SingleNestedAttribute{
						Description: "Bucket Encryption Configuration for this Blob Store",
						Required:    false,
						Optional:    true,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"encryption_type": schema.StringAttribute{
								Description: "The type of S3 server side encryption to use. Either 's3ManagedEncryption' or 'kmsManagedEncryption'",
								Required:    false,
								Optional:    true,
								Computed:    true,
							},
							"encryption_key": schema.StringAttribute{
								Description: "The encryption key",
								Required:    false,
								Optional:    true,
								Computed:    true,
								Sensitive:   true,
							},
						},
					},
					"bucket_security": schema.SingleNestedAttribute{
						Description: "Bucket Security Configuration for this Blob Store",
						Required:    false,
						Optional:    true,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"access_key_id": schema.StringAttribute{
								Description: "An IAM access key ID for granting access to the S3 bucket",
								Required:    false,
								Optional:    true,
								Computed:    true,
								Sensitive:   true,
							},
							"secret_access_key": schema.StringAttribute{
								Description: "The secret access key associated with the specified IAM access key ID",
								Required:    false,
								Optional:    true,
								Computed:    true,
								Sensitive:   true,
							},
							"role": schema.StringAttribute{
								Description: "An IAM role to assume in order to access the S3 bucket",
								Required:    false,
								Optional:    true,
								Computed:    true,
							},
							"session_token": schema.StringAttribute{
								Description: "An AWS STS session token associated with temporary security credentials which grant access to the S3 bucket",
								Required:    false,
								Optional:    true,
								Computed:    true,
								Sensitive:   true,
							},
						},
					},
					"advanced_bucket_connection": schema.SingleNestedAttribute{
						Description: "Advanced Connection Configuration for this S3 Blob Store",
						Required:    false,
						Optional:    true,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"endpoint": schema.StringAttribute{
								Description: "A custom endpoint URL for third party object stores using the S3 API",
								Required:    false,
								Optional:    true,
								Computed:    true,
							},
							"signer_type": schema.StringAttribute{
								Description: "An API signature version which may be required for third party object stores using the S3 API",
								Required:    false,
								Optional:    true,
								Computed:    true,
							},
							"force_path_style": schema.BoolAttribute{
								Description: "Setting this flag will result in path-style access being used for all requests",
								Required:    false,
								Optional:    true,
								Computed:    true,
							},
							"max_connection_pool_size": schema.Int64Attribute{
								Description: "Setting this value will override the default connection pool size of Nexus of the s3 client for this blobstore",
								Required:    false,
								Optional:    true,
								Computed:    true,
							},
						},
					},
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *s3BlobStoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.BlobStoreS3Model

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.Auth,
	)

	if data.Name.IsNull() {
		resp.Diagnostics.AddError("Name must not be empty.", "Name must be provided.")
		return
	}

	api_response, _, err := d.Client.BlobStoreAPI.GetS3BlobStore(ctx, data.Name.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to S3 Read Blob Store",
			err.Error(),
		)
		return
	}

	state := model.BlobStoreS3Model{
		Name: types.StringValue(data.Name.ValueString()),
		Type: types.StringValue(BLOB_STORE_TYPE_S3),
		BucketConfiguration: &model.BlobStoreS3BucketConfigurationModel{
			Bucket: model.BlobStoreS3BucketModel{
				Region:     types.StringValue(api_response.BucketConfiguration.Bucket.Region),
				Name:       types.StringValue(api_response.BucketConfiguration.Bucket.Name),
				Expiration: types.Int64Value(int64(api_response.BucketConfiguration.Bucket.Expiration)),
			},
		},
	}
	if api_response.SoftQuota != nil && api_response.SoftQuota.Type != nil {
		state.SoftQuota = &model.BlobStoreSoftQuota{
			Type:  types.StringValue(*api_response.SoftQuota.Type),
			Limit: types.Int64Value(*api_response.SoftQuota.Limit),
		}
	}
	if api_response.BucketConfiguration.Bucket.Prefix != nil {
		state.BucketConfiguration.Bucket.Prefix = types.StringValue(*api_response.BucketConfiguration.Bucket.Prefix)
	}
	if api_response.BucketConfiguration.Encryption != nil {
		state.BucketConfiguration.Encryption = &model.BlobStoreS3Encryption{}
		if api_response.BucketConfiguration.Encryption.EncryptionType != nil {
			state.BucketConfiguration.Encryption.EncryptionType = types.StringValue(*api_response.BucketConfiguration.Encryption.EncryptionType)
		}
		if api_response.BucketConfiguration.Encryption.EncryptionKey != nil {
			state.BucketConfiguration.Encryption.EncryptionKey = types.StringValue(*api_response.BucketConfiguration.Encryption.EncryptionKey)
		}
	}
	if api_response.BucketConfiguration.BucketSecurity != nil {
		state.BucketConfiguration.BucketSecurity = &model.BlobStoreS3BucketSecurityModel{}
		if api_response.BucketConfiguration.BucketSecurity.AccessKeyId != nil {
			state.BucketConfiguration.BucketSecurity.AccessKeyId = types.StringValue(*api_response.BucketConfiguration.BucketSecurity.AccessKeyId)
		}
		if api_response.BucketConfiguration.BucketSecurity.SecretAccessKey != nil {
			state.BucketConfiguration.BucketSecurity.SecretAccessKey = types.StringValue(*api_response.BucketConfiguration.BucketSecurity.SecretAccessKey)
		}
		if api_response.BucketConfiguration.BucketSecurity.Role != nil {
			state.BucketConfiguration.BucketSecurity.Role = types.StringValue(*api_response.BucketConfiguration.BucketSecurity.Role)
		}
		if api_response.BucketConfiguration.BucketSecurity.SessionToken != nil {
			state.BucketConfiguration.BucketSecurity.SessionToken = types.StringValue(*api_response.BucketConfiguration.BucketSecurity.SessionToken)
		}
	}
	if api_response.BucketConfiguration.AdvancedBucketConnection != nil {
		state.BucketConfiguration.AdvancedBucketConnection = &model.BlobStoreS3AdvancedBucketConnectionModel{}
		if api_response.BucketConfiguration.AdvancedBucketConnection.Endpoint != nil {
			state.BucketConfiguration.AdvancedBucketConnection.Endpoint = types.StringValue(*api_response.BucketConfiguration.AdvancedBucketConnection.Endpoint)
		}
		if api_response.BucketConfiguration.AdvancedBucketConnection.SignerType != nil {
			state.BucketConfiguration.AdvancedBucketConnection.SignerType = types.StringValue(*api_response.BucketConfiguration.AdvancedBucketConnection.SignerType)
		}
		if api_response.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle != nil {
			state.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle = types.BoolValue(*api_response.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle)
		}
		if api_response.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize != nil {
			state.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize = types.Int64Value(int64(*api_response.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize))
		}
	}

	// Set state
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
