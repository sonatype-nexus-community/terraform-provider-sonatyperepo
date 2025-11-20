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

package role

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	tfschema "github.com/sonatype-nexus-community/terraform-provider-shared/schema"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
)

// roleResource is the resource implementation.
type roleResource struct {
	common.BaseResource
}

// NewRoleResource is a helper function to simplify the provider implementation.
func NewRoleResource() resource.Resource {
	return &roleResource{}
}

// Metadata returns the resource type name.
func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the resource.
func (r *roleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	idAttr := tfschema.ResourceRequiredString("The id of the Role.\n\nThis should be unique and can be the name of an LDAP or SAML Group if you are using LDAP or SAML for authentication.\nMatching Roles based on id will automatically be granted to LDAP or SAML users.")
	idAttr.PlanModifiers = []planmodifier.String{
		stringplanmodifier.RequiresReplace(),
	}

	resp.Schema = schema.Schema{
		Description: "Manage Roles in Sonatype Nexus Repository",
		Attributes: map[string]schema.Attribute{
			"id":          idAttr,
			"name":        tfschema.ResourceRequiredString("The name of the role."),
			"description": tfschema.ResourceRequiredString("The description of this role."),
			"privileges":  tfschema.ResourceRequiredStringSet("The set of privileges assigned to this role."),
			"roles":       tfschema.ResourceRequiredStringSet("The set of roles assigned to this role."),
			"last_updated": tfschema.ResourceComputedString("The timestamp of when the resource was last updated"),
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.RoleModelResource
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
	apiBody := sonatyperepo.NewRoleXORequest()
	plan.MapToApi(apiBody)
	_, httpResponse, err := r.Client.SecurityManagementRolesAPI.Create(ctx).Body(*apiBody).Execute()

	if err != nil {
		sharederr.HandleAPIError(
			"Error creating Role",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusOK {
		sharederr.HandleAPIError(
			"Creation of Role was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}

	// Update State
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.RoleModelResource
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
	apiResponse, httpResponse, err := r.Client.SecurityManagementRolesAPI.GetRole(ctx, state.Id.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				"Role to read did not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				"Error reading Role",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Update State based on Response
	state.MapFromApi(apiResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.RoleModelResource
	var state model.RoleModelResource

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
	apiBody := sonatyperepo.NewRoleXORequest()
	plan.MapToApi(apiBody)
	httpResponse, err := r.Client.SecurityManagementRolesAPI.Update(ctx, state.Id.ValueString()).Body(*apiBody).Execute()

	if err != nil {
		sharederr.HandleAPIError(
			"Error updating Role",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		sharederr.HandleAPIError(
			"Update of Role was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.RoleModelResource

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

	httpResponse, err := r.Client.SecurityManagementRolesAPI.Delete(ctx, state.Id.ValueString()).Execute()

	// Handle Error
	if err != nil {
		sharederr.HandleAPIError(
			"Error removing Role",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		sharederr.HandleAPIError(
			"Removal of Role was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}
}

// ImportState imports the resource into Terraform state.
func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID - the import ID should be the role ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
