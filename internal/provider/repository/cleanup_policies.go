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

package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// cleanupPolicyResource is the resource implementation.
type cleanupPolicyResource struct {
	common.BaseResource
}

// NewCleanupPolicyResource is a helper function to simplify the provider implementation.
func NewCleanupPolicyResource() resource.Resource {
	return &cleanupPolicyResource{}
}

// Metadata returns the resource type name.
func (r *cleanupPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cleanup_policy"
}

// Schema defines the schema for the resource.
func (r *cleanupPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this resource to create and manage cleanup policies in Sonatype Nexus Repository Manager",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the cleanup policy",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"notes": schema.StringAttribute{
				Description: "Notes for the cleanup policy",
				Optional:    true,
			},
			"format": schema.StringAttribute{
				Description: "Repository format that this cleanup policy applies to",
				Required:    true,
			},
			"criteria": schema.SingleNestedAttribute{
				Description: "Cleanup criteria for this policy - at least one criterion must be specified",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"last_blob_updated": schema.Int64Attribute{
						Description: "Remove components that haven't been downloaded in this many days",
						Optional:    true,
					},
					"last_downloaded": schema.Int64Attribute{
						Description: "Remove components that were last downloaded more than this many days ago",
						Optional:    true,
					},
					"release_type": schema.StringAttribute{
						Description: "Remove components that match this release type (e.g., RELEASES, PRERELEASES)",
						Optional:    true,
					},
					"asset_regex": schema.StringAttribute{
						Description: "Remove components that have at least one asset name matching this regular expression",
						Optional:    true,
					},
				},
			},
			"retain": schema.Int64Attribute{
				Description: "Minimum number of component versions to retain",
				Optional:    true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// validateCriteria validates that at least one required criterion is provided
func validateCriteria(criteria *model.CleanupPolicyCriteriaModel) error {
	if criteria == nil {
		return fmt.Errorf("criteria block is required")
	}

	hasValidCriteria := !criteria.LastBlobUpdated.IsNull() ||
		!criteria.LastDownloaded.IsNull() ||
		!criteria.AssetRegex.IsNull()

	if !hasValidCriteria {
		return fmt.Errorf("at least one criterion (last_blob_updated, last_downloaded, or asset_regex) must be specified")
	}

	return nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *cleanupPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.CleanupPolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Validate criteria
	if err := validateCriteria(plan.Criteria); err != nil {
		resp.Diagnostics.AddError(
			"Invalid cleanup policy configuration",
			err.Error(),
		)
		return
	}

	// Call API to Create
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	requestPayload := sonatyperepo.CleanupPolicyResourceXO{
		Name:   plan.Name.ValueString(),
		Format: plan.Format.ValueString(),
	}

	if !plan.Notes.IsNull() {
		notes := plan.Notes.ValueString()
		requestPayload.Notes = &notes
	}

	// Set criteria fields
	if !plan.Criteria.LastBlobUpdated.IsNull() {
		lastBlobUpdated := plan.Criteria.LastBlobUpdated.ValueInt64()
		requestPayload.CriteriaLastBlobUpdated = &lastBlobUpdated
	}
	
	if !plan.Criteria.LastDownloaded.IsNull() {
		lastDownloaded := plan.Criteria.LastDownloaded.ValueInt64()
		requestPayload.CriteriaLastDownloaded = &lastDownloaded
	}
	
	if !plan.Criteria.ReleaseType.IsNull() {
		releaseType := plan.Criteria.ReleaseType.ValueString()
		requestPayload.CriteriaReleaseType = &releaseType
	}
	
	if !plan.Criteria.AssetRegex.IsNull() {
		assetRegex := plan.Criteria.AssetRegex.ValueString()
		requestPayload.CriteriaAssetRegex = &assetRegex
	}

	if !plan.Retain.IsNull() {
		retain := int32(plan.Retain.ValueInt64())
		requestPayload.Retain = &retain
	}

	apiResponse, err := r.Client.CleanupPoliciesAPI.Create1(ctx).Body(requestPayload).Execute()

	// Handle Error
	if err != nil {
		errorBody, _ := io.ReadAll(apiResponse.Body)
		resp.Diagnostics.AddError(
			"Error creating cleanup policy",
			"Could not create cleanup policy, unexpected error: "+apiResponse.Status+": "+string(errorBody),
		)
		return
	}

	if apiResponse.StatusCode == http.StatusCreated || apiResponse.StatusCode == http.StatusNoContent {
		// Set LastUpdated
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Failed to create cleanup policy",
			fmt.Sprintf("Unable to create cleanup policy: %d: %s", apiResponse.StatusCode, apiResponse.Status),
		)
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *cleanupPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.CleanupPolicyModel

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
	httpResponse, err := r.Client.CleanupPoliciesAPI.GetCleanupPolicyByName(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading cleanup policy",
				"Unable to read cleanup policy: "+err.Error(),
			)
		}
		return
	}

	// Parse response body to get the actual cleanup policy data
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading response body",
			"Could not read cleanup policy response: "+err.Error(),
		)
		return
	}

	// Unmarshal the response to get cleanup policy data
	var cleanupPolicy sonatyperepo.CleanupPolicyResourceXO
	err = json.Unmarshal(responseBody, &cleanupPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing response",
			"Could not parse cleanup policy response: "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(cleanupPolicy.Name)
	state.Format = types.StringValue(cleanupPolicy.Format)
	
	if cleanupPolicy.Notes != nil {
		state.Notes = types.StringValue(*cleanupPolicy.Notes)
	} else {
		state.Notes = types.StringNull()
	}

	// Handle criteria - ensure we always have a criteria object since it's required
	if state.Criteria == nil {
		state.Criteria = &model.CleanupPolicyCriteriaModel{}
	}
	
	if cleanupPolicy.CriteriaLastBlobUpdated != nil {
		state.Criteria.LastBlobUpdated = types.Int64Value(*cleanupPolicy.CriteriaLastBlobUpdated)
	} else {
		state.Criteria.LastBlobUpdated = types.Int64Null()
	}
	
	if cleanupPolicy.CriteriaLastDownloaded != nil {
		state.Criteria.LastDownloaded = types.Int64Value(*cleanupPolicy.CriteriaLastDownloaded)
	} else {
		state.Criteria.LastDownloaded = types.Int64Null()
	}
	
	if cleanupPolicy.CriteriaReleaseType != nil {
		state.Criteria.ReleaseType = types.StringValue(*cleanupPolicy.CriteriaReleaseType)
	} else {
		state.Criteria.ReleaseType = types.StringNull()
	}
	
	if cleanupPolicy.CriteriaAssetRegex != nil {
		state.Criteria.AssetRegex = types.StringValue(*cleanupPolicy.CriteriaAssetRegex)
	} else {
		state.Criteria.AssetRegex = types.StringNull()
	}

	if cleanupPolicy.Retain != nil {
		state.Retain = types.Int64Value(int64(*cleanupPolicy.Retain))
	} else {
		state.Retain = types.Int64Null()
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *cleanupPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.CleanupPolicyModel
	var state model.CleanupPolicyModel

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

	// Validate criteria
	if err := validateCriteria(plan.Criteria); err != nil {
		resp.Diagnostics.AddError(
			"Invalid cleanup policy configuration",
			err.Error(),
		)
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Update API Call
	requestPayload := sonatyperepo.CleanupPolicyResourceXO{
		Name:   plan.Name.ValueString(),
		Format: plan.Format.ValueString(),
	}

	if !plan.Notes.IsNull() {
		notes := plan.Notes.ValueString()
		requestPayload.Notes = &notes
	}

	// Set criteria fields
	if !plan.Criteria.LastBlobUpdated.IsNull() {
		lastBlobUpdated := plan.Criteria.LastBlobUpdated.ValueInt64()
		requestPayload.CriteriaLastBlobUpdated = &lastBlobUpdated
	}
	
	if !plan.Criteria.LastDownloaded.IsNull() {
		lastDownloaded := plan.Criteria.LastDownloaded.ValueInt64()
		requestPayload.CriteriaLastDownloaded = &lastDownloaded
	}
	
	if !plan.Criteria.ReleaseType.IsNull() {
		releaseType := plan.Criteria.ReleaseType.ValueString()
		requestPayload.CriteriaReleaseType = &releaseType
	}
	
	if !plan.Criteria.AssetRegex.IsNull() {
		assetRegex := plan.Criteria.AssetRegex.ValueString()
		requestPayload.CriteriaAssetRegex = &assetRegex
	}

	if !plan.Retain.IsNull() {
		retain := int32(plan.Retain.ValueInt64())
		requestPayload.Retain = &retain
	}

	api_request := r.Client.CleanupPoliciesAPI.Update1(ctx, state.Name.ValueString()).Body(requestPayload)

	// Call API
	api_response, err := api_request.Execute()

	// Handle Error(s)
	if err != nil {
		if api_response != nil && api_response.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"Cleanup policy to update did not exist",
				fmt.Sprintf("Unable to update cleanup policy: %d: %s", api_response.StatusCode, api_response.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Updating cleanup policy",
				fmt.Sprintf("Unable to update cleanup policy: %s", err.Error()),
			)
		}
		return
	} else if api_response.StatusCode == http.StatusNoContent || api_response.StatusCode == http.StatusOK {
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
func (r *cleanupPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.CleanupPolicyModel

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
	api_response, err := r.Client.CleanupPoliciesAPI.DeletePolicyByName(ctx, state.Name.ValueString()).Execute()

	// Handle Error(s)
	if err != nil {
		if api_response != nil && api_response.StatusCode == 404 {
			// Resource already deleted, nothing to do
			resp.Diagnostics.AddWarning(
				"Cleanup policy to delete did not exist",
				fmt.Sprintf("Cleanup policy was already deleted: %d: %s", api_response.StatusCode, api_response.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Deleting cleanup policy",
				fmt.Sprintf("Unable to delete cleanup policy: %s", err.Error()),
			)
		}
		return
	}

	if api_response.StatusCode != http.StatusNoContent && api_response.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Failed to delete cleanup policy",
			fmt.Sprintf("Unable to delete cleanup policy: %d: %s", api_response.StatusCode, api_response.Status),
		)
		return
	}
}

// ImportState imports the resource by name.
func (r *cleanupPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}