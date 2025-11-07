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

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// securityUserTokenResource is the resource implementation.
type securityUserTokenResource struct {
	common.BaseResource
}

// NewSecurityUserTokenResource is a helper function to simplify the provider implementation.
func NewSecurityUserTokenResource() resource.Resource {
	return &securityUserTokenResource{}
}

// Metadata returns the resource type name.
func (r *securityUserTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_user_token"
}

// Schema defines the schema for the resource.
func (r *securityUserTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage User Token Configuration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Resource identifier (always 'user-tokens' for this singleton resource)",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether or not User Tokens feature is enabled",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"expiration_days": schema.Int32Attribute{
				Description: "Set user token expiration days (1-999)",
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(1),
				Validators: []validator.Int32{
					int32validator.Between(1, 999),
				},
			},
			"expiration_enabled": schema.BoolAttribute{
				Description: "Enable user tokens expiration",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"protect_content": schema.BoolAttribute{
				Description: "Additionally require user tokens for repository authentication",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// ImportState imports the resource into Terraform state.
func (r *securityUserTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Since this is a singleton resource (there's only one user token configuration),
	// we don't need to parse the import ID. We just read the current configuration.

	// Set up authentication context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Read current user token settings from the API
	apiResponse, httpResponse, err := r.Client.SecurityManagementUserTokensAPI.ServiceStatus(ctx).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				"Your user is unauthorized to access this resource or feature during import.",
			)
		} else {
			errorBody, _ := io.ReadAll(httpResponse.Body)
			resp.Diagnostics.AddError(
				"Error importing User Token settings",
				"Could not read User Token settings during import, unexpected error: "+httpResponse.Status+": "+string(errorBody),
			)
		}
		return
	}

	// Create the state model with the current API values
	var state model.SecurityUserTokenModel
	state.ID = types.StringValue("user-tokens")
	state.MapFromApi(apiResponse)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully imported security user token resource")
}

// Create creates the resource and sets the initial Terraform state.
func (r *securityUserTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.SecurityUserTokenModel

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

	payload := sonatyperepo.UserTokensApiModel{}
	plan.MapToApi(&payload)

	apiResponse, httpResponse, err := r.Client.SecurityManagementUserTokensAPI.SetServiceStatus(ctx).Body(payload).Execute()

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
				"Error creating User Token settings",
				"Could not create User Token settings, unexpected error: "+httpResponse.Status+": "+string(errorBody),
			)
		}
		return
	}

	if httpResponse.StatusCode == http.StatusOK {
		plan.ID = types.StringValue("user-tokens")
		plan.MapFromApi(apiResponse)
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Failed to create User Token settings",
			fmt.Sprintf("Unable to create User Token settings: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *securityUserTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.SecurityUserTokenModel

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
	apiResponse, httpResponse, err := r.Client.SecurityManagementUserTokensAPI.ServiceStatus(ctx).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				"Your user is unauthorized to access this resource or feature.",
			)
		} else {
			errorBody, _ := io.ReadAll(httpResponse.Body)
			resp.Diagnostics.AddError(
				"Error reading User Token settings",
				"Could not read User Token settings, unexpected error: "+httpResponse.Status+": "+string(errorBody),
			)
		}
		return
	} else {
		state.ID = types.StringValue("user-tokens")
		state.MapFromApi(apiResponse)
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *securityUserTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan model.SecurityUserTokenModel

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
	payload := sonatyperepo.UserTokensApiModel{}
	plan.MapToApi(&payload)

	apiResponse, httpResponse, err := r.Client.SecurityManagementUserTokensAPI.SetServiceStatus(ctx).Body(payload).Execute()

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
				"Error updating User Token settings",
				"Could not update User Token settings, unexpected error: "+httpResponse.Status+": "+string(errorBody),
			)
		}
		return
	}

	if httpResponse.StatusCode == http.StatusOK {
		plan.ID = types.StringValue("user-tokens")
		plan.MapFromApi(apiResponse)
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Failed to update User Token settings",
			fmt.Sprintf("Unable to update User Token settings: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *securityUserTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.SecurityUserTokenModel

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

	// Instead of deleting, we disable the user token feature
	defaultExpirationDays := int32(1)
	payload := sonatyperepo.UserTokensApiModel{
		Enabled:        common.NewFalse(),
		ExpirationDays: &defaultExpirationDays,
	}

	_, httpResponse, err := r.Client.SecurityManagementUserTokensAPI.SetServiceStatus(ctx).Body(payload).Execute()

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
				"Error disabling User Token settings",
				"Could not disable User Token settings, unexpected error: "+httpResponse.Status+": "+string(errorBody),
			)
		}
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Failed to disable User Token settings",
			fmt.Sprintf("Unable to disable User Token settings: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
	}
}
