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
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

const (
	serviceAccountTokenNamePattern = `^[a-zA-Z0-9_-]+$`
	serviceAccountTokenNameMax     = 128
	serviceAccountTokenDescMax     = 256
)

// securityServiceAccountTokenResource is the resource implementation.
type securityServiceAccountTokenResource struct {
	common.BaseResource
}

// NewSecurityServiceAccountTokenResource is a helper function to simplify the provider implementation.
func NewSecurityServiceAccountTokenResource() resource.Resource {
	return &securityServiceAccountTokenResource{}
}

// Metadata returns the resource type name.
func (r *securityServiceAccountTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_service_account_token"
}

// Schema defines the schema for the resource.
func (r *securityServiceAccountTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	replaceString := []planmodifier.String{stringplanmodifier.RequiresReplace()}
	keepString := stringplanmodifier.UseStateForUnknown()

	name := schema.ResourceRequiredStringWithPlanModifier(
		fmt.Sprintf("Name of the Service Account. Up to %d characters: letters, numbers, hyphens and underscores. Must be unique. Changing this replaces the token.", serviceAccountTokenNameMax),
		replaceString,
	)
	name.Validators = []validator.String{
		stringvalidator.LengthBetween(1, serviceAccountTokenNameMax),
		stringvalidator.RegexMatches(regexp.MustCompile(serviceAccountTokenNamePattern), "must contain only letters, numbers, hyphens and underscores"),
	}

	roleId := schema.ResourceRequiredStringWithPlanModifier(
		"Id of the Role that determines the permissions of the token - e.g. the `id` of a `sonatyperepo_role`. Changing this replaces the token.",
		replaceString,
	)
	roleId.Validators = []validator.String{stringvalidator.LengthAtLeast(1)}

	description := schema.ResourceOptionalStringWithPlanModifier(
		fmt.Sprintf("Description of the Service Account, up to %d characters. Changing this replaces the token.", serviceAccountTokenDescMax),
		stringplanmodifier.RequiresReplace(),
	)
	description.Validators = []validator.String{stringvalidator.LengthAtMost(serviceAccountTokenDescMax)}

	expirationDays := schema.ResourceOptionalInt64(
		"Number of days after which the token expires. Omit for a token that never expires. Changing this replaces the token.",
	)
	expirationDays.PlanModifiers = []planmodifier.Int64{int64planmodifier.RequiresReplace()}
	expirationDays.Validators = []validator.Int64{int64validator.AtLeast(1)}

	token := schema.ResourceComputedSensitiveString(
		"The token value, starting `sat.`. Nexus Repository returns this only once, when the token is created; it cannot be retrieved afterwards. It is stored in the Terraform state.",
	)
	token.PlanModifiers = []planmodifier.String{keepString}

	resp.Schema = tfschema.Schema{
		MarkdownDescription: "Manage a Sonatype Nexus Repository Service Account Token: a named, role-scoped bearer token for " +
			"applications and automation.\n\n" +
			"**Requires Nexus Repository Pro 3.97.0 or later.** The underlying API does not exist on earlier versions - applying " +
			"this resource against an older server fails with a clear error.\n\n" +
			"~> **Immutable.** Nexus Repository cannot update a Service Account Token, so changing any argument replaces it " +
			"and issues a new token.\n\n" +
			"~> **The token cannot be read back.** Its value is only returned when the token is created and is held in the " +
			"Terraform state (marked sensitive). For this reason this resource **cannot be imported**.\n\n" +
			"A token only authenticates while the `ServiceAccountRealm` security realm is active - see " +
			"`sonatyperepo_security_realms`. Creating, listing and deleting tokens does not require it.",
		Attributes: map[string]tfschema.Attribute{
			"id": schema.ResourceComputedStringWithPlanModifier(
				"Id of the Service Account Token, assigned by Nexus Repository",
				keepString,
			),
			"name":            name,
			"role_id":         roleId,
			"description":     description,
			"expiration_days": expirationDays,
			"token":           token,
			"expires_at": schema.ResourceComputedStringWithPlanModifier(
				"Time at which the token expires (RFC 3339). Not set for a token that never expires.",
				keepString,
			),
			"created_at": schema.ResourceComputedStringWithPlanModifier(
				"Time at which the token was created (RFC 3339)",
				keepString,
			),
			"created_by": schema.ResourceComputedStringWithPlanModifier(
				"User that created the token",
				keepString,
			),
		},
	}
}

// ImportState is deliberately unsupported. The common.BaseResource this resource embeds
// defines ImportState (as a panic), which would make the framework offer `terraform import`;
// override it with a clear error instead.
func (r *securityServiceAccountTokenResource) ImportState(_ context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError(
		"Import is not supported",
		"A Service Account Token's value is returned by Nexus Repository only when the token is created and "+
			"cannot be retrieved afterwards, so an imported token would have no usable value. Create the token with "+
			"Terraform instead.",
	)
}

// Create creates the resource and sets the initial Terraform state.
func (r *securityServiceAccountTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.SecurityServiceAccountTokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	payload := sonatyperepoV397.ServiceAccountCreateXO{}
	plan.MapToApi(&payload)

	tflog.Debug(ctx, fmt.Sprintf("Creating Service Account Token %q", payload.Name))

	created, httpResponse, err := r.Services.ServiceAccountTokens.CreateServiceAccountToken(ctx, payload)
	if err != nil {
		errors.HandleAPIError(
			"Error creating Service Account Token",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Error creating Service Account Token",
			fmt.Sprintf("Unexpected Response Code whilst creating Service Account Token: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}

	plan.MapFromCreate(created)

	// The create response does not include who created the token or when. Complete the state
	// from the list; if that is not possible the token still exists, so keep what we have and
	// let the next refresh fill in the rest.
	if tokens, _, listErr := r.Services.ServiceAccountTokens.ListServiceAccountTokens(ctx); listErr == nil {
		if found := findServiceAccountToken(tokens, plan.ID.ValueString()); found != nil {
			plan.MapFromApi(found)
		}
	} else {
		tflog.Warn(ctx, fmt.Sprintf("Created Service Account Token %q but could not read it back: %v", payload.Name, listErr))
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully created Service Account Token %q", payload.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *securityServiceAccountTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.SecurityServiceAccountTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(stateDataErrorMessage, resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	// There is no get-by-id endpoint, so read the list and pick out this token.
	tokens, httpResponse, err := r.Services.ServiceAccountTokens.ListServiceAccountTokens(ctx)
	if err != nil {
		errors.HandleAPIError(
			"Error reading Service Account Tokens",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	found := findServiceAccountToken(tokens, state.ID.ValueString())
	if found == nil {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddWarning(
			"Service Account Token does not exist",
			fmt.Sprintf("Service Account Token %q (%s) no longer exists in Nexus Repository, so it will be removed from the Terraform state and recreated by the next apply.", state.Name.ValueString(), state.ID.ValueString()),
		)
		return
	}

	state.MapFromApi(found)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is never called: every argument forces replacement. It exists to fail clearly should
// the schema ever allow an in-place change.
func (r *securityServiceAccountTokenResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Service Account Tokens cannot be updated",
		"Nexus Repository does not support updating a Service Account Token. Changing an argument should have replaced it; "+
			"please report this as a provider bug.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *securityServiceAccountTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.SecurityServiceAccountTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(stateDataErrorMessage, resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	httpResponse, err := r.Services.ServiceAccountTokens.DeleteServiceAccountToken(ctx, state.ID.ValueString())
	if err != nil {
		// Already deleted outside Terraform: nothing left to do.
		if httpResponse != nil && errors.IsNotFound(httpResponse.StatusCode) {
			tflog.Info(ctx, fmt.Sprintf("Service Account Token %q was already deleted", state.Name.ValueString()))
			return
		}
		errors.HandleAPIError(
			"Error deleting Service Account Token",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error deleting Service Account Token",
			fmt.Sprintf("Unexpected Response Code whilst deleting Service Account Token: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully deleted Service Account Token %q", state.Name.ValueString()))
}

// findServiceAccountToken returns the token with the given id, or nil.
func findServiceAccountToken(tokens []sonatyperepoV397.ServiceAccountXO, id string) *sonatyperepoV397.ServiceAccountXO {
	for i := range tokens {
		if tokens[i].Id != nil && *tokens[i].Id == id {
			return &tokens[i]
		}
	}
	return nil
}
