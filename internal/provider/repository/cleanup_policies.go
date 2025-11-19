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
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	tfschema "github.com/sonatype-nexus-community/terraform-provider-shared/schema"
	"io"
	"net/http"
	"regexp"
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

const cleanupPolicyNamePattern = `^[a-zA-Z0-9\-]{1}[a-zA-Z0-9_\-\.]*$`

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
			"name": func() schema.StringAttribute {
				attr := tfschema.RequiredStringWithRegexAndLength(
					"Name of the cleanup policy",
					regexp.MustCompile(cleanupPolicyNamePattern),
					"Name must start with an alphanumeric character or hyphen, and can only contain alphanumeric characters, underscores, hyphens, and periods",
					1,
					255,
				)
				attr.PlanModifiers = []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				}
				return attr
			}(),
			"notes": tfschema.OptionalString("Notes for the cleanup policy"),
			"format": tfschema.RequiredStringEnum(
				"Repository format that this cleanup policy applies to",
				"apt", "bower", "cocoapods", "conan", "conda", "docker", "gitlfs", "go", "helm", "maven2", "npm", "nuget", "p2", "pypi", "r", "raw", "rubygems", "yum",
			),
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

	requestPayload := buildRequestPayload(plan)

	apiResponse, err := r.Client.CleanupPoliciesAPI.Create1(ctx).Body(requestPayload).Execute()

	// Handle Error
	if err != nil {
		sharederr.HandleAPIError(
			"Error creating cleanup policy",
			&err,
			apiResponse,
			&resp.Diagnostics,
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
		sharederr.HandleAPIError(
			"Creation of cleanup policy was not successful",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
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

	ctx = r.GetAuthContext(ctx)

	// Fetch cleanup policy from API
	cleanupPolicy, err := r.fetchCleanupPolicy(ctx, state.Name.ValueString())
	if err != nil {
		// Check if this is a 404 error by attempting to get the HTTP response
		httpResponse, _ := r.Client.CleanupPoliciesAPI.GetCleanupPolicyByName(ctx, state.Name.ValueString()).Execute()
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				"Cleanup policy to read did not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
			return
		}

		sharederr.HandleAPIError(
			"Error reading cleanup policy",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	// Update state from API response
	updateStateFromAPI(&state, *cleanupPolicy)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *cleanupPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan, state model.CleanupPolicyModel

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
		resp.Diagnostics.AddError("Invalid cleanup policy configuration", err.Error())
		return
	}

	ctx = r.GetAuthContext(ctx)

	// Build request payload and make API call
	requestPayload := buildRequestPayload(plan)
	apiRequest := r.Client.CleanupPoliciesAPI.Update2(ctx, state.Name.ValueString()).Body(requestPayload)
	apiResponse, err := apiRequest.Execute()

	// Handle API response
	if err != nil {
		r.handleUpdateError(resp, apiResponse, err)
		return
	}

	if apiResponse.StatusCode == http.StatusNoContent || apiResponse.StatusCode == http.StatusOK {
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
	apiResponse, err := r.Client.CleanupPoliciesAPI.DeletePolicyByName(ctx, state.Name.ValueString()).Execute()

	// Handle Error(s)
	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			// Resource already deleted, nothing to do
			sharederr.HandleAPIWarning(
				"Cleanup policy to delete did not exist",
				&err,
				apiResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				"Error deleting cleanup policy",
				&err,
				apiResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	if apiResponse.StatusCode != http.StatusNoContent && apiResponse.StatusCode != http.StatusOK {
		sharederr.HandleAPIError(
			"Deletion of cleanup policy was not successful",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
		return
	}
}

// ImportState imports the resource by name.
func (r *cleanupPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Helper functions

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

// buildRequestPayload creates the API request payload from the model
func buildRequestPayload(plan model.CleanupPolicyModel) sonatyperepo.CleanupPolicyResourceXO {
	requestPayload := sonatyperepo.CleanupPolicyResourceXO{
		Name:   plan.Name.ValueString(),
		Format: plan.Format.ValueString(),
	}

	if !plan.Notes.IsNull() {
		notes := plan.Notes.ValueStringPointer()
		requestPayload.Notes = notes
	}

	setCriteriaFields(&requestPayload, plan.Criteria)

	if !plan.Retain.IsNull() {
		retain := int32(plan.Retain.ValueInt64())
		requestPayload.Retain = &retain
	}

	return requestPayload
}

// setCriteriaFields sets the criteria fields in the request payload
func setCriteriaFields(payload *sonatyperepo.CleanupPolicyResourceXO, criteria *model.CleanupPolicyCriteriaModel) {
	if !criteria.LastBlobUpdated.IsNull() {
		lastBlobUpdated := criteria.LastBlobUpdated.ValueInt64()
		payload.CriteriaLastBlobUpdated = &lastBlobUpdated
	}

	if !criteria.LastDownloaded.IsNull() {
		lastDownloaded := criteria.LastDownloaded.ValueInt64()
		payload.CriteriaLastDownloaded = &lastDownloaded
	}

	if !criteria.ReleaseType.IsNull() {
		releaseType := criteria.ReleaseType.ValueString()
		payload.CriteriaReleaseType = &releaseType
	}

	if !criteria.AssetRegex.IsNull() {
		assetRegex := criteria.AssetRegex.ValueString()
		payload.CriteriaAssetRegex = &assetRegex
	}
}

// updateStateFromAPI updates the state model from the API response
func updateStateFromAPI(state *model.CleanupPolicyModel, cleanupPolicy sonatyperepo.CleanupPolicyResourceXO) {
	state.Name = types.StringValue(cleanupPolicy.Name)
	state.Format = types.StringValue(cleanupPolicy.Format)

	if cleanupPolicy.Notes != nil {
		state.Notes = types.StringValue(*cleanupPolicy.Notes)
	} else {
		state.Notes = types.StringNull()
	}

	// Ensure we always have a criteria object since it's required
	if state.Criteria == nil {
		state.Criteria = &model.CleanupPolicyCriteriaModel{}
	}

	updateCriteriaFromAPI(state.Criteria, cleanupPolicy)

	if cleanupPolicy.Retain != nil {
		state.Retain = types.Int64Value(int64(*cleanupPolicy.Retain))
	} else {
		state.Retain = types.Int64Null()
	}
}

// updateCriteriaFromAPI updates the criteria model from the API response
func updateCriteriaFromAPI(criteria *model.CleanupPolicyCriteriaModel, cleanupPolicy sonatyperepo.CleanupPolicyResourceXO) {
	if cleanupPolicy.CriteriaLastBlobUpdated != nil {
		criteria.LastBlobUpdated = types.Int64Value(*cleanupPolicy.CriteriaLastBlobUpdated)
	} else {
		criteria.LastBlobUpdated = types.Int64Null()
	}

	if cleanupPolicy.CriteriaLastDownloaded != nil {
		criteria.LastDownloaded = types.Int64Value(*cleanupPolicy.CriteriaLastDownloaded)
	} else {
		criteria.LastDownloaded = types.Int64Null()
	}

	if cleanupPolicy.CriteriaReleaseType != nil {
		criteria.ReleaseType = types.StringValue(*cleanupPolicy.CriteriaReleaseType)
	} else {
		criteria.ReleaseType = types.StringNull()
	}

	if cleanupPolicy.CriteriaAssetRegex != nil {
		criteria.AssetRegex = types.StringValue(*cleanupPolicy.CriteriaAssetRegex)
	} else {
		criteria.AssetRegex = types.StringNull()
	}
}

// fetchCleanupPolicy retrieves and parses the cleanup policy from the API
func (r *cleanupPolicyResource) fetchCleanupPolicy(ctx context.Context, name string) (*sonatyperepo.CleanupPolicyResourceXO, error) {
	httpResponse, err := r.Client.CleanupPoliciesAPI.GetCleanupPolicyByName(ctx, name).Execute()
	if err != nil {
		return nil, err
	}

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	var cleanupPolicy sonatyperepo.CleanupPolicyResourceXO
	err = json.Unmarshal(responseBody, &cleanupPolicy)
	if err != nil {
		return nil, fmt.Errorf("could not parse response: %w", err)
	}

	return &cleanupPolicy, nil
}

// handleUpdateError handles errors from the update API call
func (r *cleanupPolicyResource) handleUpdateError(resp *resource.UpdateResponse, apiResponse *http.Response, err error) {
	if apiResponse != nil && apiResponse.StatusCode == 404 {
		resp.State.RemoveResource(context.Background())
		sharederr.HandleAPIWarning(
			"Cleanup policy to update did not exist",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
	} else {
		sharederr.HandleAPIError(
			"Error updating cleanup policy",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
	}
}
