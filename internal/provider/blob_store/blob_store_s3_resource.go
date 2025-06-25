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
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// blobStoreS3Resource is the resource implementation.
type blobStoreS3Resource struct {
	common.BaseResource
}

// NewBlobStoreS3Resource is a helper function to simplify the provider implementation.
func NewBlobStoreS3Resource() resource.Resource {
	return &blobStoreS3Resource{}
}

// Metadata returns the resource type name.
func (r *blobStoreS3Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blob_store_s3"
}

// Schema defines the schema for the resource.
func (r *blobStoreS3Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Default:     stringdefault.StaticString(BLOB_STORE_TYPE_S3),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"soft_quota": schema.SingleNestedAttribute{
				Description: "Soft Quota for this Blob Store",
				Required:    false,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Soft Quota type",
						Required:    true,
						Optional:    false,
					},
					"limit": schema.Int64Attribute{
						Description: "Quota limit",
						Required:    false,
						Optional:    true,
					},
				},
			},
			"bucket_configuration": schema.SingleNestedAttribute{
				Description: "Bucket Configuration for this Blob Store",
				Required:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"bucket": schema.SingleNestedAttribute{
						Description: "Main Bucket Configuration for this Blob Store",
						Required:    true,
						Optional:    false,
						Attributes: map[string]schema.Attribute{
							"region": schema.StringAttribute{
								Description: "The AWS region to create a new S3 bucket in or an existing S3 bucket's region",
								Required:    true,
								Optional:    false,
							},
							"name": schema.StringAttribute{
								Description: "The name of the S3 bucket",
								Required:    true,
								Optional:    false,
							},
							"prefix": schema.StringAttribute{
								Description: "The S3 blob store (i.e S3 object) key prefix",
								Required:    false,
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString(""),
							},
							"expiration": schema.Int64Attribute{
								Description: "How many days until deleted blobs are finally removed from the S3 bucket (-1 to disable)",
								Required:    true,
								Optional:    false,
							},
						},
					},
					"encryption": schema.SingleNestedAttribute{
						Description: "Bucket Encryption Configuration for this Blob Store",
						Required:    false,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"encryption_type": schema.StringAttribute{
								Description: "The type of S3 server side encryption to use. Either 's3ManagedEncryption' or 'kmsManagedEncryption'",
								Required:    false,
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("s3ManagedEncryption", "kmsManagedEncryption"),
								},
							},
							"encryption_key": schema.StringAttribute{
								Description: "The encryption key",
								Required:    false,
								Optional:    true,
								Sensitive:   true,
							},
						},
					},
					"bucket_security": schema.SingleNestedAttribute{
						Description: "Bucket Security Configuration for this Blob Store",
						Required:    false,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"access_key_id": schema.StringAttribute{
								Description: "An IAM access key ID for granting access to the S3 bucket",
								Required:    false,
								Optional:    true,
								Sensitive:   true,
							},
							"secret_access_key": schema.StringAttribute{
								Description: "The secret access key associated with the specified IAM access key ID",
								Required:    false,
								Optional:    true,
								Sensitive:   true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"role": schema.StringAttribute{
								Description: "An IAM role to assume in order to access the S3 bucket",
								Required:    false,
								Optional:    true,
							},
							"session_token": schema.StringAttribute{
								Description: "An AWS STS session token associated with temporary security credentials which grant access to the S3 bucket",
								Required:    false,
								Optional:    true,
								Sensitive:   true,
							},
						},
					},
					"advanced_bucket_connection": schema.SingleNestedAttribute{
						Description: "Advanced Connection Configuration for this S3 Blob Store",
						Required:    false,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"endpoint": schema.StringAttribute{
								Description: "A custom endpoint URL for third party object stores using the S3 API",
								Required:    false,
								Optional:    true,
							},
							"signer_type": schema.StringAttribute{
								Description: "An API signature version which may be required for third party object stores using the S3 API",
								Required:    false,
								Optional:    true,
							},
							"force_path_style": schema.BoolAttribute{
								Description: "Setting this flag will result in path-style access being used for all requests",
								Required:    false,
								Optional:    true,
							},
							"max_connection_pool_size": schema.Int64Attribute{
								Description: "Setting this value will override the default connection pool size of Nexus of the s3 client for this blobstore",
								Required:    false,
								Optional:    true,
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

// Create creates the resource and sets the initial Terraform state.
func (r *blobStoreS3Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.BlobStoreS3Model

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Call API to Create
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	request_payload := sonatyperepo.S3BlobStoreApiModel{
		Name: *plan.Name.ValueStringPointer(),
		BucketConfiguration: sonatyperepo.S3BlobStoreApiBucketConfiguration{
			Bucket: sonatyperepo.S3BlobStoreApiBucket{
				Region:     *plan.BucketConfiguration.Bucket.Region.ValueStringPointer(),
				Name:       *plan.BucketConfiguration.Bucket.Name.ValueStringPointer(),
				Expiration: int32(plan.BucketConfiguration.Bucket.Expiration.ValueInt64()),
			},
		},
	}
	if !plan.BucketConfiguration.Bucket.Prefix.IsNull() {
		request_payload.BucketConfiguration.Bucket.Prefix = plan.BucketConfiguration.Bucket.Prefix.ValueStringPointer()
	}
	if plan.BucketConfiguration.Encryption != nil {
		request_payload.BucketConfiguration.Encryption = &sonatyperepo.S3BlobStoreApiEncryption{}
		if !plan.BucketConfiguration.Encryption.EncryptionType.IsNull() {
			request_payload.BucketConfiguration.Encryption.EncryptionType = plan.BucketConfiguration.Encryption.EncryptionType.ValueStringPointer()
		}
		if !plan.BucketConfiguration.Encryption.EncryptionKey.IsNull() {
			request_payload.BucketConfiguration.Encryption.EncryptionKey = plan.BucketConfiguration.Encryption.EncryptionKey.ValueStringPointer()
		}
	}
	if plan.BucketConfiguration.BucketSecurity != nil {
		request_payload.BucketConfiguration.BucketSecurity = &sonatyperepo.S3BlobStoreApiBucketSecurity{}
		if !plan.BucketConfiguration.BucketSecurity.AccessKeyId.IsNull() {
			request_payload.BucketConfiguration.BucketSecurity.AccessKeyId = plan.BucketConfiguration.BucketSecurity.AccessKeyId.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.SecretAccessKey.IsNull() {
			request_payload.BucketConfiguration.BucketSecurity.SecretAccessKey = plan.BucketConfiguration.BucketSecurity.SecretAccessKey.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.Role.IsNull() {
			request_payload.BucketConfiguration.BucketSecurity.Role = plan.BucketConfiguration.BucketSecurity.Role.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.SessionToken.IsNull() {
			request_payload.BucketConfiguration.BucketSecurity.SessionToken = plan.BucketConfiguration.BucketSecurity.SessionToken.ValueStringPointer()
		}
	}
	if plan.BucketConfiguration.AdvancedBucketConnection != nil {
		request_payload.BucketConfiguration.AdvancedBucketConnection = &sonatyperepo.S3BlobStoreApiAdvancedBucketConnection{}
		if !plan.BucketConfiguration.AdvancedBucketConnection.Endpoint.IsNull() {
			request_payload.BucketConfiguration.AdvancedBucketConnection.Endpoint = plan.BucketConfiguration.AdvancedBucketConnection.Endpoint.ValueStringPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.SignerType.IsNull() {
			request_payload.BucketConfiguration.AdvancedBucketConnection.SignerType = plan.BucketConfiguration.AdvancedBucketConnection.SignerType.ValueStringPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle.IsNull() {
			request_payload.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle = plan.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle.ValueBoolPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize.IsNull() {
			max_connection_pool_size := int32(plan.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize.ValueInt64())
			request_payload.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize = &max_connection_pool_size
		}
	}
	if plan.SoftQuota != nil {
		request_payload.SoftQuota = &sonatyperepo.BlobStoreApiSoftQuota{
			Limit: plan.SoftQuota.Limit.ValueInt64Pointer(),
			Type:  plan.SoftQuota.Type.ValueStringPointer(),
		}
	}

	api_response, err := r.Client.BlobStoreAPI.CreateS3BlobStore(ctx).Body(request_payload).Execute()

	// Handle Error
	if err != nil {
		error_body, _ := io.ReadAll(api_response.Body)
		resp.Diagnostics.AddError(
			"Error creating S3 Blob Store",
			"Could not create S3 Blob Store, unexpected error: "+api_response.Status+": "+string(error_body),
		)
		return
	}

	if api_response.StatusCode == http.StatusCreated {
		// Set LastUpdated
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		// Inject some defaults that are not in the API response
		plan.Type = types.StringValue(BLOB_STORE_TYPE_S3)
		plan.BucketConfiguration.Bucket.Prefix = types.StringValue("")

		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Failed to create S3 Blob Store",
			fmt.Sprintf("Unable to create S3 Blob Store: %d: %s", api_response.StatusCode, api_response.Status),
		)
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *blobStoreS3Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.BlobStoreS3Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Read API Call
	api_response, httpResponse, err := r.Client.BlobStoreAPI.GetS3BlobStore(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading S3 Blob Store",
				fmt.Sprintf("Unable to read S3 Blob Store: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		}
		return
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(state.Name.ValueString())
	state.Type = types.StringValue(BLOB_STORE_TYPE_S3)
	state.BucketConfiguration.Bucket.Region = types.StringValue(api_response.BucketConfiguration.Bucket.Region)
	state.BucketConfiguration.Bucket.Name = types.StringValue(api_response.BucketConfiguration.Bucket.Name)
	state.BucketConfiguration.Bucket.Expiration = types.Int64Value(int64(api_response.BucketConfiguration.Bucket.Expiration))
	if api_response.BucketConfiguration.Bucket.Prefix != nil {
		state.BucketConfiguration.Bucket.Prefix = types.StringValue(*api_response.BucketConfiguration.Bucket.Prefix)
	}
	if api_response.BucketConfiguration.Encryption != nil {
		if state.BucketConfiguration.Encryption == nil {
			state.BucketConfiguration.Encryption = &model.BlobStoreS3Encryption{}
		}
		if api_response.BucketConfiguration.Encryption.EncryptionType != nil {
			state.BucketConfiguration.Encryption.EncryptionType = types.StringValue(*api_response.BucketConfiguration.Encryption.EncryptionType)
		}
		if api_response.BucketConfiguration.Encryption.EncryptionKey != nil {
			state.BucketConfiguration.Encryption.EncryptionKey = types.StringValue(*api_response.BucketConfiguration.Encryption.EncryptionKey)
		}
	}
	if api_response.BucketConfiguration.BucketSecurity != nil {
		if state.BucketConfiguration.BucketSecurity == nil {
			state.BucketConfiguration.BucketSecurity = &model.BlobStoreS3BucketSecurityModel{}
		}
		if api_response.BucketConfiguration.BucketSecurity.AccessKeyId != nil {
			state.BucketConfiguration.BucketSecurity.AccessKeyId = types.StringValue(*api_response.BucketConfiguration.BucketSecurity.AccessKeyId)
		}
		// API does not echo back AWS Secret Access Key
		if api_response.BucketConfiguration.BucketSecurity.Role != nil {
			state.BucketConfiguration.BucketSecurity.Role = types.StringValue(*api_response.BucketConfiguration.BucketSecurity.Role)
		}
		if api_response.BucketConfiguration.BucketSecurity.SessionToken != nil {
			state.BucketConfiguration.BucketSecurity.SessionToken = types.StringValue(*api_response.BucketConfiguration.BucketSecurity.SessionToken)
		}
	}
	if api_response.BucketConfiguration.AdvancedBucketConnection != nil {
		if state.BucketConfiguration.AdvancedBucketConnection == nil {
			state.BucketConfiguration.AdvancedBucketConnection = &model.BlobStoreS3AdvancedBucketConnectionModel{}
		}
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
	if api_response.SoftQuota != nil {
		state.SoftQuota = &model.BlobStoreSoftQuota{
			Type:  types.StringValue(*api_response.SoftQuota.Type),
			Limit: types.Int64Value(*api_response.SoftQuota.Limit),
		}
	} else {
		state.SoftQuota = nil
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *blobStoreS3Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.BlobStoreS3Model
	var state model.BlobStoreS3Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Update API Call
	request_payload := sonatyperepo.S3BlobStoreApiModel{
		Name: *plan.Name.ValueStringPointer(),
		BucketConfiguration: sonatyperepo.S3BlobStoreApiBucketConfiguration{
			Bucket: sonatyperepo.S3BlobStoreApiBucket{
				Region:     *plan.BucketConfiguration.Bucket.Region.ValueStringPointer(),
				Name:       *plan.BucketConfiguration.Bucket.Name.ValueStringPointer(),
				Expiration: int32(plan.BucketConfiguration.Bucket.Expiration.ValueInt64()),
			},
		},
	}
	if !plan.BucketConfiguration.Bucket.Prefix.IsNull() {
		request_payload.BucketConfiguration.Bucket.Prefix = plan.BucketConfiguration.Bucket.Prefix.ValueStringPointer()
	}
	if plan.BucketConfiguration.Encryption != nil {
		request_payload.BucketConfiguration.Encryption = &sonatyperepo.S3BlobStoreApiEncryption{}
		if !plan.BucketConfiguration.Encryption.EncryptionType.IsNull() {
			request_payload.BucketConfiguration.Encryption.EncryptionType = plan.BucketConfiguration.Encryption.EncryptionType.ValueStringPointer()
		}
		if !plan.BucketConfiguration.Encryption.EncryptionKey.IsNull() {
			request_payload.BucketConfiguration.Encryption.EncryptionKey = plan.BucketConfiguration.Encryption.EncryptionKey.ValueStringPointer()
		}
	}
	if plan.BucketConfiguration.BucketSecurity != nil {
		request_payload.BucketConfiguration.BucketSecurity = &sonatyperepo.S3BlobStoreApiBucketSecurity{}
		if !plan.BucketConfiguration.BucketSecurity.AccessKeyId.IsNull() {
			request_payload.BucketConfiguration.BucketSecurity.AccessKeyId = plan.BucketConfiguration.BucketSecurity.AccessKeyId.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.SecretAccessKey.IsNull() {
			request_payload.BucketConfiguration.BucketSecurity.SecretAccessKey = plan.BucketConfiguration.BucketSecurity.SecretAccessKey.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.Role.IsNull() {
			request_payload.BucketConfiguration.BucketSecurity.Role = plan.BucketConfiguration.BucketSecurity.Role.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.SessionToken.IsNull() {
			request_payload.BucketConfiguration.BucketSecurity.SessionToken = plan.BucketConfiguration.BucketSecurity.SessionToken.ValueStringPointer()
		}
	}
	if plan.BucketConfiguration.AdvancedBucketConnection != nil {
		request_payload.BucketConfiguration.AdvancedBucketConnection = &sonatyperepo.S3BlobStoreApiAdvancedBucketConnection{}
		if !plan.BucketConfiguration.AdvancedBucketConnection.Endpoint.IsNull() {
			request_payload.BucketConfiguration.AdvancedBucketConnection.Endpoint = plan.BucketConfiguration.AdvancedBucketConnection.Endpoint.ValueStringPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.SignerType.IsNull() {
			request_payload.BucketConfiguration.AdvancedBucketConnection.SignerType = plan.BucketConfiguration.AdvancedBucketConnection.SignerType.ValueStringPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle.IsNull() {
			request_payload.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle = plan.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle.ValueBoolPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize.IsNull() {
			max_connection_pool_size := int32(plan.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize.ValueInt64())
			request_payload.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize = &max_connection_pool_size
		}
	}
	if plan.SoftQuota != nil {
		request_payload.SoftQuota = &sonatyperepo.BlobStoreApiSoftQuota{
			Limit: plan.SoftQuota.Limit.ValueInt64Pointer(),
			Type:  plan.SoftQuota.Type.ValueStringPointer(),
		}
	}

	api_request := r.Client.BlobStoreAPI.UpdateS3BlobStore(ctx, state.Name.ValueString()).Body(request_payload)

	// Call API
	api_response, err := api_request.Execute()

	// Handle Error(s)
	if err != nil {
		if api_response.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"S3 Blob Store to update did not exist",
				fmt.Sprintf("Unable to update S3 Blob Store: %d: %s", api_response.StatusCode, api_response.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Updating S3 Blob Store",
				fmt.Sprintf("Unable to update S3 Blob Store: %d: %s", api_response.StatusCode, api_response.Status),
			)
		}
		return
	} else if api_response.StatusCode == http.StatusNoContent {
		// Map response body to schema and populate Computed attribute values
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		// Set state to fully populated data
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *blobStoreS3Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.BlobStoreS3Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Delete API Call
	DeleteBlobStore(r.Client, &ctx, state.Name.ValueString(), resp)
}
