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

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

// blobStoreGroupResource is the resource implementation.
type blobStoreGroupResource struct {
	common.BaseResource
}

// NewBlobStoreGroupResource is a helper function to simplify the provider implementation.
func NewBlobStoreGroupResource() resource.Resource {
	return &blobStoreGroupResource{}
}

// Metadata returns the resource type name.
func (r *blobStoreGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blob_store_group"
}

// Schema defines the schema for the resource.
func (r *blobStoreGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Description: `Manage Blob Store Groups.
		
This resource does not support promoting a Blob Store to becoming a Group - see examples.`,
		Attributes: map[string]tfschema.Attribute{
			"name": schema.ResourceRequiredString("Name of the Blob Store"),
			"soft_quota": schema.ResourceOptionalSingleNestedAttribute(
				"Soft Quota for this Blob Store",
				map[string]tfschema.Attribute{
					"type":  schema.ResourceRequiredString("Soft Quota type"),
					"limit": schema.ResourceOptionalInt64("Quota limit"),
				},
			),
			"members": schema.ResourceRequiredStringList("List of the names of blob stores that are members of this group"),
			"fill_policy": schema.ResourceRequiredStringEnum(
				"Fill Policy for this Blob Store - see [official documentation](https://help.sonatype.com/en/blob-stores.html#what-is-a-fill-policy-).",
				BLOB_STORE_FILL_POLICY_ROUND_ROBIN,
				BLOB_STORE_FILL_POLICY_WRITE_FIRST,
			),
			"last_updated": schema.ResourceLastUpdated(),
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *blobStoreGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from state
	var plan model.BlobStoreGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Call API to Create
	ctx = r.AuthContext(ctx)
	apiBody := sonatyperepo.NewGroupBlobStoreApiCreateRequestWithDefaults()
	plan.MapToApiCreate(apiBody)
	httpResponse, err := r.Client.BlobStoreAPI.CreateGroupBlobStore(ctx).Body(*apiBody).Execute()

	if err != nil {
		errors.HandleAPIError(
			"Error creating Blobstore Group",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		errors.HandleAPIError(
			"Creation of Blobstore Group was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}

	// Update State
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *blobStoreGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.BlobStoreGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Call API to Create
	ctx = r.AuthContext(ctx)
	apiResponse, httpResponse, err := r.Client.BlobStoreAPI.GetGroupBlobStoreConfiguration(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			errors.HandleAPIWarning(
				"Blobstore Group to read did not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
				"Error reading Blobstore Group",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Update State based on Response
	state.MapFromApi(apiResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *blobStoreGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.BlobStoreGroupModel
	var state model.BlobStoreGroupModel

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
	ctx = r.AuthContext(ctx)
	apiBody := sonatyperepo.NewGroupBlobStoreApiUpdateRequestWithDefaults()
	plan.MapToApiUpdate(apiBody)
	httpResponse, err := r.Client.BlobStoreAPI.UpdateGroupBlobStore(ctx, state.Name.ValueString()).Body(*apiBody).Execute()

	if err != nil {
		errors.HandleAPIError(
			"Error updating Blobstore Group",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		errors.HandleAPIError(
			"Update of Blobstore Group was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *blobStoreGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.BlobStoreGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Add Auth Context
	ctx = r.AuthContext(ctx)

	// Delete API Call
	DeleteBlobStore(r.Client, &ctx, state.Name.ValueString(), resp)
}

// This allows users to import existing Tasks into Terraform state.
func (r *blobStoreGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the Blob Store Name as the import identifier
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
