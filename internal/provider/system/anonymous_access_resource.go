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
	"io"
	"net/http"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// anonymousAccessSystemResource is the resource implementation.
type anonymousAccessSystemResource struct {
	common.BaseResource
}

// NewAnonymousAccessSystemResource is a helper function to simplify the provider implementation.
func NewAnonymousAccessSystemResource() resource.Resource {
	return &anonymousAccessSystemResource{}
}

// Metadata returns the resource type name.
func (r *anonymousAccessSystemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_anonymous_access"
}

// Schema defines the schema for the resource.
func (r *anonymousAccessSystemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage Anonymous Access",
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: "Whether or not Anonymous Access is enabled",
				Required:    true,
			},
			"realm_name": schema.StringAttribute{
				Description: "The name of the authentication realm for the anonymous account",
				Required:    true,
			},
			"user_id": schema.StringAttribute{
				Description: "The username of the anonymous account",
				Required:    true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *anonymousAccessSystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.AnonymousAccessModel

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

	payload := sonatyperepo.AnonymousAccessSettingsXO{
		Enabled:   plan.Enabled.ValueBoolPointer(),
		RealmName: plan.RealmName.ValueStringPointer(),
		UserId:    plan.UserId.ValueStringPointer(),
	}

	_, httpResponse, err := r.Client.SecurityManagementAnonymousAccessAPI.Update1(ctx).Body(payload).Execute()

	// Handle Error
	if err != nil {
		if httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				"Your user is unauthorized to access this resource or feature.",
			)
		} else {
			errorBody, _ := io.ReadAll(httpResponse.Body)
			resp.Diagnostics.AddError(
				"Error updating Anonymous Access settings",
				"Could not update Anonymous Access settings, unexpected error: "+httpResponse.Status+": "+string(errorBody),
			)
		}
		return
	}

	if httpResponse.StatusCode == http.StatusOK {
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Failed to update Anonymous Access settings",
			fmt.Sprintf("Unable to update Anonymous Access settings: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *anonymousAccessSystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.AnonymousAccessModel

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
	apiResponse, httpResponse, err := r.Client.SecurityManagementAnonymousAccessAPI.Read1(ctx).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				"Your user is unauthorized to access this resource or feature.",
			)
		} else {
			errorBody, _ := io.ReadAll(httpResponse.Body)
			resp.Diagnostics.AddError(
				"Error reading Anonymous Access settings",
				"Could not read Anonymous Access settings, unexpected error: "+httpResponse.Status+": "+string(errorBody),
			)
		}
		return
	} else {
		state.Enabled = types.BoolValue(*apiResponse.Enabled)
		state.RealmName = types.StringValue(*apiResponse.RealmName)
		state.UserId = types.StringValue(*apiResponse.UserId)
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *anonymousAccessSystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.AnonymousAccessModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Update API Call
	payload := sonatyperepo.AnonymousAccessSettingsXO{
		Enabled:   plan.Enabled.ValueBoolPointer(),
		RealmName: plan.RealmName.ValueStringPointer(),
		UserId:    plan.UserId.ValueStringPointer(),
	}

	apiResponse, httpResponse, err := r.Client.SecurityManagementAnonymousAccessAPI.Update1(ctx).Body(payload).Execute()

	// Handle Error
	if err != nil {
		if httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				"Your user is unauthorized to access this resource or feature.",
			)
		} else {
			errorBody, _ := io.ReadAll(httpResponse.Body)
			resp.Diagnostics.AddError(
				"Error updating Anonymous Access settings",
				"Could not update Anonymous Access settings, unexpected error: "+httpResponse.Status+": "+string(errorBody),
			)
		}
		return
	}

	if httpResponse.StatusCode == http.StatusOK {
		plan.Enabled = types.BoolValue(*apiResponse.Enabled)
		plan.RealmName = types.StringValue(*apiResponse.RealmName)
		plan.UserId = types.StringValue(*apiResponse.UserId)
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Failed to update Anonymous Access settings",
			fmt.Sprintf("Unable to update Anonymous Access settings: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *anonymousAccessSystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.AnonymousAccessModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	//
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Update API Call
	payload := sonatyperepo.AnonymousAccessSettingsXO{
		Enabled:   common.NewFalse(),
		RealmName: common.StringPointer(common.DEFAULT_REALM_NAME),
		UserId:    common.StringPointer(common.DEFAULT_ANONYMOUS_USERNAME),
	}

	_, httpResponse, err := r.Client.SecurityManagementAnonymousAccessAPI.Update1(ctx).Body(payload).Execute()

	// Handle Error
	if err != nil {
		if httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				"Your user is unauthorized to access this resource or feature.",
			)
		} else {
			errorBody, _ := io.ReadAll(httpResponse.Body)
			resp.Diagnostics.AddError(
				"Error removing Anonymous Access settings",
				"Could not remove Anonymous Access settings, unexpected error: "+httpResponse.Status+": "+string(errorBody),
			)
		}
		return
	}
}
