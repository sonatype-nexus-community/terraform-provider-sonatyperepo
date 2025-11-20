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

package system

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Constants for error messages to avoid duplication
const (
	errorSettingSecurityRealms           = "Error setting Security Realms configuration"
	errorSettingSecurityRealmsWithDetail = "Error setting Security Realms configuration: %s"
	unexpectedResponseCode               = "Unexpected Response Code whilst setting Security Realms configuration: %d: %s"
	gettingStateDataHasErrors            = "Getting state data has errors: %v"
	gettingPlanDataHasErrors             = "Getting plan data has errors: %v"
	gettingRequestDataHasErrors          = "Getting request data has errors: %v"
)

// securityRealmsResource is the resource implementation.
type securityRealmsResource struct {
	common.BaseResource
}

// SecurityRealmsConfiguration defines the structure for the API request body
type SecurityRealmsConfiguration struct {
	Active []string `json:"active"`
}

// NewSecurityRealmsResource is a helper function to simplify the provider implementation.
func NewSecurityRealmsResource() resource.Resource {
	return &securityRealmsResource{}
}

// Metadata returns the resource type name.
func (r *securityRealmsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_realms"
}

// Schema defines the schema for the resource.
func (r *securityRealmsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Activate and order Sontaype Nexus Repository Security realms. This resource manages the configuration of active security realms and their order.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Resource identifier",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"active": schema.ListAttribute{
				Description: "Specify active security realms in usage order. At least one realm must be specified.",
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.UniqueValues(),
				},
			},
		},
	}
}

// ImportState imports the resource state from the remote system.
// For security realms, we use a static ID since this is a singleton resource.
func (r *securityRealmsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	ctx = r.GetAuthContext(ctx)

	// Read current configuration from the API
	apiResponse, httpResponse, err := r.Client.SecurityManagementRealmsAPI.GetActiveRealms(ctx).Execute()
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			resp.Diagnostics.AddError(
				"Security Realms Configuration Not Found",
				"Unable to import Security Realms Configuration: configuration does not exist on the server.",
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading Security Realms Configuration",
				fmt.Sprintf("Unable to read Security Realms Configuration during import: %s", err),
			)
		}
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully imported security realms configuration: %v", apiResponse))

	// Convert API response to Terraform state
	var state model.SecurityRealmsModel

	// Set the ID
	state.ID = types.StringValue("security_realms")

	// Convert API response to Terraform List
	if apiResponse != nil {
		activeElements := make([]attr.Value, len(apiResponse))
		for i, realm := range apiResponse {
			activeElements[i] = types.StringValue(realm)
		}

		activeList, diags := types.ListValue(types.StringType, activeElements)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Active = activeList
	} else {
		// If no active realms, create empty list
		emptyList, diags := types.ListValue(types.StringType, []attr.Value{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Active = emptyList
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Info(ctx, "Successfully imported security realms configuration")
}

// Create creates the resource and sets the initial Terraform state.
func (r *securityRealmsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.SecurityRealmsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(gettingRequestDataHasErrors, resp.Diagnostics.Errors()))
		return
	}

	// Set up authentication context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Convert list to appropriate format
	activeRealms, err := r.convertActiveRealms(plan.Active)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid element type",
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating security realms configuration with realms: %v", activeRealms))

	// Call API to Create
	apiResponse, err := r.Client.SecurityManagementRealmsAPI.SetActiveRealms(ctx).Body(activeRealms).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			errorSettingSecurityRealms,
			fmt.Sprintf(errorSettingSecurityRealmsWithDetail, err.Error()),
		)
		return
	} else if apiResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			errorSettingSecurityRealms,
			fmt.Sprintf(unexpectedResponseCode, apiResponse.StatusCode, apiResponse.Status),
		)
	}

	tflog.Info(ctx, "Successfully created security realms configuration")

	// Set the ID - for a singleton resource, we can use a static ID
	plan.ID = types.StringValue("security_realms")

	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *securityRealmsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.SecurityRealmsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(gettingStateDataHasErrors, resp.Diagnostics.Errors()))
		return
	}

	ctx = r.GetAuthContext(ctx)

	// Read API Call - GetActiveRealms returns []string directly
	apiResponse, httpResponse, err := r.Client.SecurityManagementRealmsAPI.GetActiveRealms(ctx).Execute()
	if err != nil {
		r.handleReadError(ctx, resp, httpResponse, err)
		return
	}

	tflog.Debug(ctx, "Successfully read security realms configuration from API")

	// Convert API response and update state
	r.updateStateFromAPIResponse(apiResponse, &state, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure ID is set and save state
	r.finalizeReadState(ctx, &state, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *securityRealmsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.SecurityRealmsModel
	var state model.SecurityRealmsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(gettingPlanDataHasErrors, resp.Diagnostics.Errors()))
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(gettingStateDataHasErrors, resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Convert list to appropriate format
	activeRealms, err := r.convertActiveRealms(plan.Active)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid element type",
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating security realms configuration with realms: %v", activeRealms))

	// Call API to Update
	apiResponse, err := r.Client.SecurityManagementRealmsAPI.SetActiveRealms(ctx).Body(activeRealms).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			errorSettingSecurityRealms,
			fmt.Sprintf(errorSettingSecurityRealmsWithDetail, err.Error()),
		)
		return
	} else if apiResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			errorSettingSecurityRealms,
			fmt.Sprintf(unexpectedResponseCode, apiResponse.StatusCode, apiResponse.Status),
		)
	}

	tflog.Info(ctx, "Successfully updated security realms configuration")

	// Copy the ID from state to plan to maintain consistency
	plan.ID = state.ID
	if plan.ID.IsNull() || plan.ID.IsUnknown() {
		plan.ID = types.StringValue("security_realms")
	}

	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// Note: This resets the security realms to default configuration rather than
// truly deleting them, as Nexus always requires at least one security realm.
func (r *securityRealmsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.SecurityRealmsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(gettingStateDataHasErrors, resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Reset to default security realms - Nexus requires at least one realm
	defaultRealms := []string{"NexusAuthenticatingRealm"}

	tflog.Debug(ctx, fmt.Sprintf("Resetting security realms to default configuration: %v", defaultRealms))

	// Call API to reset to defaults
	apiResponse, err := r.Client.SecurityManagementRealmsAPI.SetActiveRealms(ctx).Body(defaultRealms).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			errorSettingSecurityRealms,
			fmt.Sprintf(errorSettingSecurityRealmsWithDetail, err.Error()),
		)
		return
	} else if apiResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			errorSettingSecurityRealms,
			fmt.Sprintf(unexpectedResponseCode, apiResponse.StatusCode, apiResponse.Status),
		)
	}

	tflog.Info(ctx, "Successfully reset security realms configuration to defaults")
}

// Helper functions

// convertActiveRealms converts a Terraform List to a slice of strings
func (r *securityRealmsResource) convertActiveRealms(active types.List) ([]string, error) {
	activeRealms := make([]string, len(active.Elements()))
	for i, elem := range active.Elements() {
		if strVal, ok := elem.(types.String); ok {
			activeRealms[i] = strVal.ValueString()
		} else {
			return nil, fmt.Errorf("expected string element at index %d", i)
		}
	}
	return activeRealms, nil
}

// handleReadError processes errors from the GetActiveRealms API call
func (r *securityRealmsResource) handleReadError(ctx context.Context, resp *resource.ReadResponse, httpResponse *http.Response, err error) {
	if httpResponse != nil && httpResponse.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddWarning(
			"Security Realms Configuration does not exist",
			fmt.Sprintf("Unable to read Security Realms Configuration: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}

	errorMsg := "Unable to read Security Realms Configuration"
	if httpResponse != nil {
		errorMsg = fmt.Sprintf("%s: %s", errorMsg, httpResponse.Status)
	}
	resp.Diagnostics.AddError(
		"Error Reading Security Realms Configuration",
		fmt.Sprintf("%s: %s", errorMsg, err),
	)
}

// updateStateFromAPIResponse converts API response to Terraform types and updates state
func (r *securityRealmsResource) updateStateFromAPIResponse(apiResponse []string, state *model.SecurityRealmsModel, resp *resource.ReadResponse) {
	if apiResponse == nil {
		return
	}

	activeElements := make([]attr.Value, len(apiResponse))
	for i, realm := range apiResponse {
		activeElements[i] = types.StringValue(realm)
	}

	activeList, diags := types.ListValue(types.StringType, activeElements)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Active = activeList
}

// finalizeReadState ensures ID is set and saves the final state
func (r *securityRealmsResource) finalizeReadState(ctx context.Context, state *model.SecurityRealmsModel, resp *resource.ReadResponse) {
	if state.ID.IsNull() || state.ID.IsUnknown() {
		state.ID = types.StringValue("security_realms")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
