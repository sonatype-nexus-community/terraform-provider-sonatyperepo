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
			"enabled":      tfschema.ResourceRequiredBool("Whether or not Anonymous Access is enabled"),
			"realm_name":   tfschema.ResourceRequiredString("The name of the authentication realm for the anonymous account"),
			"user_id":      tfschema.ResourceRequiredString("The username of the anonymous account"),
			"last_updated": tfschema.ResourceComputedString("The timestamp of when the resource was last updated"),
		},
	}
}

// ImportState imports the resource into Terraform state.
func (r *anonymousAccessSystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Since this is a singleton resource (there's only one anonymous access configuration),
	// we don't need to parse the import ID. We just read the current configuration.

	// Set up authentication context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Read current anonymous access settings from the API
	apiResponse, httpResponse, err := r.Client.SecurityManagementAnonymousAccessAPI.Read1(ctx).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				"Your user is unauthorized to access this resource or feature during import.",
			)
		} else {
			sharederr.HandleAPIError(
				"Error importing Anonymous Access settings",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Create the state model with the current API values
	var state model.AnonymousAccessModel
	state.Enabled = types.BoolValue(*apiResponse.Enabled)
	state.RealmName = types.StringValue(*apiResponse.RealmName)
	state.UserId = types.StringValue(*apiResponse.UserId)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully imported anonymous access system resource")
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
			sharederr.HandleAPIError(
				"Error updating Anonymous Access settings",
				&err,
				httpResponse,
				&resp.Diagnostics,
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
		sharederr.HandleAPIError(
			"Update of Anonymous Access settings was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
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
			sharederr.HandleAPIError(
				"Error reading Anonymous Access settings",
				&err,
				httpResponse,
				&resp.Diagnostics,
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
			sharederr.HandleAPIError(
				"Error updating Anonymous Access settings",
				&err,
				httpResponse,
				&resp.Diagnostics,
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
			sharederr.HandleAPIError(
				"Error removing Anonymous Access settings",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}
}
