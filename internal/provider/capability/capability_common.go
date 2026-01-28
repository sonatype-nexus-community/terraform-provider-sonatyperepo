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

package capability

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"

	capabilitytype "terraform-provider-sonatyperepo/internal/provider/capability/capability_type"
)

const (
	CAPABILITY_ERROR_RESPONSE_PREFIX           = "Error response: "
	CAPABILITY_GENERAL_ERROR_RESPONSE_GENERAL  = CAPABILITY_ERROR_RESPONSE_PREFIX + " %s"
	CAPABILITY_GENERAL_ERROR_RESPONSE_WITH_ERR = CAPABILITY_ERROR_RESPONSE_PREFIX + " %s - %s"
	CAPABILITY_ERROR_DID_NOT_EXIST             = "%s (ID=%s) Capability did not exist to %s"
)

// Generic to all Task Resources
type capabilityResource struct {
	common.BaseResource
	CapabilityType capabilitytype.CapabilityTypeI
}

// Metadata returns the resource type name.
func (c *capabilityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, c.CapabilityType.ResourceName())
}

// Set Schema for this Resource
func (c *capabilityResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = capabilitySchema(c.CapabilityType)
}

// This allows import of existing capabilities into Terraform state.
func (c *capabilityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the Capability ID as the import identifier
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (c *capabilityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	plan, diags := c.CapabilityType.PlanAsModel(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting Plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Request Context
	ctx = c.AuthContext(ctx)

	// Make API requet
	capabilityCreateResponse, httpResponse, err := c.CapabilityType.DoCreateRequest(plan, c.Client, ctx, c.NxrmVersion)

	// Handle Errors
	if err != nil {
		errors.HandleAPIError(
			fmt.Sprintf("Error creating %s Capability", c.CapabilityType.GetType().String()),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}
	if !slices.Contains(c.CapabilityType.ApiCreateSuccessResponseCodes(), httpResponse.StatusCode) {
		errors.HandleAPIError(
			fmt.Sprintf("Creation of %s Capability was not successful", c.CapabilityType.GetType().String()),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}

	stateModel := c.CapabilityType.UpdateStateFromApi(plan, capabilityCreateResponse)
	stateModel = c.CapabilityType.MapFromPlanToState(plan, stateModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *capabilityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	stateModel, diags := c.CapabilityType.StateAsModel(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	// Handle any errors
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Set API Context
	ctx = c.AuthContext(ctx)

	// Make API Request
	capabilityId, shouldReturn := capabilityIdFromState(stateModel, &resp.Diagnostics)
	if shouldReturn {
		return
	}
	capability, httpResponse, err := c.readCapabilityById(capabilityId.ValueString(), ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			errors.HandleAPIWarning(
				fmt.Sprintf(CAPABILITY_ERROR_DID_NOT_EXIST, c.CapabilityType.Key(), capabilityId.ValueString(), "read"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
				fmt.Sprintf(CAPABILITY_ERROR_DID_NOT_EXIST, c.CapabilityType.Key(), capabilityId.ValueString(), "read"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	if capability == nil {
		resp.State.RemoveResource(ctx)
		errors.HandleAPIWarning(
			fmt.Sprintf(CAPABILITY_ERROR_DID_NOT_EXIST, c.CapabilityType.Key(), capabilityId.ValueString(), "read"),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	currentStateModel := c.CapabilityType.UpdateStateFromApi(stateModel, capability)
	currentStateModel = c.CapabilityType.MapFromPlanToState(stateModel, currentStateModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, &currentStateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *capabilityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	planModel, diags := c.CapabilityType.PlanAsModel(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)

	// Retrieve values from state
	stateModel, diags := c.CapabilityType.StateAsModel(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	// Request Context
	ctx = c.AuthContext(ctx)

	// Make API requet
	capabilityId, shouldReturn := capabilityIdFromState(stateModel, &resp.Diagnostics)
	if shouldReturn {
		return
	}
	httpResponse, err := c.CapabilityType.DoUpdateRequest(planModel, capabilityId.ValueString(), c.Client, ctx, c.NxrmVersion)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			errors.HandleAPIWarning(
				fmt.Sprintf(CAPABILITY_ERROR_DID_NOT_EXIST, c.CapabilityType.GetType().String(), capabilityId.ValueString(), "update"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
				fmt.Sprintf(CAPABILITY_ERROR_DID_NOT_EXIST, c.CapabilityType.GetType().String(), capabilityId.ValueString(), "update"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Now Read from API
	capability, httpResponse, err := c.readCapabilityById(capabilityId.ValueString(), ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			errors.HandleAPIWarning(
				fmt.Sprintf(CAPABILITY_ERROR_DID_NOT_EXIST, c.CapabilityType.Key(), capabilityId.ValueString(), "update"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
				fmt.Sprintf(CAPABILITY_ERROR_DID_NOT_EXIST, c.CapabilityType.Key(), capabilityId.ValueString(), "update"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	if capability == nil {
		resp.State.RemoveResource(ctx)
		errors.HandleAPIWarning(
			fmt.Sprintf(CAPABILITY_ERROR_DID_NOT_EXIST, c.CapabilityType.Key(), capabilityId.ValueString(), "update"),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	stateModel = c.CapabilityType.UpdateStateFromApi(stateModel, capability)
	stateModel = c.CapabilityType.MapFromPlanToState(planModel, stateModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *capabilityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	state, diags := c.CapabilityType.StateAsModel(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	// Handle any errors
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Request Context
	ctx = c.AuthContext(ctx)

	// Make API request
	capabilityId, shouldReturn := capabilityIdFromState(state, &resp.Diagnostics)
	if shouldReturn {
		return
	}

	attempts := 1
	maxAttempts := 3
	success := false

	for !success && attempts < maxAttempts {
		httpResponse, err := c.Client.CapabilitiesAPI.Delete4(ctx, capabilityId.ValueString()).Execute()

		// Trap 500 Error as they occur when Repo is not in appropriate internal state
		if httpResponse.StatusCode == http.StatusInternalServerError {
			tflog.Info(ctx, fmt.Sprintf("Unexpected response when deleting Capability %s (attempt %d)", c.CapabilityType.GetType().String(), attempts))
			attempts++
			continue
		}

		if err != nil {
			if httpResponse.StatusCode == http.StatusNotFound {
				resp.State.RemoveResource(ctx)
				errors.HandleAPIWarning(
					fmt.Sprintf(CAPABILITY_ERROR_DID_NOT_EXIST, c.CapabilityType.GetType().String(), capabilityId.ValueString(), "delete"),
					&err,
					httpResponse,
					&resp.Diagnostics,
				)
			} else {
				errors.HandleAPIError(
					fmt.Sprintf(CAPABILITY_ERROR_DID_NOT_EXIST, c.CapabilityType.GetType().String(), capabilityId.ValueString(), "delete"),
					&err,
					httpResponse,
					&resp.Diagnostics,
				)
			}
			return
		}
		if httpResponse.StatusCode != http.StatusNoContent {
			errors.HandleAPIError(
				fmt.Sprintf("Unexpected response when deleting %s Capability (attempt %d)", c.CapabilityType.GetType().String(), attempts),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)

			time.Sleep(1 * time.Second)
			attempts++
		} else {
			success = true
		}
	}
}

func (c *capabilityResource) readCapabilityById(capabilityId string, ctx context.Context) (*v3.CapabilityDTO, *http.Response, error) {
	// Ensure API Context has authentication
	ctx = c.AuthContext(ctx)

	// Make API Request
	apiResponse, httpResponse, err := c.Client.CapabilitiesAPI.List(ctx).Execute()

	// Handle any errors
	if err != nil {
		return nil, httpResponse, err
	}

	// Find the actual Capability from the list returned
	var capability *v3.CapabilityDTO
	for _, cap := range apiResponse {
		if cap.Id != nil && *cap.Id == capabilityId {
			capability = &cap
			break
		}
	}

	return capability, httpResponse, nil
}

func capabilityIdFromState(state any, respDiags *diag.Diagnostics) (basetypes.StringValue, bool) {
	capabilityIdStructField := reflect.Indirect(reflect.ValueOf(state)).FieldByName("Id").Interface()
	capabilityId, ok := capabilityIdStructField.(basetypes.StringValue)
	if !ok {
		respDiags.AddError(
			"Failed to determine Capability ID to delete from State",
			fmt.Sprintf("%s %s", CAPABILITY_ERROR_RESPONSE_PREFIX, capabilityIdStructField),
		)
		return basetypes.StringValue{}, true
	}
	return capabilityId, false
}

func capabilitySchema(ct capabilitytype.CapabilityTypeI) tfschema.Schema {
	propertiesAttributes := ct.PropertiesSchema()

	baseSchema := tfschema.Schema{
		MarkdownDescription: ct.GetMarkdownDescription() + `
		
**NOTE:** Requires Sonatype Nexus Repostiory 3.84.0 or later.`,
		Attributes: map[string]tfschema.Attribute{
			"id":           schema.ResourceComputedString("The internal ID of the Capability."),
			"notes":        schema.ResourceOptionalStringWithDefault("Optional notes about configured capability.", ""),
			"enabled":      schema.ResourceRequiredBool("Whether the Capability is enabled."),
			"last_updated": schema.ResourceLastUpdated(),
		},
	}

	if len(propertiesAttributes) > 0 {
		baseSchema.Attributes["properties"] = schema.ResourceRequiredSingleNestedAttribute("Properties specific to this Capability type", propertiesAttributes)
	}

	return baseSchema
}
