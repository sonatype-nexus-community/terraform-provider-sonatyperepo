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
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"net/http"
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
		Description: "Manage Local and non-Local Users",
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				Description: "The userid which is required for login. This value cannot be changed.",
				Required:    true,
				Optional:    false,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: `The first name of the user.
				
**Note:** This can only be managed for local users - and not LDAP, CROWD or SAML users.`,
				Required: true,
				Optional: false,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: `The last name of the user.

**Note:** This can only be managed for local users - and not LDAP, CROWD or SAML users.`,
				Required: true,
				Optional: false,
			},
			"email_address": schema.StringAttribute{
				MarkdownDescription: `The email address associated with the user.
				
**Note:** This can only be managed for local users - and not LDAP, CROWD or SAML users.`,
				Required: true,
				Optional: false,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: `The password for the user.
				
**Note:** This is required for LOCAL users and must not be supplied for LDAP, CROWD or SAML users.`,
				Required:  false,
				Optional:  true,
				Sensitive: true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: `The user's status.
				
**Note:** This can only be managed for local users - and not LDAP, CROWD or SAML users.`,
				Required: true,
				Optional: false,
				Validators: []validator.String{
					stringvalidator.OneOf(
						common.USER_STATUS_ACTIVE,
						common.USER_STATUS_LOCKED,
						common.USER_STATUS_DISABLED,
						common.USER_STATUS_CHANGE_PASSWORD,
					),
				},
			},
			"roles": schema.SetAttribute{
				Description: "The list of roles assigned to this User.",
				Required:    true,
				Optional:    false,
				ElementType: types.StringType,
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

	if err != nil || httpResponse.StatusCode != http.StatusOK {
		sharederr.HandleAPIError(
			fmt.Sprintf("Error creating User: %s", plan.UserId.ValueString()),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
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
	apiResponse, httpResponse, err := r.Client.SecurityManagementUsersAPI.GetUsers(ctx).UserId(state.UserId.ValueString()).Source(state.Source.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				"User to read did not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				"Error reading User",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	if len(apiResponse) == 0 {
		resp.Diagnostics.AddError(
			"No User with requested User ID",
			fmt.Sprintf("No user found for %s@%s", state.UserId.ValueString(), state.Source.ValueString()),
		)
		return
	}

	var actualUser *sonatyperepo.ApiUser
	for _, u := range apiResponse {
		if *u.UserId == state.UserId.ValueString() && *u.Source == state.Source.ValueString() {
			tflog.Debug(ctx,
				fmt.Sprintf(
					"Matched User: %s=%s and %s=%s",
					*u.UserId,
					state.UserId.ValueString(),
					*u.Source,
					state.Source.ValueString(),
				),
			)
			actualUser = &u
		}
	}

	if actualUser == nil {
		// No user with the exact User ID and Source
		resp.Diagnostics.AddError(
			"User does not exist",
			fmt.Sprintf("No user returned: %s: %s", httpResponse.Status, err),
		)
		return
	}

	// Update State based on Response
	state.MapFromApi(actualUser)
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

// This allows users to import existing Users into Terraform state using the blob store name as the identifier.
func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <user_id>,<source> - e.g. admin,SAML. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source"), idParts[1])...)
}
