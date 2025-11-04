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
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

const (
	stateDataErrorMessage = "Getting state data has errors: %v"
)

// securitySamlResource is the resource implementation.
type securitySamlResource struct {
	common.BaseResource
}

// SecuritySamlConfiguration defines the structure for the API request body
type SecuritySamlConfiguration struct {
	Active []string `json:"active"`
}

// NewSecuritySamlResource is a helper function to simplify the provider implementation.
func NewSecuritySamlResource() resource.Resource {
	return &securitySamlResource{}
}

// Metadata returns the resource type name.
func (r *securitySamlResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_saml"
}

// Schema defines the schema for the resource.
func (r *securitySamlResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configure Sonatype Nexus Repository Security SAML.",
		Attributes: map[string]schema.Attribute{
			"idp_metadata": schema.StringAttribute{
				Description: "SAML Identity Provider Metadata XML",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(10),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`(?s)^\s*<.*>\s*$`),
						"must be valid XML format",
					),
				},
			},
			"username_attribute": schema.StringAttribute{
				Description: "IdP field mappings for username",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"first_name_attribute": schema.StringAttribute{
				Description: "IdP field mappings for user's given name",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"last_name_attribute": schema.StringAttribute{
				Description: "IdP field mappings for user's family name",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"email_attribute": schema.StringAttribute{
				Description: "IdP field mappings for user's email",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"groups_attribute": schema.StringAttribute{
				Description: "IdP field mappings for user's groups",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"validate_response_signature": schema.BoolAttribute{
				Description: "Validate SAML response signature",
				Optional: true,
			},
			"validate_assertion_signature": schema.BoolAttribute{
				Description: "By default, if a signing key is found in the IdP metadata, then Sonatype Nexus Repository Manager will attempt to validate signatures on the assertions.",
				Optional: true,
			},
			"entity_id": schema.StringAttribute{
				Description: "SAML Entity ID (typically a URI)",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

// ImportState imports the resource state.
func (r *securitySamlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing SAML configuration", map[string]interface{}{
		"import_id": req.ID,
	})

	// Set up authentication context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Get the current SAML configuration from the API
	httpResponse, err := r.Client.SecurityManagementSAMLAPI.GetSamlConfiguration(ctx).Execute()
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			resp.Diagnostics.AddError(
				"SAML Configuration not found",
				"No SAML configuration exists to import",
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading SAML Configuration during import",
				fmt.Sprintf("Unable to read SAML Configuration: %s", err),
			)
		}
		return
	}

	// Parse the response body to get the SamlConfigurationXO
	var samlConfig sonatyperepo.SamlConfigurationXO
	if err := json.NewDecoder(httpResponse.Body).Decode(&samlConfig); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing SAML Configuration response",
			fmt.Sprintf("Unable to parse SAML Configuration response: %s", err),
		)
		return
	}

	var state model.SecuritySamlModel
	state.MapFromApi(&samlConfig)

	// Set the populated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	tflog.Info(ctx, "Successfully imported SAML configuration")
}

// Create creates the resource and sets the initial Terraform state.
func (r *securitySamlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.SecuritySamlModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Set up authentication context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	requestPayload := sonatyperepo.SamlConfigurationXO{
		IdpMetadata: plan.IdpMetadata.ValueString(),
		UsernameAttribute: plan.UsernameAttribute.ValueString(),
		FirstNameAttribute: plan.FirstNameAttribute.ValueStringPointer(),
		LastNameAttribute: plan.LastNameAttribute.ValueStringPointer(),
		EmailAttribute: plan.EmailAttribute.ValueStringPointer(),
		GroupsAttribute: plan.GroupsAttribute.ValueStringPointer(),
		ValidateResponseSignature: plan.ValidateResponseSignature.ValueBoolPointer(),
		ValidateAssertionSignature: plan.ValidateAssertionSignature.ValueBoolPointer(),
		EntityId: plan.EntityId.ValueStringPointer(),
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating security saml configuration with : %v", requestPayload))

	apiResponse, err := r.Client.SecurityManagementSAMLAPI.PutSamlConfiguration(ctx).Body(requestPayload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Security SAML configuration",
			fmt.Sprintf("Error creating Security SAML configuration: %s", err.Error()),
		)
		return
	} else if apiResponse.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Error creating Security SAML configuration",
			fmt.Sprintf("Unexpected Response Code whilst creating Security SAML configuration: %d: %s", apiResponse.StatusCode, apiResponse.Status),
		)
		return
	}

	tflog.Info(ctx, "Successfully created security saml configuration")
	
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *securitySamlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.SecuritySamlModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(stateDataErrorMessage, resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	httpResponse, err := r.Client.SecurityManagementSAMLAPI.GetSamlConfiguration(ctx).Execute()
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"SAML Configuration does not exist",
				fmt.Sprintf("Unable to read SAML Configuration: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		} else {
			errorMsg := "Unable to read SAML Configuration"
			if httpResponse != nil {
				errorMsg = fmt.Sprintf("%s: %s", errorMsg, httpResponse.Status)
			}
			resp.Diagnostics.AddError(
				"Error Reading Security SAML Configuration",
				fmt.Sprintf("%s: %s", errorMsg, err),
			)
		}
		return
	}

	// Parse the response body to get the SamlConfigurationXO
	var samlConfig sonatyperepo.SamlConfigurationXO
	if err := json.NewDecoder(httpResponse.Body).Decode(&samlConfig); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing SAML Configuration response",
			fmt.Sprintf("Unable to parse SAML Configuration response: %s", err),
		)
		return
	}

	// Update state with values from API using MapFromApi
	state.MapFromApi(&samlConfig)

	tflog.Debug(ctx, "Successfully read security SAML configuration from API")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *securitySamlResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.SecuritySamlModel
	var state model.SecuritySamlModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(stateDataErrorMessage, resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	requestPayload := sonatyperepo.SamlConfigurationXO{
		IdpMetadata: plan.IdpMetadata.ValueString(),
		UsernameAttribute: plan.UsernameAttribute.ValueString(),
		FirstNameAttribute: plan.FirstNameAttribute.ValueStringPointer(),
		LastNameAttribute: plan.LastNameAttribute.ValueStringPointer(),
		EmailAttribute: plan.EmailAttribute.ValueStringPointer(),
		GroupsAttribute: plan.GroupsAttribute.ValueStringPointer(),
		ValidateResponseSignature: plan.ValidateResponseSignature.ValueBoolPointer(),
		ValidateAssertionSignature: plan.ValidateAssertionSignature.ValueBoolPointer(),
		EntityId: plan.EntityId.ValueStringPointer(),
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating security SAML configuration with : %v", requestPayload))

	// Call API to Update
	apiResponse, err := r.Client.SecurityManagementSAMLAPI.PutSamlConfiguration(ctx).Body(requestPayload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Security SAML configuration",
			fmt.Sprintf("Error updating Security SAML configuration: %s", err.Error()),
		)
		return
	} else if apiResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error updating Security SAML configuration",
			fmt.Sprintf("Unexpected Response Code whilst updating Security SAML configuration: %d: %s", apiResponse.StatusCode, apiResponse.Status),
		)
		return
	}

	tflog.Info(ctx, "Successfully updated security SAML configuration")
	
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *securitySamlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.SecuritySamlModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf(stateDataErrorMessage, resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	apiResponse, err := r.Client.SecurityManagementSAMLAPI.DeleteSamlConfiguration(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Security SAML configuration",
			fmt.Sprintf("Error deleting Security SAML configuration: %s", err.Error()),
		)
		return
	} else if apiResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error deleting Security SAML configuration",
			fmt.Sprintf("Unexpected Response Code whilst deleting Security SAML configuration: %d: %s", apiResponse.StatusCode, apiResponse.Status),
		)
		return
	}

	tflog.Info(ctx, "Successfully deleted security SAML configuration")
}