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
			"type": schema.ResourceOptionalStringWithDefault(
				fmt.Sprintf("Type of this Blob Store - will always be '%s'", common.BLOB_STORE_TYPE_S3),
				common.BLOB_STORE_TYPE_S3,
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
					"pre_signed_url_enabled": schema.ResourceOptionalBoolWithDefault(
						"Whether pre-signed URL is enabled or not. **Requires Sonatype Nexus Repository Manager 3.79.0 PRO or later**",
						false,
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
	ctx = r.AuthContext(ctx)

	apiBody := sonatyperepo.NewS3BlobStoreApiModelWithDefaults()
	plan.MapToApi(apiBody)
	httpResponse, err := r.Client.BlobStoreAPI.CreateS3BlobStore(ctx).Body(*apiBody).Execute()

	// Handle Error
	if err != nil {
		errors.HandleAPIError(
			"Error creating S3 Blob Store",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	// Handle any unexpected errors
	if httpResponse.StatusCode != http.StatusCreated {
		errors.HandleAPIError(
			"Creation of S3 Blob Store was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	// Read API Call
	apiResponse, httpResponse, err := r.Client.BlobStoreAPI.GetS3BlobStore(ctx, plan.Name.ValueString()).Execute()

	if err != nil || httpResponse.StatusCode != http.StatusOK {
		if httpResponse.StatusCode == http.StatusNotFound {
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

	// Map response to state
	plan.MapFromApi(apiResponse)
	// state.Type = types.StringValue(common.BLOB_STORE_TYPE_S3)

	// Set LastUpdated
	plan.Type = types.StringValue(common.BLOB_STORE_TYPE_S3)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Update State
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

	ctx = r.AuthContext(ctx)

	// Read API Call
	apiResponse, httpResponse, err := r.Client.BlobStoreAPI.GetS3BlobStore(ctx, state.Name.ValueString()).Execute()

	if err != nil || httpResponse.StatusCode != http.StatusOK {
		if httpResponse.StatusCode == http.StatusNotFound {
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

	// Map response to state
	state.MapFromApi(apiResponse)
	state.Type = types.StringValue(common.BLOB_STORE_TYPE_S3)

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

	ctx = r.AuthContext(ctx)

	apiBody := sonatyperepo.NewS3BlobStoreApiModelWithDefaults()
	plan.MapToApi(apiBody)

	// Call API
	apiResponse, err := r.Client.BlobStoreAPI.UpdateS3BlobStore(ctx, state.Name.ValueString()).Body(*apiBody).Execute()

	// Handle Error(s)
	if err != nil || apiResponse.StatusCode != http.StatusNoContent {
		if apiResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			errors.HandleAPIWarning(
				"S3 Blob Store to update did not exist",
				&err,
				apiResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
				"Error updating S3 Blob Store",
				&err,
				apiResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Type = types.StringValue(common.BLOB_STORE_TYPE_S3)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
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
