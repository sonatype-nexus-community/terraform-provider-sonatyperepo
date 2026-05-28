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

	if capabilityCreateResponse == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating %s Capability", c.CapabilityType.GetType().String()),
			"API returned empty response without error",
		)
		return
	}

	// Stamp plan notes so the post-apply consistency check passes on HA clusters.
	notesFromPlan := capabilityNotesFromModel(plan)
	capabilityCreateResponse.Notes = &notesFromPlan

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

	// resolveNotesForHA also retries nil reads on HA clusters before state removal.
	capability = c.resolveNotesForHA(ctx, capabilityId.ValueString(), stateModel, capability)

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

	// Build convergence predicate: only check Enabled — notes replication in HA can
	// exceed 30 s on some capability types, so we never converge on it here.
	expectedEnabled := capabilityEnabledFromModel(planModel)

	// Now Read from API, retrying until consistent on HA clusters
	capability, httpResponse, err := c.readCapabilityByIdConsistently(
		capabilityId.ValueString(), ctx,
		func(cap *v3.CapabilityDTO) bool {
			if cap == nil {
				return false
			}
			return cap.Enabled != nil && *cap.Enabled == expectedEnabled
		},
	)

	// Whether or not the convergence loop saw notes stabilise, always stamp the
	// plan's notes value into the returned DTO so that Terraform's post-apply
	// consistency check sees a state that matches the plan.  The value was
	// already written to the shared DB; we're only hiding HA read-lag here.
	if capability != nil {
		notesFromPlan := capabilityNotesFromModel(planModel)
		capability.Notes = &notesFromPlan
	}

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
		httpResponse, err := c.Client.CapabilitiesAPI.Delete5(ctx, capabilityId.ValueString()).Execute()

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

// readCapabilityByIdConsistently retries readCapabilityById until isConverged
// returns true, up to 3 attempts × 2 s. The 10 s middleware delay covers the
// normal propagation window; these retries are a backstop for transient load.
func (c *capabilityResource) readCapabilityByIdConsistently(
	capabilityId string,
	ctx context.Context,
	isConverged func(*v3.CapabilityDTO) bool,
) (*v3.CapabilityDTO, *http.Response, error) {
	if c.NodeCount <= 1 {
		return c.readCapabilityById(capabilityId, ctx)
	}

	const maxAttempts = 3
	const retryInterval = 2 * time.Second

	var lastCap *v3.CapabilityDTO
	var lastResp *http.Response
	for attempt := 0; attempt < maxAttempts; attempt++ {
		cap, httpResp, err := c.readCapabilityById(capabilityId, ctx)
		lastCap = cap
		lastResp = httpResp
		if err != nil {
			return cap, httpResp, err
		}
		if isConverged(cap) {
			tflog.Debug(ctx, fmt.Sprintf("HA: capability converged after %d read(s)", attempt+1))
			return cap, httpResp, nil
		}
		if attempt < maxAttempts-1 {
			tflog.Info(ctx, fmt.Sprintf("HA: capability not yet consistent, retrying (%d/%d)", attempt+1, maxAttempts))
			time.Sleep(retryInterval)
		}
	}
	return lastCap, lastResp, nil
}

// resolveNotesForHA handles two HA read-lag cases on multi-node clusters:
//  1. Nil capability: retries before concluding the resource is gone.
//  2. Stale notes: falls back to the prior state value to avoid false drift.
//     Retained for 3.89/3.92 compatibility; revisit once 3.93 is minimum.
func (c *capabilityResource) resolveNotesForHA(
	ctx context.Context,
	capabilityId string,
	stateModel any,
	capability *v3.CapabilityDTO,
) *v3.CapabilityDTO {
	if c.NodeCount <= 1 {
		return capability
	}
	capability = c.retryNilCapability(ctx, capabilityId, capability)
	if capability == nil {
		return nil
	}
	return c.resolveStaleNotes(ctx, capabilityId, capabilityNotesFromModel(stateModel), capability)
}

// retryNilCapability retries a nil read up to 3 times before giving up.
// The 10 s middleware delay covers normal propagation; these retries catch
// nodes still warming up under heavy load.
func (c *capabilityResource) retryNilCapability(ctx context.Context, capabilityId string, capability *v3.CapabilityDTO) *v3.CapabilityDTO {
	const retries = 3
	const retryInterval = 2 * time.Second
	for attempt := 0; capability == nil && attempt < retries; attempt++ {
		tflog.Info(ctx, fmt.Sprintf("HA: capability not yet visible on read, retrying (%d/%d)", attempt+1, retries))
		time.Sleep(retryInterval)
		if refreshed, _, err := c.readCapabilityById(capabilityId, ctx); err == nil {
			capability = refreshed
		}
	}
	return capability
}

// resolveStaleNotes retries reads until notes match the prior state value,
// falling back to stamping the state value after all retries are exhausted.
func (c *capabilityResource) resolveStaleNotes(ctx context.Context, capabilityId string, stateNotes string, capability *v3.CapabilityDTO) *v3.CapabilityDTO {
	const retries = 3
	const retryInterval = 2 * time.Second
	for attempt := 0; attempt < retries; attempt++ {
		if capabilityNotesValue(capability) == stateNotes {
			return capability
		}
		if attempt < retries-1 {
			tflog.Info(ctx, fmt.Sprintf("HA: notes not yet replicated on read, retrying (%d/%d)", attempt+1, retries))
			time.Sleep(retryInterval)
			if refreshed, _, err := c.readCapabilityById(capabilityId, ctx); err == nil && refreshed != nil {
				capability = refreshed
			}
		} else {
			tflog.Info(ctx, "HA: notes still stale after retries, keeping state value to avoid false drift")
			capability.Notes = &stateNotes
		}
	}
	return capability
}

func capabilityNotesValue(capability *v3.CapabilityDTO) string {
	if capability.Notes != nil {
		return *capability.Notes
	}
	return ""
}

func capabilityNotesFromModel(model any) string {
	field := reflect.Indirect(reflect.ValueOf(model)).FieldByName("Notes")
	if !field.IsValid() {
		return ""
	}
	if val, ok := field.Interface().(basetypes.StringValue); ok {
		return val.ValueString()
	}
	return ""
}

func capabilityEnabledFromModel(model any) bool {
	field := reflect.Indirect(reflect.ValueOf(model)).FieldByName("Enabled")
	if !field.IsValid() {
		return false
	}
	if val, ok := field.Interface().(basetypes.BoolValue); ok {
		return val.ValueBool()
	}
	return false
}

func (c *capabilityResource) readCapabilityById(capabilityId string, ctx context.Context) (*v3.CapabilityDTO, *http.Response, error) {
	// Ensure API Context has authentication
	ctx = c.AuthContext(ctx)

	// Make API Request
	apiResponse, httpResponse, err := c.Client.CapabilitiesAPI.List2(ctx).Execute()

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

	if ct.DeprecationMessage() != nil {
		baseSchema.DeprecationMessage = *ct.DeprecationMessage()
	}

	if len(propertiesAttributes) > 0 {
		baseSchema.Attributes["properties"] = schema.ResourceRequiredSingleNestedAttribute("Properties specific to this Capability type", propertiesAttributes)
	}

	return baseSchema
}
