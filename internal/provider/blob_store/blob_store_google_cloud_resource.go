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
	"regexp"
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

type blobStoreGoogleCloudResource struct {
	common.BaseResource
}

func NewBlobStoreGoogleCloudResource() resource.Resource {
	return &blobStoreGoogleCloudResource{}
}

func (r *blobStoreGoogleCloudResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blob_store_google_cloud"
}

func (r *blobStoreGoogleCloudResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this resource to create a Google Cloud Storage Blob Store",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the Blob Store",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_]*[a-zA-Z0-9]$|^[a-zA-Z0-9]$`),
						"Name must contain only letters, numbers, hyphens, and underscores. Must start and end with a letter or number.",
					),
				},
			},
			"type": schema.StringAttribute{
				Description: fmt.Sprintf("Type of this Blob Store - will always be '%s'", BLOB_STORE_TYPE_GOOGLE_CLOUD),
				Required:    false,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(BLOB_STORE_TYPE_GOOGLE_CLOUD),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"soft_quota": schema.SingleNestedBlock{
				Description: "Soft Quota for this Blob Store",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Soft Quota type",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("spaceUsedQuota", "spaceRemainingQuota"),
						},
					},
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
							"name": schema.StringAttribute{
								Description: "The name of the Google Cloud Storage bucket",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.LengthBetween(3, 63),
									stringvalidator.RegexMatches(
										regexp.MustCompile(`^[a-z0-9][a-z0-9\-]*[a-z0-9]$|^[a-z0-9]$`),
										"Bucket name must contain only lowercase letters, numbers, and hyphens. Must start and end with a letter or number.",
									),
								},
							},
							"prefix": schema.StringAttribute{
								Description: "The path within your Cloud Storage bucket where blob data should be stored",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString(""),
								Validators: []validator.String{
									stringvalidator.LengthAtMost(1024),
									stringvalidator.RegexMatches(
										regexp.MustCompile(`^[a-zA-Z0-9\-_/]*$`),
										"Prefix must contain only letters, numbers, hyphens, underscores, and forward slashes.",
									),
								},
							},
							"region": schema.StringAttribute{
								Description: "The Google Cloud region for the bucket",
								Optional:    true,
							},
						},
					},
					"authentication": schema.SingleNestedBlock{
						Description: "Authentication Configuration for Google Cloud Storage",
						Attributes: map[string]schema.Attribute{
							"authentication_method": schema.StringAttribute{
								Description: "The type of Google Cloud authentication to use",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("accountKey", "applicationDefault"),
								},
							},
							"account_key": schema.StringAttribute{
								Description: "The credentials JSON file content",
								Optional:    true,
								Sensitive:   true,
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(10),
								},
							},
						},
					},
					"encryption": schema.SingleNestedBlock{
						Description: "Encryption Configuration for Google Cloud Storage",
						Attributes: map[string]schema.Attribute{
							"encryption_type": schema.StringAttribute{
								Description: "The type of GCP server side encryption to use",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("kms", "sse"),
								},
							},
							"encryption_key": schema.StringAttribute{
								Description: "CryptoKey ID for KMS encryption",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
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

	ctx = context.WithValue(ctx, sonatyperepo.ContextBasicAuth, r.Auth)

	requestPayload := r.buildRequestPayload(ctx, &plan, "create")
	apiResponse, err := r.Client.BlobStoreAPI.CreateBlobStore1(ctx).Body(requestPayload).Execute()

	if err != nil {
		errorBody, _ := io.ReadAll(apiResponse.Body)
		resp.Diagnostics.AddError(
			"Error creating Google Cloud Storage Blob Store",
			"Could not create Google Cloud Storage Blob Store, unexpected error: "+apiResponse.Status+": "+string(errorBody),
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
		resp.Diagnostics.AddError(
			"Failed to create Google Cloud Storage Blob Store",
			fmt.Sprintf("Unable to create Google Cloud Storage Blob Store: %d: %s", apiResponse.StatusCode, apiResponse.Status),
		)
	}
}

func (r *blobStoreGoogleCloudResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.BlobStoreGoogleCloudModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = context.WithValue(ctx, sonatyperepo.ContextBasicAuth, r.Auth)

	apiResponse, httpResponse, err := r.Client.BlobStoreAPI.GetBlobStore1(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading Google Cloud Storage Blob Store",
				fmt.Sprintf("Unable to read Google Cloud Storage Blob Store: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		}
		return
	}

	// Set basic fields
	state.Type = types.StringValue(BLOB_STORE_TYPE_GOOGLE_CLOUD)
	state.BucketConfiguration.Bucket.Name = types.StringValue(apiResponse.BucketConfiguration.Bucket.Name)
	
	// Update bucket configuration
	if apiResponse.BucketConfiguration.Bucket.Prefix != nil {
		state.BucketConfiguration.Bucket.Prefix = types.StringValue(*apiResponse.BucketConfiguration.Bucket.Prefix)
	}
	if apiResponse.BucketConfiguration.Bucket.Region != nil {
		state.BucketConfiguration.Bucket.Region = types.StringValue(*apiResponse.BucketConfiguration.Bucket.Region)
	}

	// Handle encryption configuration from API response
	if apiResponse.BucketConfiguration.Encryption != nil {
		if state.BucketConfiguration.Encryption == nil {
			state.BucketConfiguration.Encryption = &model.BlobStoreGoogleCloudEncryption{}
		}
		if apiResponse.BucketConfiguration.Encryption.EncryptionType != nil {
			state.BucketConfiguration.Encryption.EncryptionType = types.StringValue(*apiResponse.BucketConfiguration.Encryption.EncryptionType)
		}
		if apiResponse.BucketConfiguration.Encryption.EncryptionKey != nil {
			state.BucketConfiguration.Encryption.EncryptionKey = types.StringValue(*apiResponse.BucketConfiguration.Encryption.EncryptionKey)
		}
	} else {
		state.BucketConfiguration.Encryption = nil
	}

	// Handle soft quota configuration
	if apiResponse.SoftQuota != nil {
		state.SoftQuota = &model.BlobStoreSoftQuota{
			Type:  types.StringValue(*apiResponse.SoftQuota.Type),
			Limit: types.Int64Value(*apiResponse.SoftQuota.Limit),
		}
	} else {
		state.SoftQuota = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *blobStoreGoogleCloudResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state model.BlobStoreGoogleCloudModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = context.WithValue(ctx, sonatyperepo.ContextBasicAuth, r.Auth)

	requestPayload := r.buildRequestPayload(ctx, &plan, "update")
	apiResponse, err := r.Client.BlobStoreAPI.UpdateBlobStore1(ctx, state.Name.ValueString()).Body(requestPayload).Execute()

	if err != nil {
		if apiResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"Google Cloud Storage Blob Store to update did not exist",
				fmt.Sprintf("Unable to update Google Cloud Storage Blob Store: %d: %s", apiResponse.StatusCode, apiResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Updating Google Cloud Storage Blob Store",
				fmt.Sprintf("Unable to update Google Cloud Storage Blob Store: %d: %s", apiResponse.StatusCode, apiResponse.Status),
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

	ctx = context.WithValue(ctx, sonatyperepo.ContextBasicAuth, r.Auth)
	DeleteBlobStore(r.Client, &ctx, state.Name.ValueString(), resp)
}

// ====================== SHARED HELPER METHODS ======================

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