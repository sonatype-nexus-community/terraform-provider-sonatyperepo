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

package user

import (
	"context"
	"fmt"
	"net/http"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// userResource is the resource implementation.
type userResource struct {
	common.BaseResource
}

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configure the Sonatype IQ Server Connection",
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				Description: "The userid which is required for login. This value cannot be changed.",
				Required:    true,
				Optional:    false,
			},
			"first_name": schema.StringAttribute{
				Description: "The first name of the user.",
				Required:    true,
				Optional:    false,
			},
			"last_name": schema.StringAttribute{
				Description: "The last name of the user.",
				Required:    true,
				Optional:    false,
			},
			"email_address": schema.StringAttribute{
				Description: "The email address associated with the user.",
				Required:    true,
				Optional:    false,
			},
			"password": schema.StringAttribute{
				Description: "The password for the user.",
				Required:    true,
				Optional:    false,
				Sensitive:   true,
			},
			"status": schema.StringAttribute{
				Description: "The user's status.",
				Required:    true,
				Optional:    false,
				Validators: []validator.String{
					stringvalidator.OneOf(
						common.USER_STATUS_ACTIVE,
						common.USER_STATUS_LOCKED,
						common.USER_STATUS_DISABLED,
						common.USER_STATUS_CHANGE_PASSWORD,
					),
				},
			},
			"roles": schema.ListAttribute{
				Description: "The list of roles assigned to this User.",
				Required:    true,
				Optional:    false,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
				},
			},
			"read_only": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"source": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.UserModelResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

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
	apiBody := sonatyperepo.NewApiCreateUser(plan.Status.ValueString())
	plan.MapToCreateApi(apiBody)
	apiResponse, httpResponse, err := r.Client.SecurityManagementUsersAPI.CreateUser(ctx).Body(*apiBody).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating User",
			fmt.Sprintf("Error creating User: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	} else if httpResponse.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error creating User",
			fmt.Sprintf("Unexpected Response Code whilst creating User: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
	}

	// Update State
	plan.MapFromApi(apiResponse)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.UserModelResource
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
	apiResponse, httpResponse, err := r.Client.SecurityManagementUsersAPI.GetUsers(ctx).UserId(state.UserId.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"User does not exist",
				fmt.Sprintf("Unable to read User: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading User",
				fmt.Sprintf("Unable to read User: %s: %s", httpResponse.Status, err),
			)
		}
		return
	}

	if len(apiResponse) == 0 {
		resp.Diagnostics.AddError(
			"No User with requested User ID",
			fmt.Sprintf("No user returned: %s: %s", httpResponse.Status, err),
		)
		return
	}

	// Update State based on Response
	state.MapFromApi(&apiResponse[0])
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.UserModelResource
	var state model.UserModelResource

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
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)
	apiBody := sonatyperepo.NewApiUser(plan.Status.ValueString())
	plan.MapToApi(apiBody)
	httpResponse, err := r.Client.SecurityManagementUsersAPI.UpdateUser(ctx, state.UserId.ValueString()).Body(*apiBody).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating User",
			fmt.Sprintf("Error updating User: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error updating User",
			fmt.Sprintf("Unexpected Response Code whilst updating User: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
	}

	// If Passowrd is required to be changed, make that additional API call now
	if !plan.Password.Equal(state.Password) {
		httpResponse, err = r.Client.SecurityManagementUsersAPI.ChangePassword(ctx, state.UserId.ValueString()).Body(plan.Password.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating User password",
				fmt.Sprintf("Error updating User password: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
			return
		} else if httpResponse.StatusCode != http.StatusNoContent {
			resp.Diagnostics.AddError(
				"Error updating User password",
				fmt.Sprintf("Unexpected Response Code whilst updating User password: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		}
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.UserModelResource

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

	httpResponse, err := r.Client.SecurityManagementUsersAPI.DeleteUser(ctx, state.UserId.ValueString()).Execute()

	// Handle Error
	if err != nil {
		resp.Diagnostics.AddError(
			"Error removing User",
			fmt.Sprintf("Error removing User: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error removing User",
			fmt.Sprintf("Unexpected Response Code whilst removing User: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
	}
}
