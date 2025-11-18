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
	"fmt"
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
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

const (
	// RoutingRuleModeAllow represents the ALLOW mode for routing rules
	RoutingRuleModeAllow = "ALLOW"
	// RoutingRuleModeBlock represents the BLOCK mode for routing rules
	RoutingRuleModeBlock = "BLOCK"
)

// routingRuleResource is the resource implementation.
type routingRuleResource struct {
	common.BaseResource
}

// NewRoutingRuleResource is a helper function to simplify the provider implementation.
func NewRoutingRuleResource() resource.Resource {
	return &routingRuleResource{}
}

// Metadata returns the resource type name.
func (r *routingRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_rule"
}

// Schema defines the schema for the resource.
func (r *routingRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this resource to create and manage routing rules in Sonatype Nexus Repository Manager",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the routing rule",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9\-]{1}[a-zA-Z0-9_\-\.]*$`),
						"Name must start with an alphanumeric character or hyphen, and can only contain alphanumeric characters, underscores, hyphens, and periods",
					),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the routing rule (required by Nexus API)",
				Required:    true,
			},
			"mode": schema.StringAttribute{
				Description: "Determines what should be done with requests when their path matches any of the matchers. Valid values: ALLOW, BLOCK",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(RoutingRuleModeAllow, RoutingRuleModeBlock),
				},
			},
			"matchers": schema.SetAttribute{
				Description: "Regular expressions used to identify request paths that are allowed or blocked (depending on mode)",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *routingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.RoutingRuleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Call API to Create
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	requestPayload := sonatyperepo.RoutingRuleXO{}
	plan.MapToApi(&requestPayload)

	apiResponse, err := r.Client.RoutingRulesAPI.CreateRoutingRule(ctx).Body(requestPayload).Execute()

	// Handle Error
	if err != nil {
		sharederr.HandleAPIError(
			"Error creating routing rule",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
		return
	}

	if apiResponse.StatusCode == http.StatusNoContent {
		// Set LastUpdated
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		sharederr.HandleAPIError(
			"Failed to create routing rule",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *routingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.RoutingRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.GetAuthContext(ctx)

	// Fetch routing rule from API
	routingRule, httpResponse, err := r.Client.RoutingRulesAPI.GetRoutingRule(ctx, state.Name.ValueString()).Execute()
	if err != nil {
		// Check if this is a 404 error
		if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				"Routing rule not found",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
			return
		}

		sharederr.HandleAPIError(
			"Error reading routing rule",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	// Update state from API response
	state.MapFromApi(routingRule)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *routingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan, state model.RoutingRuleModel

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

	// Build request payload and make API call
	requestPayload := sonatyperepo.RoutingRuleXO{}
	plan.MapToApi(&requestPayload)
	apiResponse, err := r.Client.RoutingRulesAPI.UpdateRoutingRule(ctx, state.Name.ValueString()).Body(requestPayload).Execute()

	// Handle API response
	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				"Routing rule to update did not exist",
				&err,
				apiResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				"Error updating routing rule",
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

// Delete deletes the resource and removes the Terraform state on success.
func (r *routingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.RoutingRuleModel

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
	apiResponse, err := r.Client.RoutingRulesAPI.DeleteRoutingRule(ctx, state.Name.ValueString()).Execute()

	// Handle Error(s)
	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == http.StatusNotFound {
			// Resource already deleted, nothing to do
			sharederr.HandleAPIWarning(
				"Routing rule to delete did not exist",
				&err,
				apiResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				"Error deleting routing rule",
				&err,
				apiResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	if apiResponse.StatusCode != http.StatusNoContent && apiResponse.StatusCode != http.StatusOK {
		sharederr.HandleAPIError(
			"Failed to delete routing rule",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
		return
	}
}

// ImportState imports the resource by name.
func (r *routingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
