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
	"strings"
	"time"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
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
	resp.Schema = tfschema.Schema{
		Description: "Manage Local and non-Local Users",
		Attributes: map[string]tfschema.Attribute{
			"user_id": schema.ResourceRequiredStringWithPlanModifier(
				"The userid which is required for login. This value cannot be changed.",
				[]planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			),
			"first_name": schema.ResourceRequiredString(`The first name of the user.

  **Note:** This can only be managed for local users - and not LDAP, CROWD or SAML users.`),
			"last_name": schema.ResourceRequiredString(`The last name of the user.

  **Note:** This can only be managed for local users - and not LDAP, CROWD or SAML users.`),
			"email_address": schema.ResourceRequiredString(`The email address associated with the user.

  **Note:** This can only be managed for local users - and not LDAP, CROWD or SAML users.`),
			"password": schema.ResourceSensitiveString(`The password for the user.
			
  **Note:** This is required for LOCAL users and must not be supplied for LDAP, CROWD or SAML users.`),
			"status": schema.ResourceRequiredStringWithValidators(
				`The user's status.
				
  **Note:** This can only be managed for local users - and not LDAP, CROWD or SAML users.`,
				stringvalidator.OneOf(
					common.AllUserStatusTypes()...,
				),
			),
			"roles":        schema.ResourceRequiredStringSet("The list of roles assigned to this User."),
			"read_only":    schema.ResourceComputedBool("Whether the user is read-only"),
			"source":       schema.ResourceComputedString("Source system managing this user"),
			"last_updated": schema.ResourceLastUpdated(),
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
		errors.HandleAPIError(
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
			errors.HandleAPIWarning(
				"User to read did not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
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
