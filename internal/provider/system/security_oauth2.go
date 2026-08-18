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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

// securityOAuth2Resource is the resource implementation.
type securityOAuth2Resource struct {
	common.BaseResource
}

// NewSecurityOAuth2Resource is a helper function to simplify the provider implementation.
func NewSecurityOAuth2Resource() resource.Resource {
	return &securityOAuth2Resource{}
}

// Metadata returns the resource type name.
func (r *securityOAuth2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_oauth2"
}

// Schema defines the schema for the resource.
func (r *securityOAuth2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	// The API always returns these as an empty (not omitted/null) JSON object when unset, so
	// they need a matching empty-map default - otherwise every plan after Create/Read drifts
	// null -> {} the same way SSRF Protection's sets needed a static empty-set default.
	emptyMap := types.MapValueMust(types.StringType, map[string]attr.Value{})
	authorizationCustomParams := schema.ResourceOptionalStringMap("Additional custom parameters sent with the authorization request")
	authorizationCustomParams.Computed = true
	authorizationCustomParams.Default = mapdefault.StaticValue(emptyMap)
	tokenRequestCustomParams := schema.ResourceOptionalStringMap("Additional custom parameters sent with the token request")
	tokenRequestCustomParams.Computed = true
	tokenRequestCustomParams.Default = mapdefault.StaticValue(emptyMap)
	exactMatchClaims := schema.ResourceOptionalStringMap("Claims that must exactly match the given value for authentication to succeed")
	exactMatchClaims.Computed = true
	exactMatchClaims.Default = mapdefault.StaticValue(emptyMap)

	resp.Schema = tfschema.Schema{
		MarkdownDescription: "Configure Sonatype Nexus Repository OAuth2 / OpenID Connect (OIDC) authentication.\n\n" +
			"**Requires Nexus Repository Pro 3.94.0 or later**, with `nexus.security.oauth2.enabled=true` and " +
			"`nexus.jwt.enabled=true` set in `nexus.properties`. The underlying API does not exist on earlier " +
			"versions - applying this resource against an older server fails with a clear error rather than a panic.",
		Attributes: map[string]tfschema.Attribute{
			"idp_jws_algorithm": schema.ResourceRequiredStringWithLengthAtLeast(
				"JWT signature algorithm used to validate tokens issued by the OpenID Provider - e.g. RS256",
				1,
			),
			"username_claim": schema.ResourceRequiredStringWithLengthAtLeast(
				"Claim used as the unique identifier for the user",
				1,
			),
			"client_id": schema.ResourceOptionalStringWithDefault(
				"Client ID issued by the OpenID Provider for this Nexus Repository instance", "",
			),
			"client_secret": schema.ResourceSensitiveString(
				"Client Secret issued by the OpenID Provider. The API never returns this value on read, so Terraform cannot detect out-of-band changes to it.",
			),
			// The API always returns these as an empty string (not omitted/null) when unset,
			// so - like the *_custom_params maps and use_trust_store above - they need a
			// matching "" default or every refresh drifts null -> "".
			"idp_authorization_url": schema.ResourceOptionalStringWithDefault("Authorization endpoint of the OpenID Provider", ""),
			"idp_token_url":         schema.ResourceOptionalStringWithDefault("Token endpoint of the OpenID Provider", ""),
			"idp_logout_url":        schema.ResourceOptionalStringWithDefault("Logout endpoint of the OpenID Provider", ""),
			"idp_jwks_url":          schema.ResourceOptionalStringWithDefault("JSON Web Key Set (JWKS) endpoint of the OpenID Provider", ""),
			"idp_jwks":              schema.ResourceOptionalStringWithDefault("Literal JSON Web Key Set (JWKS) document, as an alternative to idp_jwks_url", ""),
			"first_name_claim":      schema.ResourceOptionalStringWithDefault("Claim containing the user's first name", ""),
			"last_name_claim":       schema.ResourceOptionalStringWithDefault("Claim containing the user's last name", ""),
			"email_claim":           schema.ResourceOptionalStringWithDefault("Claim containing the user's email address", ""),
			"groups_claim":          schema.ResourceOptionalStringWithDefault("Claim containing the user's group memberships, required for role mapping", ""),
			"use_trust_store": schema.ResourceOptionalBoolWithDefault(
				"Whether to validate the OpenID Provider's certificate using Nexus Repository's own truststore",
				false,
			),
			"authorization_custom_params": authorizationCustomParams,
			"token_request_custom_params": tokenRequestCustomParams,
			"exact_match_claims":          exactMatchClaims,
			"last_updated":                schema.ResourceLastUpdated(),
		},
	}
}

// ImportState imports the resource state.
func (r *securityOAuth2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing OAuth2/OIDC configuration", map[string]interface{}{
		"import_id": req.ID,
	})

	ctx = r.AuthContext(ctx)

	apiResponse, httpResponse, err := r.Services.OAuth2.GetOAuth2Configuration(ctx)
	if err != nil {
		if httpResponse != nil && errors.IsNotFound(httpResponse.StatusCode) {
			resp.Diagnostics.AddError(
				"OAuth2/OIDC Configuration not found",
				"No OAuth2/OIDC configuration exists to import",
			)
		} else {
			errors.HandleAPIError(
				"Error reading OAuth2/OIDC configuration during import",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	var state model.SecurityOAuth2Model
	state.MapFromApi(apiResponse)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully imported OAuth2/OIDC configuration")
}

// Create creates the resource and sets the initial Terraform state.
func (r *securityOAuth2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.SecurityOAuth2Model
	// Read from Plan, not Config: several attributes (use_trust_store, the *_custom_params
	// maps) are Optional+Computed with a static default, which only resolves into the Plan -
	// Config always reflects the raw, possibly-null HCL. Writing state back from a
	// Config-sourced value here would conflict with the defaulted Plan value Terraform is
	// expecting and fail with "Provider produced inconsistent result after apply".
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	payload := sonatyperepoV395.OAuth2OidcConfigurationXO{}
	plan.MapToApi(&payload)

	tflog.Debug(ctx, fmt.Sprintf("Creating OAuth2/OIDC configuration with : %v", payload))

	// The PUT endpoint always returns 204, even on first creation - there is no separate
	// 201 response documented or observed for this API, unlike security_saml.go's PUT.
	httpResponse, err := r.Services.OAuth2.PutOAuth2Configuration(ctx, payload)
	if err != nil {
		errors.HandleAPIError(
			"Error creating OAuth2/OIDC configuration",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error creating OAuth2/OIDC configuration",
			fmt.Sprintf("Unexpected Response Code whilst creating OAuth2/OIDC configuration: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}

	tflog.Info(ctx, "Successfully created OAuth2/OIDC configuration")

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *securityOAuth2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.SecurityOAuth2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(stateDataErrorMessage, resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	apiResponse, httpResponse, err := r.Services.OAuth2.GetOAuth2Configuration(ctx)
	if err != nil {
		if httpResponse != nil && errors.IsNotFound(httpResponse.StatusCode) {
			resp.State.RemoveResource(ctx)
			errors.HandleAPIWarning(
				"OAuth2/OIDC configuration does not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
				"Error reading OAuth2/OIDC configuration",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Unlike GetSamlConfiguration, this API's GET does not document a 404 response for the
	// unconfigured case. Fall back to treating an empty pair of required fields as "removed
	// out of band" - PUT can never persist an empty usernameClaim/idpJwsAlgorithm, so a 200
	// response with both empty can only mean there is no configuration to read.
	if apiResponse.UsernameClaim == "" && apiResponse.IdpJwsAlgorithm == "" {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddWarning(
			"OAuth2/OIDC configuration does not exist",
			"The OAuth2/OIDC configuration appears to have been removed outside of Terraform.",
		)
		return
	}

	state.MapFromApi(apiResponse)

	tflog.Debug(ctx, "Successfully read OAuth2/OIDC configuration from API")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *securityOAuth2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.SecurityOAuth2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	payload := sonatyperepoV395.OAuth2OidcConfigurationXO{}
	plan.MapToApi(&payload)

	tflog.Debug(ctx, fmt.Sprintf("Updating OAuth2/OIDC configuration with : %v", payload))

	httpResponse, err := r.Services.OAuth2.PutOAuth2Configuration(ctx, payload)
	if err != nil {
		errors.HandleAPIError(
			"Error updating OAuth2/OIDC configuration",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error updating OAuth2/OIDC configuration",
			fmt.Sprintf("Unexpected Response Code whilst updating OAuth2/OIDC configuration: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}

	tflog.Info(ctx, "Successfully updated OAuth2/OIDC configuration")

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *securityOAuth2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.SecurityOAuth2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(stateDataErrorMessage, resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	httpResponse, err := r.Services.OAuth2.DeleteOAuth2Configuration(ctx)
	if err != nil {
		errors.HandleAPIError(
			"Error deleting OAuth2/OIDC configuration",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error deleting OAuth2/OIDC configuration",
			fmt.Sprintf("Unexpected Response Code whilst deleting OAuth2/OIDC configuration: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}

	tflog.Info(ctx, "Successfully deleted OAuth2/OIDC configuration")
}
