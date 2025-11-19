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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	tfschema "github.com/sonatype-nexus-community/terraform-provider-shared/schema"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
)

// blobStoreFileResource is the resource implementation.
type blobStoreFileResource struct {
	common.BaseResource
}

// NewBlobStoreFileResource is a helper function to simplify the provider implementation.
func NewBlobStoreFileResource() resource.Resource {
	return &blobStoreFileResource{}
}

// Metadata returns the resource type name.
func (r *blobStoreFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blob_store_file"
}

// Schema defines the schema for the resource.
func (r *blobStoreFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get a specific File Blob Store by it's name",
		Attributes: map[string]schema.Attribute{
			"name": tfschema.RequiredString("Name of the Blob Store"),
			"path": tfschema.RequiredString("The Path on disk of this File Blob Store"),
			"soft_quota": tfschema.OptionalSingleNestedAttribute("Soft Quota for this Blob Store", map[string]schema.Attribute{
				"type":  tfschema.RequiredString("Soft Quota type"),
				"limit": tfschema.OptionalInt64("Quota limit"),
			}),
			"last_updated": tfschema.Timestamp(),
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *blobStoreFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.BlobStoreFileModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Call API to Create
	ctx = r.GetAuthContext(ctx)

	request_payload := sonatyperepo.FileBlobStoreApiCreateRequest{
		Name: plan.Name.ValueStringPointer(),
		Path: plan.Path.ValueStringPointer(),
	}
	if plan.SoftQuota != nil {
		request_payload.SoftQuota = &sonatyperepo.BlobStoreApiSoftQuota{
			Limit: plan.SoftQuota.Limit.ValueInt64Pointer(),
			Type:  plan.SoftQuota.Type.ValueStringPointer(),
		}
	}

	create_request := r.Client.BlobStoreAPI.CreateFileBlobStore(ctx).Body(request_payload)
	api_response, err := create_request.Execute()

	// Handle Error
	if err != nil {
		sharederr.HandleAPIError(
			"Error creating Blob Store File",
			&err,
			api_response,
			&resp.Diagnostics,
		)
		return
	}

	if api_response.StatusCode == http.StatusNoContent {
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		sharederr.HandleAPIError(
			"Creation of Blob Store File was not successful",
			&err,
			api_response,
			&resp.Diagnostics,
		)
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *blobStoreFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.BlobStoreFileModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.GetAuthContext(ctx)

	// Read API Call
	blobStoreApiResponse, httpResponse, err := r.Client.BlobStoreAPI.GetFileBlobStoreConfiguration(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				"Blob Store File to read did not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				"Error reading Blob Store File",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	} else {
		// Overwrite items with refreshed state
		state.Path = types.StringValue(*blobStoreApiResponse.Path)

		if blobStoreApiResponse.SoftQuota != nil {
			state.SoftQuota = &model.BlobStoreSoftQuota{
				Type:  types.StringValue(*blobStoreApiResponse.SoftQuota.Type),
				Limit: types.Int64Value(*blobStoreApiResponse.SoftQuota.Limit),
			}
		} else {
			state.SoftQuota = nil
		}
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *blobStoreFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.BlobStoreFileModel
	var state model.BlobStoreFileModel

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
	request_payload := sonatyperepo.FileBlobStoreApiUpdateRequest{
		Path: plan.Path.ValueStringPointer(),
	}
	if plan.SoftQuota != nil {
		request_payload.SoftQuota = &sonatyperepo.BlobStoreApiSoftQuota{
			Limit: plan.SoftQuota.Limit.ValueInt64Pointer(),
			Type:  plan.SoftQuota.Type.ValueStringPointer(),
		}
	}
	apiUpdateRequest := r.Client.BlobStoreAPI.UpdateFileBlobStore(ctx, state.Name.ValueString()).Body(request_payload)

	// Call API
	httpResponse, err := apiUpdateRequest.Execute()

	// Handle Error(s)
	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				"Blob Store File to update did not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				"Error updating Blob Store File",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	} else if httpResponse.StatusCode == http.StatusNoContent {
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
func (r *blobStoreFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.BlobStoreFileModel

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
