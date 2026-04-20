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

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// blobStoreAcsResource is the resource implementation.
type blobStoreAcsResource struct {
	common.BaseResource
}

// NewBlobStoreAcsResource is a helper function to simplify the provider implementation.
func NewBlobStoreAcsResource() resource.Resource {
	return &blobStoreAcsResource{}
}

// Metadata returns the resource type name.
func (r *blobStoreAcsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blob_store_acs"
}

// Schema defines the schema for the resource.
func (r *blobStoreAcsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Description: "Use this resource to manage an Azure Cloud Storage (ACS) Blob Store",
		Attributes: map[string]tfschema.Attribute{
			"name": schema.ResourceRequiredString("Name of the Blob Store"),
			// "type": schema.ResourceOptionalStringWithDefault(
			// 	fmt.Sprintf("Type of this Blob Store - will always be '%s'", common.BLOB_STORE_TYPE_S3),
			// 	common.BLOB_STORE_TYPE_S3,
			// ),
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
					"account_name":   schema.ResourceRequiredString("Account name found under Access keys for the storage account."),
					"container_name": schema.ResourceRequiredString("The name of an existing container to be used for storage."),
					"authentication": schema.ResourceRequiredSingleNestedAttribute(
						"Authentication to Azure for this Blob Store",
						map[string]tfschema.Attribute{
							"authentication_method": schema.ResourceRequiredStringEnum(
								"The type of Azure authentication to use.",
								common.BLOB_STORE_ACS_AUTH_METHOD_ACCOUNT_KEY,
								common.BLOB_STORE_ACS_AUTH_METHOD_ENVIRONMENT_VARIABLE,
								common.BLOB_STORE_ACS_AUTH_METHOD_MANAGED_IDENTITY,
							),
							"account_key": schema.ResourceSensitiveString("The account key"),
						},
					),
				},
			),
			"last_updated": schema.ResourceLastUpdated(),
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *blobStoreAcsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.BlobStoreAcsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Call API to Create
	httpResponse, err := r.Client.BlobStoreAPI.CreateBlobStore1(r.AuthContext(ctx)).Body(*plan.MapToApi()).Execute()

	// Handle Error
	if err != nil {
		errors.HandleAPIError(
			"Error creating Azure Cloud Storage Blob Store",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusCreated {
		errors.HandleAPIError(
			"Creation of Azure Cloud Storage Blob Store was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	// Read API Call
	apiResonse := r.readAcsBlobStore(ctx, plan.Name.ValueString(), &resp.Diagnostics, &resp.State)

	if apiResonse == nil || resp.Diagnostics.HasError() {
		return
	}

	plan.MapFromApi(apiResonse)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Update State
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *blobStoreAcsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.BlobStoreAcsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Read API Call
	apiResonse := r.readAcsBlobStore(ctx, state.Name.ValueString(), &resp.Diagnostics, &resp.State)

	if apiResonse == nil || resp.Diagnostics.HasError() {
		return
	}

	state.MapFromApi(apiResonse)

	// Update State
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *blobStoreAcsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.BlobStoreAcsModel
	var state model.BlobStoreAcsModel

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

	// Call API to Update
	httpResponse, err := r.Client.BlobStoreAPI.UpdateBlobStore1(r.AuthContext(ctx), state.Name.ValueString()).Body(*plan.MapToApi()).Execute()

	// Handle Error
	if err != nil {
		errors.HandleAPIError(
			"Error updating Azure Cloud Storage Blob Store",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		errors.HandleAPIError(
			"Updating Azure Cloud Storage Blob Store was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	// Read API Call
	apiResonse := r.readAcsBlobStore(ctx, plan.Name.ValueString(), &resp.Diagnostics, &resp.State)

	if apiResonse == nil || resp.Diagnostics.HasError() {
		return
	}

	plan.MapFromApi(apiResonse)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Update State
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *blobStoreAcsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.BlobStoreAcsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	// Delete API Call
	DeleteBlobStore(r.Client, &ctx, state.Name.ValueString(), resp)
}

// This allows users to import existing S3 Blob Stores into Terraform state.
func (r *blobStoreAcsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the Blob Store Name as the import identifier
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *blobStoreAcsResource) readAcsBlobStore(ctx context.Context, blobStoreName string, respDiagnostics *diag.Diagnostics, respState *tfsdk.State) *sonatyperepo.AzureBlobStoreApiModel {
	// Call Read API
	apiResponse, httpResponse, err := r.Client.BlobStoreAPI.GetBlobStore1(r.AuthContext(ctx), blobStoreName).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			respState.RemoveResource(ctx)
			errors.HandleAPIWarning(
				"Azure Cloud Storage Blob Store to read did not exist",
				&err,
				httpResponse,
				respDiagnostics,
			)
		} else {
			errors.HandleAPIError(
				"Error reading Azure Cloud Storage Blob Store",
				&err,
				httpResponse,
				respDiagnostics,
			)
		}
		return nil
	}

	return apiResponse
}
