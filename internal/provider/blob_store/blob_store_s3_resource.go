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
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"

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
	resp.Schema = tfschema.Schema{
		Description: "Use this data source to get a specific S3 Blob Store by it's name",
		Attributes: map[string]tfschema.Attribute{
			"name": schema.ResourceRequiredString("Name of the Blob Store"),
			"type": schema.ResourceComputedStringWithDefault(
				fmt.Sprintf("Type of this Blob Store - will always be '%s'", BLOB_STORE_TYPE_S3),
				BLOB_STORE_TYPE_S3,
			),
			"soft_quota": schema.ResourceOptionalSingleNestedAttribute(
				"Soft Quota for this Blob Store",
				map[string]tfschema.Attribute{
					"type":  schema.ResourceRequiredString("Soft Quota type"),
					"limit": schema.ResourceOptionalInt64("Quota limit"),
				},
			),
			"bucket_configuration": schema.ResourceRequiredSingleNestedAttribute(
				"Bucket Configuration for this Blob Store",
				map[string]tfschema.Attribute{
					"bucket": schema.ResourceRequiredSingleNestedAttribute(
						"Main Bucket Configuration for this Blob Store",
						map[string]tfschema.Attribute{
							"region": schema.ResourceRequiredString("The AWS region to create a new S3 bucket in or an existing S3 bucket's region"),
							"name":   schema.ResourceRequiredString("The name of the S3 bucket"),
							"prefix": schema.ResourceStringWithDefault(
								"The S3 blob store (i.e S3 object) key prefix",
								"",
							),
						},
					),
					"encryption": schema.ResourceOptionalSingleNestedAttribute(
						"Bucket Encryption Configuration for this Blob Store",
						map[string]tfschema.Attribute{
							"encryption_type": schema.ResourceStringEnum(
								"The type of S3 server side encryption to use",
								"s3ManagedEncryption",
								"kmsManagedEncryption",
							),
							"encryption_key": schema.ResourceOptionalSensitiveStringWithLengthAtLeast("The encryption key", 1),
						},
					),
					"bucket_security": schema.ResourceOptionalSingleNestedAttribute(
						"Bucket Security Configuration for this Blob Store",
						map[string]tfschema.Attribute{
							"access_key_id": schema.ResourceOptionalSensitiveStringWithLengthAtLeast("An IAM access key ID for granting access to the S3 bucket", 1),
							"secret_access_key": schema.ResourceOptionalSensitiveStringWithLengthAtLeast(
								"The secret access key associated with the specified IAM access key ID",
								1,
							),
							"role":          schema.ResourceOptionalString("An IAM role to assume in order to access the S3 bucket"),
							"session_token": schema.ResourceOptionalSensitiveStringWithLengthAtLeast("An AWS STS session token associated with temporary security credentials which grant access to the S3 bucket", 1),
						},
					),
					"advanced_bucket_connection": schema.ResourceOptionalSingleNestedAttribute(
						"Advanced Connection Configuration for this S3 Blob Store",
						map[string]tfschema.Attribute{
							"endpoint":    schema.ResourceOptionalString("A custom endpoint URL for third party object stores using the S3 API"),
							"signer_type": schema.ResourceOptionalString("An API signature version which may be required for third party object stores using the S3 API"),
							"force_path_style": schema.ResourceOptionalBool(
								"Setting this flag will result in path-style access being used for all requests",
							),
							"max_connection_pool_size": schema.ResourceOptionalInt64(
								"Setting this value will override the default connection pool size of Nexus of the s3 client for this blobstore",
							),
						},
					),
				},
			),
			"last_updated": schema.ResourceLastUpdated(),
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
	ctx = r.GetAuthContext(ctx)

	requestPayload := sonatyperepo.S3BlobStoreApiModel{
		Name: *plan.Name.ValueStringPointer(),
		BucketConfiguration: sonatyperepo.S3BlobStoreApiBucketConfiguration{
			Bucket: sonatyperepo.S3BlobStoreApiBucket{
				Region: *plan.BucketConfiguration.Bucket.Region.ValueStringPointer(),
				Name:   *plan.BucketConfiguration.Bucket.Name.ValueStringPointer(),
			},
		},
	}
	if !plan.BucketConfiguration.Bucket.Prefix.IsNull() {
		requestPayload.BucketConfiguration.Bucket.Prefix = plan.BucketConfiguration.Bucket.Prefix.ValueStringPointer()
	}
	if plan.BucketConfiguration.Encryption != nil {
		requestPayload.BucketConfiguration.Encryption = &sonatyperepo.S3BlobStoreApiEncryption{}
		if !plan.BucketConfiguration.Encryption.EncryptionType.IsNull() {
			requestPayload.BucketConfiguration.Encryption.EncryptionType = plan.BucketConfiguration.Encryption.EncryptionType.ValueStringPointer()
		}
		if !plan.BucketConfiguration.Encryption.EncryptionKey.IsNull() {
			requestPayload.BucketConfiguration.Encryption.EncryptionKey = plan.BucketConfiguration.Encryption.EncryptionKey.ValueStringPointer()
		}
	}
	if plan.BucketConfiguration.BucketSecurity != nil {
		requestPayload.BucketConfiguration.BucketSecurity = &sonatyperepo.S3BlobStoreApiBucketSecurity{}
		if !plan.BucketConfiguration.BucketSecurity.AccessKeyId.IsNull() {
			requestPayload.BucketConfiguration.BucketSecurity.AccessKeyId = plan.BucketConfiguration.BucketSecurity.AccessKeyId.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.SecretAccessKey.IsNull() {
			requestPayload.BucketConfiguration.BucketSecurity.SecretAccessKey = plan.BucketConfiguration.BucketSecurity.SecretAccessKey.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.Role.IsNull() {
			requestPayload.BucketConfiguration.BucketSecurity.Role = plan.BucketConfiguration.BucketSecurity.Role.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.SessionToken.IsNull() {
			requestPayload.BucketConfiguration.BucketSecurity.SessionToken = plan.BucketConfiguration.BucketSecurity.SessionToken.ValueStringPointer()
		}
	}
	if plan.BucketConfiguration.AdvancedBucketConnection != nil {
		requestPayload.BucketConfiguration.AdvancedBucketConnection = &sonatyperepo.S3BlobStoreApiAdvancedBucketConnection{}
		if !plan.BucketConfiguration.AdvancedBucketConnection.Endpoint.IsNull() {
			requestPayload.BucketConfiguration.AdvancedBucketConnection.Endpoint = plan.BucketConfiguration.AdvancedBucketConnection.Endpoint.ValueStringPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.SignerType.IsNull() {
			requestPayload.BucketConfiguration.AdvancedBucketConnection.SignerType = plan.BucketConfiguration.AdvancedBucketConnection.SignerType.ValueStringPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle.IsNull() {
			requestPayload.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle = plan.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle.ValueBoolPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize.IsNull() {
			max_connection_pool_size := int32(plan.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize.ValueInt64())
			requestPayload.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize = &max_connection_pool_size
		}
	}
	if plan.SoftQuota != nil {
		requestPayload.SoftQuota = &sonatyperepo.BlobStoreApiSoftQuota{
			Limit: plan.SoftQuota.Limit.ValueInt64Pointer(),
			Type:  plan.SoftQuota.Type.ValueStringPointer(),
		}
	}

	apiResponse, err := r.Client.BlobStoreAPI.CreateS3BlobStore(ctx).Body(requestPayload).Execute()

	// Handle Error
	if err != nil {
		errors.HandleAPIError(
			"Error creating S3 Blob Store",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
		return
	}

	if apiResponse.StatusCode == http.StatusCreated {
		// Set LastUpdated
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		// Inject some defaults that are not in the API response
		plan.Type = types.StringValue(BLOB_STORE_TYPE_S3)
		if plan.BucketConfiguration.Bucket.Prefix.IsNull() {
			plan.BucketConfiguration.Bucket.Prefix = types.StringValue("")
		}

		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		errors.HandleAPIError(
			"Creation of S3 Blob Store was not successful",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
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

	ctx = r.GetAuthContext(ctx)

	// Read API Call
	apiResponse, httpResponse, err := r.Client.BlobStoreAPI.GetS3BlobStore(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			errors.HandleAPIWarning(
				"S3 Blob Store to read did not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
				"Error reading S3 Blob Store",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(state.Name.ValueString())
	state.Type = types.StringValue(BLOB_STORE_TYPE_S3)
	state.BucketConfiguration.Bucket.Region = types.StringValue(apiResponse.BucketConfiguration.Bucket.Region)
	state.BucketConfiguration.Bucket.Name = types.StringValue(apiResponse.BucketConfiguration.Bucket.Name)
	if apiResponse.BucketConfiguration.Bucket.Prefix != nil {
		state.BucketConfiguration.Bucket.Prefix = types.StringValue(*apiResponse.BucketConfiguration.Bucket.Prefix)
	}
	if apiResponse.BucketConfiguration.Encryption != nil {
		if state.BucketConfiguration.Encryption == nil {
			state.BucketConfiguration.Encryption = &model.BlobStoreS3Encryption{}
		}
		if apiResponse.BucketConfiguration.Encryption.EncryptionType != nil {
			state.BucketConfiguration.Encryption.EncryptionType = types.StringValue(*apiResponse.BucketConfiguration.Encryption.EncryptionType)
		}
		if apiResponse.BucketConfiguration.Encryption.EncryptionKey != nil {
			state.BucketConfiguration.Encryption.EncryptionKey = types.StringValue(*apiResponse.BucketConfiguration.Encryption.EncryptionKey)
		}
	}
	if apiResponse.BucketConfiguration.BucketSecurity != nil {
		if state.BucketConfiguration.BucketSecurity == nil {
			state.BucketConfiguration.BucketSecurity = &model.BlobStoreS3BucketSecurityModel{}
		}
		if apiResponse.BucketConfiguration.BucketSecurity.AccessKeyId != nil {
			state.BucketConfiguration.BucketSecurity.AccessKeyId = types.StringValue(*apiResponse.BucketConfiguration.BucketSecurity.AccessKeyId)
		}
		// API does not echo back AWS Secret Access Key
		if apiResponse.BucketConfiguration.BucketSecurity.Role != nil {
			state.BucketConfiguration.BucketSecurity.Role = types.StringValue(*apiResponse.BucketConfiguration.BucketSecurity.Role)
		}
		if apiResponse.BucketConfiguration.BucketSecurity.SessionToken != nil {
			state.BucketConfiguration.BucketSecurity.SessionToken = types.StringValue(*apiResponse.BucketConfiguration.BucketSecurity.SessionToken)
		}
	}
	if apiResponse.BucketConfiguration.AdvancedBucketConnection != nil {
		if state.BucketConfiguration.AdvancedBucketConnection == nil {
			state.BucketConfiguration.AdvancedBucketConnection = &model.BlobStoreS3AdvancedBucketConnectionModel{}
		}
		if apiResponse.BucketConfiguration.AdvancedBucketConnection.Endpoint != nil {
			state.BucketConfiguration.AdvancedBucketConnection.Endpoint = types.StringValue(*apiResponse.BucketConfiguration.AdvancedBucketConnection.Endpoint)
		}
		if apiResponse.BucketConfiguration.AdvancedBucketConnection.SignerType != nil {
			state.BucketConfiguration.AdvancedBucketConnection.SignerType = types.StringValue(*apiResponse.BucketConfiguration.AdvancedBucketConnection.SignerType)
		}
		if apiResponse.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle != nil {
			state.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle = types.BoolValue(*apiResponse.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle)
		}
		if apiResponse.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize != nil {
			state.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize = types.Int64Value(int64(*apiResponse.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize))
		}
	}
	if apiResponse.SoftQuota != nil {
		state.SoftQuota = &model.BlobStoreSoftQuota{
			Type:  types.StringValue(*apiResponse.SoftQuota.Type),
			Limit: types.Int64Value(*apiResponse.SoftQuota.Limit),
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

	ctx = r.GetAuthContext(ctx)

	// Update API Call
	requestPayload := sonatyperepo.S3BlobStoreApiModel{
		Name: *plan.Name.ValueStringPointer(),
		BucketConfiguration: sonatyperepo.S3BlobStoreApiBucketConfiguration{
			Bucket: sonatyperepo.S3BlobStoreApiBucket{
				Region: *plan.BucketConfiguration.Bucket.Region.ValueStringPointer(),
				Name:   *plan.BucketConfiguration.Bucket.Name.ValueStringPointer(),
			},
		},
	}
	if !plan.BucketConfiguration.Bucket.Prefix.IsNull() {
		requestPayload.BucketConfiguration.Bucket.Prefix = plan.BucketConfiguration.Bucket.Prefix.ValueStringPointer()
	}
	if plan.BucketConfiguration.Encryption != nil {
		requestPayload.BucketConfiguration.Encryption = &sonatyperepo.S3BlobStoreApiEncryption{}
		if !plan.BucketConfiguration.Encryption.EncryptionType.IsNull() {
			requestPayload.BucketConfiguration.Encryption.EncryptionType = plan.BucketConfiguration.Encryption.EncryptionType.ValueStringPointer()
		}
		if !plan.BucketConfiguration.Encryption.EncryptionKey.IsNull() {
			requestPayload.BucketConfiguration.Encryption.EncryptionKey = plan.BucketConfiguration.Encryption.EncryptionKey.ValueStringPointer()
		}
	}
	if plan.BucketConfiguration.BucketSecurity != nil {
		requestPayload.BucketConfiguration.BucketSecurity = &sonatyperepo.S3BlobStoreApiBucketSecurity{}
		if !plan.BucketConfiguration.BucketSecurity.AccessKeyId.IsNull() {
			requestPayload.BucketConfiguration.BucketSecurity.AccessKeyId = plan.BucketConfiguration.BucketSecurity.AccessKeyId.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.SecretAccessKey.IsNull() {
			requestPayload.BucketConfiguration.BucketSecurity.SecretAccessKey = plan.BucketConfiguration.BucketSecurity.SecretAccessKey.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.Role.IsNull() {
			requestPayload.BucketConfiguration.BucketSecurity.Role = plan.BucketConfiguration.BucketSecurity.Role.ValueStringPointer()
		}
		if !plan.BucketConfiguration.BucketSecurity.SessionToken.IsNull() {
			requestPayload.BucketConfiguration.BucketSecurity.SessionToken = plan.BucketConfiguration.BucketSecurity.SessionToken.ValueStringPointer()
		}
	}
	if plan.BucketConfiguration.AdvancedBucketConnection != nil {
		requestPayload.BucketConfiguration.AdvancedBucketConnection = &sonatyperepo.S3BlobStoreApiAdvancedBucketConnection{}
		if !plan.BucketConfiguration.AdvancedBucketConnection.Endpoint.IsNull() {
			requestPayload.BucketConfiguration.AdvancedBucketConnection.Endpoint = plan.BucketConfiguration.AdvancedBucketConnection.Endpoint.ValueStringPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.SignerType.IsNull() {
			requestPayload.BucketConfiguration.AdvancedBucketConnection.SignerType = plan.BucketConfiguration.AdvancedBucketConnection.SignerType.ValueStringPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle.IsNull() {
			requestPayload.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle = plan.BucketConfiguration.AdvancedBucketConnection.ForcePathStyle.ValueBoolPointer()
		}
		if !plan.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize.IsNull() {
			max_connection_pool_size := int32(plan.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize.ValueInt64())
			requestPayload.BucketConfiguration.AdvancedBucketConnection.MaxConnectionPoolSize = &max_connection_pool_size
		}
	}
	if plan.SoftQuota != nil {
		requestPayload.SoftQuota = &sonatyperepo.BlobStoreApiSoftQuota{
			Limit: plan.SoftQuota.Limit.ValueInt64Pointer(),
			Type:  plan.SoftQuota.Type.ValueStringPointer(),
		}
	}

	api_request := r.Client.BlobStoreAPI.UpdateS3BlobStore(ctx, state.Name.ValueString()).Body(requestPayload)

	// Call API
	api_response, err := api_request.Execute()

	// Handle Error(s)
	if err != nil {
		if api_response.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			errors.HandleAPIWarning(
				"S3 Blob Store to update did not exist",
				&err,
				api_response,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
				"Error updating S3 Blob Store",
				&err,
				api_response,
				&resp.Diagnostics,
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
