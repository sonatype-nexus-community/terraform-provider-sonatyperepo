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

package privilege

import (
	"context"
	"fmt"
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	tfschema "github.com/sonatype-nexus-community/terraform-provider-shared/schema"
	"maps"
	"net/http"
	"reflect"
	"regexp"
	"slices"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/privilege/privilege_type"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

const privilegeNamePattern = `^[a-zA-Z0-9\-]{1}[a-zA-Z0-9_\-\.]*$`

const (
	PRIVILEGE_ERROR_RESPONSE_PREFIX           = "Error response: "
	PRIVILEGE_GENERAL_ERROR_RESPONSE_GENERAL  = PRIVILEGE_ERROR_RESPONSE_PREFIX + " %s"
	PRIVILEGE_GENERAL_ERROR_RESPONSE_WITH_ERR = PRIVILEGE_ERROR_RESPONSE_PREFIX + " %s - %s"
	PRIVILEGE_ERROR_DID_NOT_EXIST             = "%s Privilege did not exist to %s"
)

// Generic Resource for all Privilege Types
type privilegeResource struct {
	common.BaseResource
	PrivilegeType     privilege_type.PrivilegeType
	PrivilegeTypeType privilege_type.PrivilegeTypeType
}

// Metadata returns the resource type name.
func (r *privilegeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, r.PrivilegeType.GetResourceName(r.PrivilegeTypeType))
}

// Schema defines the schema for the resource.
func (r *privilegeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := getBasePrivilegeSchema(r.PrivilegeTypeType)
	maps.Copy(schema.Attributes, r.PrivilegeType.GetPrivilegeTypeSchemaAttributes())
	if r.PrivilegeType.IsDeprecated() {
		schema.DeprecationMessage = "Groovy scripting has been disbaled by default since Sonatype Nexus Repository 3.21.2 - see https://help.sonatype.com/en/script-api.html"
	}
	resp.Schema = schema
}

// Create creates the resource and sets the initial Terraform state.
func (r *privilegeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	plan, diags := r.PrivilegeType.GetPlanAsModel(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting Plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Request Context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Make API requet
	httpResponse, err := r.PrivilegeType.DoCreateRequest(plan, r.Client, ctx)

	// Handle Errors
	if err != nil {
		sharederr.HandleAPIError(
			fmt.Sprintf("Error creating %s Privilege", r.PrivilegeTypeType.String()),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}
	if !slices.Contains(r.PrivilegeType.GetApiCreateSuccessResponseCodes(), httpResponse.StatusCode) {
		sharederr.HandleAPIError(
			fmt.Sprintf("Creation of %s Privilege was not successful", r.PrivilegeTypeType.String()),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}

	stateModel := r.PrivilegeType.UpdatePlanForState(plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *privilegeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	stateModel, diags := r.PrivilegeType.GetStateAsModel(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	// Handle any errors
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Set API Context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Make API Request
	apiResponse, httpResponse, err := r.PrivilegeType.DoReadRequest(stateModel, r.Client, ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				fmt.Sprintf(PRIVILEGE_ERROR_DID_NOT_EXIST, r.PrivilegeTypeType.String(), "read"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				fmt.Sprintf(PRIVILEGE_ERROR_DID_NOT_EXIST, r.PrivilegeTypeType.String(), "read"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Update State from Response
	r.PrivilegeType.UpdateStateFromApi(stateModel, apiResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *privilegeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	planModel, diags := r.PrivilegeType.GetPlanAsModel(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)

	// Retrieve values from state
	stateModel, diags := r.PrivilegeType.GetStateAsModel(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	// Request Context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Make API requet
	httpResponse, err := r.PrivilegeType.DoUpdateRequest(planModel, stateModel, r.Client, ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				fmt.Sprintf(PRIVILEGE_ERROR_DID_NOT_EXIST, r.PrivilegeTypeType.String(), "update"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				fmt.Sprintf(PRIVILEGE_ERROR_DID_NOT_EXIST, r.PrivilegeTypeType.String(), "update"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	stateModel = r.PrivilegeType.UpdatePlanForState(planModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *privilegeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	state, diags := r.PrivilegeType.GetStateAsModel(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	// Handle any errors
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Request Context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Make API request
	privilegeNameStructField := reflect.Indirect(reflect.ValueOf(state)).FieldByName("Name").Interface()
	privilegeName, ok := privilegeNameStructField.(basetypes.StringValue)
	if !ok {
		resp.Diagnostics.AddError(
			"Failed to determine Privilege Name to delete from State",
			fmt.Sprintf("%s %s", PRIVILEGE_ERROR_RESPONSE_PREFIX, privilegeNameStructField),
		)
		return
	}
	httpResponse, err := r.PrivilegeType.DoDeleteRequest(privilegeName.ValueString(), r.Client, ctx)

	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				fmt.Sprintf(PRIVILEGE_ERROR_DID_NOT_EXIST, r.PrivilegeTypeType.String(), "delete"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				fmt.Sprintf(PRIVILEGE_ERROR_DID_NOT_EXIST, r.PrivilegeTypeType.String(), "delete"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}
	if httpResponse.StatusCode != http.StatusNoContent {
		sharederr.HandleAPIError(
			fmt.Sprintf("Unexpected response when deleting %s Privilege", r.PrivilegeTypeType.String()),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}
}

func getBasePrivilegeSchema(privilegeTypeType privilege_type.PrivilegeTypeType) schema.Schema {
	return schema.Schema{
		Description: fmt.Sprintf("Manage a Privilege of type %s", privilegeTypeType.String()),
		Attributes: map[string]schema.Attribute{
			"name": tfschema.RequiredStringWithRegex(
				"The name of the privilege. This value cannot be changed.",
				regexp.MustCompile(privilegeNamePattern),
				`Please provide a name that complies with the Regular Expression: '^[a-zA-Z0-9\-]{1}[a-zA-Z0-9_\-\.]*$'`,
			),
			"description": tfschema.ResourceRequiredString("Friendly description of this Privilege"),
			"read_only": tfschema.ResourceComputedBoolWithDefault("Indicates whether the privilege can be changed. External values supplied to this will be ignored by the system.", false),
			"type": tfschema.ResourceRequiredStringEnum(
				"The type of privilege, each type covers different portions of the system. External values supplied to this will be ignored by the system.",
				privilegeTypeType.String(),
			),
			"last_updated": tfschema.ResourceComputedString(""),
		},
	}
}
