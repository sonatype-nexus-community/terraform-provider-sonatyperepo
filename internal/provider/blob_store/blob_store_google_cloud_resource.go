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
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	tfschema "github.com/sonatype-nexus-community/terraform-provider-shared/schema"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

const (
	blobStoreNamePattern = `^[a-zA-Z0-9][a-zA-Z0-9._-]*$`
	gcsBucketPattern     = `^[a-z0-9][a-z0-9\-]*[a-z0-9]$|^[a-z0-9]$`
)

type blobStoreGoogleCloudResource struct {
	common.BaseResource
}

func NewBlobStoreGoogleCloudResource() resource.Resource {
	return &blobStoreGoogleCloudResource{}
}

func (r *blobStoreGoogleCloudResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blob_store_gcs"
}

func (r *blobStoreGoogleCloudResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this resource to create a Google Cloud Storage Blob Store",
		Attributes: map[string]schema.Attribute{
			"name": tfschema.RequiredStringWithRegexAndLength(
				"Name of the Blob Store",
				regexp.MustCompile(blobStoreNamePattern),
				"Name must contain only letters, digits, underscores(_), hyphens(-), and dots(.). May not start with underscore or dot.",
				1,
				200,
			),
			"type": tfschema.ResourceComputedStringWithDefault(
				fmt.Sprintf("Type of this Blob Store - will always be '%s'", BLOB_STORE_TYPE_GOOGLE_CLOUD),
				BLOB_STORE_TYPE_GOOGLE_CLOUD,
			),
			"last_updated": tfschema.ResourceComputedString("The timestamp of when the resource was last updated"),
		},
		Blocks: map[string]schema.Block{
			"soft_quota": schema.SingleNestedBlock{
				Description: "Soft Quota for this Blob Store",
				Attributes: map[string]schema.Attribute{
					"type": tfschema.ResourceStringEnum(
						"Soft Quota type",
						"spaceUsedQuota",
						"spaceRemainingQuota",
					),
					"limit": schema.Int64Attribute{
						Description: "Quota limit in bytes",
						Optional:    true,
					},
				},
			},
			"bucket_configuration": schema.SingleNestedBlock{
				Description: "Bucket Configuration for this Google Cloud Storage Blob Store",
				Blocks: map[string]schema.Block{
					"bucket": schema.SingleNestedBlock{
						Description: "Main Bucket Configuration for this Blob Store",
						Attributes: map[string]schema.Attribute{
							"name": tfschema.RequiredStringWithRegexAndLength(
								"The name of the Google Cloud Storage bucket",
								regexp.MustCompile(gcsBucketPattern),
								"Bucket name must contain only lowercase letters, numbers, and hyphens. Must start and end with a letter or number.",
								3,
								63,
							),
							"prefix": tfschema.OptionalStringWithLengthAtMost(
								"The path within your Cloud Storage bucket where blob data should be stored",
								1024,
							),
							"region":     tfschema.ResourceOptionalString("The Google Cloud region for the bucket"),
							"project_id": tfschema.ResourceOptionalString("The Google Cloud project id for the bucket"),
						},
					},
					"authentication": schema.SingleNestedBlock{
						Description: "Authentication Configuration for Google Cloud Storage",
						Attributes: map[string]schema.Attribute{
							"authentication_method": tfschema.ResourceStringEnum(
								"The type of Google Cloud authentication to use",
								"accountKey",
								"applicationDefault",
							),
							"account_key": tfschema.OptionalSensitiveStringWithLengthAtLeast(
								"The credentials JSON file content",
								10,
							),
						},
					},
					"encryption": schema.SingleNestedBlock{
						Description: "Encryption Configuration for Google Cloud Storage",
						Attributes: map[string]schema.Attribute{
							"encryption_type": tfschema.ResourceStringEnum(
								"The type of GCP server side encryption to use",
								"kmsManagedEncryption",
								"default",
							),
							"encryption_key": tfschema.OptionalStringWithLengthAtLeast(
								"CryptoKey ID for KMS encryption",
								1,
							),
						},
					},
				},
			},
		},
	}
}

func (r *blobStoreGoogleCloudResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.BlobStoreGoogleCloudModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = r.GetAuthContext(ctx)

	requestPayload := r.buildRequestPayload(ctx, &plan, "create")
	apiResponse, err := r.Client.BlobStoreAPI.CreateBlobStore2(ctx).Body(requestPayload).Execute()

	if err != nil {
		sharederr.HandleAPIError(
			"Error creating Google Cloud Storage Blob Store",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
		return
	}

	if apiResponse.StatusCode == http.StatusCreated {
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		plan.Type = types.StringValue(BLOB_STORE_TYPE_GOOGLE_CLOUD)

		if plan.BucketConfiguration.Bucket.Prefix.IsNull() {
			plan.BucketConfiguration.Bucket.Prefix = types.StringValue("")
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	} else {
		sharederr.HandleAPIError(
			"Creation of Google Cloud Storage Blob Store was not successful",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
	}
}

func (r *blobStoreGoogleCloudResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.BlobStoreGoogleCloudModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = r.GetAuthContext(ctx)

	apiResponse, httpResponse, err := r.Client.BlobStoreAPI.GetBlobStore2(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				"Google Cloud Storage Blob Store to read did not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				"Error reading Google Cloud Storage Blob Store",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Set basic fields
	state.Type = types.StringValue(BLOB_STORE_TYPE_GOOGLE_CLOUD)

	// Populate bucket configuration from API response
	r.setBucketConfigurationFromResponse(&state, apiResponse)

	// Handle authentication, encryption and soft quota configuration
	r.setAuthenticationFromResponse(&state, apiResponse)
	r.setEncryptionFromResponse(&state, apiResponse)
	r.setSoftQuotaFromResponse(&state, apiResponse)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *blobStoreGoogleCloudResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state model.BlobStoreGoogleCloudModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = r.GetAuthContext(ctx)

	requestPayload := r.buildRequestPayload(ctx, &plan, "update")
	apiResponse, err := r.Client.BlobStoreAPI.UpdateBlobStore2(ctx, state.Name.ValueString()).Body(requestPayload).Execute()

	if err != nil {
		if apiResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				"Google Cloud Storage Blob Store to update did not exist",
				&err,
				apiResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				"Error updating Google Cloud Storage Blob Store",
				&err,
				apiResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	if apiResponse.StatusCode == http.StatusNoContent {
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (r *blobStoreGoogleCloudResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.BlobStoreGoogleCloudModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = r.GetAuthContext(ctx)
	DeleteBlobStore(r.Client, &ctx, state.Name.ValueString(), resp)
}

// ImportState implements the import functionality for the Google Cloud Storage Blob Store resource.
// This allows users to import existing blob stores into Terraform state using the blob store name as the identifier.
func (r *blobStoreGoogleCloudResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the blob store name as the import identifier
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// ====================== HELPER METHODS ======================

// buildRequestPayload constructs the API request payload for both create and update operations
func (r *blobStoreGoogleCloudResource) buildRequestPayload(ctx context.Context, plan *model.BlobStoreGoogleCloudModel, operation string) sonatyperepo.GoogleCloudBlobstoreApiModel {
	bucketConfig := r.buildBucketConfiguration(plan)

	requestPayload := sonatyperepo.GoogleCloudBlobstoreApiModel{
		Name:                *plan.Name.ValueStringPointer(),
		Type:                sonatyperepo.PtrString(BLOB_STORE_TYPE_GOOGLE_CLOUD),
		BucketConfiguration: bucketConfig,
	}

	r.configureBucketSecurity(ctx, plan, &requestPayload, operation)
	r.configureBucketEncryption(ctx, plan, &requestPayload, operation)
	r.configureSoftQuota(plan, &requestPayload)

	return requestPayload
}

// buildBucketConfiguration creates the basic bucket configuration (shared between Create/Update)
func (r *blobStoreGoogleCloudResource) buildBucketConfiguration(plan *model.BlobStoreGoogleCloudModel) sonatyperepo.GoogleCloudBlobStoreApiBucketConfiguration {
	bucketConfig := sonatyperepo.GoogleCloudBlobStoreApiBucketConfiguration{
		Bucket: sonatyperepo.GoogleCloudBlobStoreApiBucket{
			Name: *plan.BucketConfiguration.Bucket.Name.ValueStringPointer(),
		},
	}

	if !plan.BucketConfiguration.Bucket.Prefix.IsNull() {
		bucketConfig.Bucket.Prefix = plan.BucketConfiguration.Bucket.Prefix.ValueStringPointer()
	}
	if !plan.BucketConfiguration.Bucket.Region.IsNull() {
		bucketConfig.Bucket.Region = plan.BucketConfiguration.Bucket.Region.ValueStringPointer()
	}
	if !plan.BucketConfiguration.Bucket.ProjectId.IsNull() {
		bucketConfig.Bucket.ProjectId = plan.BucketConfiguration.Bucket.ProjectId.ValueStringPointer()
	}

	return bucketConfig
}

// configureBucketSecurity sets up authentication configuration (shared between Create/Update)
func (r *blobStoreGoogleCloudResource) configureBucketSecurity(ctx context.Context, plan *model.BlobStoreGoogleCloudModel, requestPayload *sonatyperepo.GoogleCloudBlobstoreApiModel, operation string) {
	if plan.BucketConfiguration.Authentication == nil || plan.BucketConfiguration.Authentication.AuthenticationMethod.IsNull() {
		return
	}

	auth := &sonatyperepo.GoogleCloudBlobStoreApiBucketAuthentication{
		AuthenticationMethod: plan.BucketConfiguration.Authentication.AuthenticationMethod.ValueString(),
	}

	if !plan.BucketConfiguration.Authentication.AccountKey.IsNull() {
		auth.AccountKey = plan.BucketConfiguration.Authentication.AccountKey.ValueStringPointer()
	}

	requestPayload.BucketConfiguration.BucketSecurity = auth

	logMsg := fmt.Sprintf("Authentication configured: %s", auth.AuthenticationMethod)
	if operation == "update" {
		logMsg = fmt.Sprintf("Authentication configured for update: %s", auth.AuthenticationMethod)
	}
	tflog.Info(ctx, logMsg)
}

// configureBucketEncryption sets up encryption configuration (shared between Create/Update)
func (r *blobStoreGoogleCloudResource) configureBucketEncryption(ctx context.Context, plan *model.BlobStoreGoogleCloudModel, requestPayload *sonatyperepo.GoogleCloudBlobstoreApiModel, operation string) {
	if plan.BucketConfiguration.Encryption == nil {
		return
	}

	if plan.BucketConfiguration.Encryption.EncryptionType.IsNull() && plan.BucketConfiguration.Encryption.EncryptionKey.IsNull() {
		return
	}

	encryption := &sonatyperepo.GoogleCloudBlobStoreApiEncryption{}

	if !plan.BucketConfiguration.Encryption.EncryptionType.IsNull() {
		encryption.EncryptionType = plan.BucketConfiguration.Encryption.EncryptionType.ValueStringPointer()
	}
	if !plan.BucketConfiguration.Encryption.EncryptionKey.IsNull() {
		encryption.EncryptionKey = plan.BucketConfiguration.Encryption.EncryptionKey.ValueStringPointer()
	}

	requestPayload.BucketConfiguration.Encryption = encryption

	encType := ""
	if encryption.EncryptionType != nil {
		encType = *encryption.EncryptionType
	}

	logMsg := fmt.Sprintf("Encryption configured: type=%s", encType)
	if operation == "update" {
		logMsg = fmt.Sprintf("Encryption configured for update: type=%s", encType)
	}
	tflog.Info(ctx, logMsg)
}

// configureSoftQuota sets up soft quota configuration (shared between Create/Update)
func (r *blobStoreGoogleCloudResource) configureSoftQuota(plan *model.BlobStoreGoogleCloudModel, requestPayload *sonatyperepo.GoogleCloudBlobstoreApiModel) {
	if plan.SoftQuota == nil {
		return
	}

	requestPayload.SoftQuota = &sonatyperepo.BlobStoreApiSoftQuota{
		Limit: plan.SoftQuota.Limit.ValueInt64Pointer(),
		Type:  plan.SoftQuota.Type.ValueStringPointer(),
	}
}

// setBucketConfigurationFromResponse handles bucket configuration from API response
func (r *blobStoreGoogleCloudResource) setBucketConfigurationFromResponse(state *model.BlobStoreGoogleCloudModel, apiResponse *sonatyperepo.GoogleCloudBlobstoreApiModel) {
	// Initialize BucketConfiguration if it's nil
	if state.BucketConfiguration == nil {
		state.BucketConfiguration = &model.BlobStoreGoogleCloudBucketConfiguration{}
	}

	// Set bucket name (required field)
	state.BucketConfiguration.Bucket.Name = types.StringValue(apiResponse.BucketConfiguration.Bucket.Name)

	// Set bucket prefix - use empty string as default if not provided
	if apiResponse.BucketConfiguration.Bucket.Prefix != nil {
		state.BucketConfiguration.Bucket.Prefix = types.StringValue(*apiResponse.BucketConfiguration.Bucket.Prefix)
	} else {
		state.BucketConfiguration.Bucket.Prefix = types.StringValue("")
	}

	// Set bucket region if provided
	if apiResponse.BucketConfiguration.Bucket.Region != nil {
		state.BucketConfiguration.Bucket.Region = types.StringValue(*apiResponse.BucketConfiguration.Bucket.Region)
	} else {
		state.BucketConfiguration.Bucket.Region = types.StringNull()
	}

	// Set bucket project ID if provided
	if apiResponse.BucketConfiguration.Bucket.ProjectId != nil {
		state.BucketConfiguration.Bucket.ProjectId = types.StringValue(*apiResponse.BucketConfiguration.Bucket.ProjectId)
	} else {
		state.BucketConfiguration.Bucket.ProjectId = types.StringNull()
	}
}

// setAuthenticationFromResponse handles authentication configuration from API response
func (r *blobStoreGoogleCloudResource) setAuthenticationFromResponse(state *model.BlobStoreGoogleCloudModel, apiResponse *sonatyperepo.GoogleCloudBlobstoreApiModel) {
	// Initialize BucketConfiguration if it's nil
	if state.BucketConfiguration == nil {
		state.BucketConfiguration = &model.BlobStoreGoogleCloudBucketConfiguration{}
	}

	if apiResponse.BucketConfiguration.BucketSecurity == nil {
		state.BucketConfiguration.Authentication = nil
		return
	}

	// Initialize Authentication if it's nil
	if state.BucketConfiguration.Authentication == nil {
		state.BucketConfiguration.Authentication = &model.BlobStoreGoogleCloudAuthentication{}
	}

	state.BucketConfiguration.Authentication.AuthenticationMethod = types.StringValue(apiResponse.BucketConfiguration.BucketSecurity.AuthenticationMethod)

	// Note: We don't read back the account key for security reasons - it's write-only
	// The account key will remain in state from the configuration but won't be updated from the API response
}

// setEncryptionFromResponse handles encryption configuration from API response
func (r *blobStoreGoogleCloudResource) setEncryptionFromResponse(state *model.BlobStoreGoogleCloudModel, apiResponse *sonatyperepo.GoogleCloudBlobstoreApiModel) {
	// Initialize BucketConfiguration if it's nil
	if state.BucketConfiguration == nil {
		state.BucketConfiguration = &model.BlobStoreGoogleCloudBucketConfiguration{}
	}

	if apiResponse.BucketConfiguration.Encryption == nil {
		state.BucketConfiguration.Encryption = nil
		return
	}

	// Initialize Encryption if it's nil
	if state.BucketConfiguration.Encryption == nil {
		state.BucketConfiguration.Encryption = &model.BlobStoreGoogleCloudEncryption{}
	}

	if apiResponse.BucketConfiguration.Encryption.EncryptionType != nil {
		state.BucketConfiguration.Encryption.EncryptionType = types.StringValue(*apiResponse.BucketConfiguration.Encryption.EncryptionType)
	}
	if apiResponse.BucketConfiguration.Encryption.EncryptionKey != nil {
		state.BucketConfiguration.Encryption.EncryptionKey = types.StringValue(*apiResponse.BucketConfiguration.Encryption.EncryptionKey)
	}
}

// setSoftQuotaFromResponse handles soft quota configuration from API response
func (r *blobStoreGoogleCloudResource) setSoftQuotaFromResponse(state *model.BlobStoreGoogleCloudModel, apiResponse *sonatyperepo.GoogleCloudBlobstoreApiModel) {
	if apiResponse.SoftQuota == nil {
		state.SoftQuota = nil
		return
	}

	state.SoftQuota = &model.BlobStoreSoftQuota{
		Type:  types.StringValue(*apiResponse.SoftQuota.Type),
		Limit: types.Int64Value(*apiResponse.SoftQuota.Limit),
	}
}
