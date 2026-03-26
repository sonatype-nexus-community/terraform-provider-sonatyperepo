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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

// Ensure resource satisfies various resource interfaces.
var (
	_ resource.Resource                = &systemConfigHttpResource{}
	_ resource.ResourceWithImportState = &systemConfigHttpResource{}
)

// systemConfigHttpResource is the resource implementation.
type systemConfigHttpResource struct {
	common.BaseResource
}

// NewSystemConfigHttpResource is a helper function to simplify the provider implementation.
func NewSystemConfigHttpResource() resource.Resource {
	return &systemConfigHttpResource{}
}

// Metadata returns the resource type name.
func (r *systemConfigHttpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_config_http"
}

// Schema defines the schema for the resource.
func (r *systemConfigHttpResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	authenticationAttribute := schema.ResourceOptionalSingleNestedAttribute(
		"Proxy Authentication settings",
		map[string]tfschema.Attribute{
			"enabled":     schema.ResourceOptionalBoolWithDefault("Proxy Authentication enabled", false),
			"username":    schema.ResourceOptionalString("Proxy Username"),
			"password":    schema.ResourceOptionalSensitiveStringWithLengthAtLeast("Proxy Password", 1),
			"ntlm_host":   schema.ResourceOptionalStringWithDefault("Proxy NTLM Host", ""),
			"ntlm_domain": schema.ResourceOptionalStringWithDefault("Proxy NTLM Domain", ""),
		},
	)

	proxyAttributes := map[string]tfschema.Attribute{
		"enabled":        schema.ResourceRequiredBool("Whether enabled"),
		"host":           schema.ResourceOptionalString("Proxy Server Hostname"),
		"port":           schema.ResourceOptionalInt32("Proxy Server Port"),
		"authentication": authenticationAttribute,
	}

	nonProxyHosts := schema.ResourceOptionalStringSet("Hosts to exclude from HTTP/HTTPS Proxy")
	nonProxyHosts.Computed = true
	emptySet, err := types.SetValue(types.StringType, []attr.Value{})
	if err != nil {
		resp.Diagnostics.AddError("Failed to generate schema", "Could not generate empty set for non_proxy_hosts")
		return
	}
	nonProxyHosts.Default = setdefault.StaticValue(emptySet)

	resp.Schema = tfschema.Schema{
		Description: "Configure the System HTTP settings",
		Attributes: map[string]tfschema.Attribute{
			"http_proxy": schema.ResourceRequiredSingleNestedAttribute(
				"HTTP Proxy settings",
				proxyAttributes,
			),
			"https_proxy": schema.ResourceRequiredSingleNestedAttribute(
				"HTTPS Proxy settings",
				proxyAttributes,
			),
			"non_proxy_hosts": nonProxyHosts,
			"retries": schema.ResourceOptionalInt32WithDefault(
				"Maximum number of retry attempts if the initial connection attempt suffers a timeout",
				common.HTTP_SETTINGS_DEFAULT_RETRIES,
			),
			"timeout": schema.ResourceOptionalInt32WithDefault(
				"Time (seconds) to wait for activity before stopping and retrying the connection",
				common.HTTP_SETTINGS_DEFAULT_TIMEOUT,
			),
			"user_agent": schema.ResourceOptionalStringWithDefault(
				"Custom fragment to append to “User-Agent” header in HTTP requests",
				"",
			),
			"last_updated": schema.ResourceLastUpdated(),
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *systemConfigHttpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.HttpConfigurationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Apply schema defaults for fields that are null
	if plan.Retries.IsNull() {
		plan.Retries = types.Int32Value(common.HTTP_SETTINGS_DEFAULT_RETRIES)
	}
	if plan.Timeout.IsNull() {
		plan.Timeout = types.Int32Value(common.HTTP_SETTINGS_DEFAULT_TIMEOUT)
	}

	r.updateHttpSettings(ctx, &plan, &resp.Diagnostics, &resp.State)
}

// Read refreshes the Terraform state with the latest data.
func (r *systemConfigHttpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.HttpConfigurationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Call API to Create
	ctx = r.AuthContext(ctx)
	apiResponse, httpResponse, err := r.Client.ManageSonatypeHTTPSystemSettingsAPI.GetHttpSettings(ctx).Execute()

	// Handle any errors
	if err != nil {
		errors.HandleAPIError(
			"Failed to read System HTTP Settings",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	// Update State from Response
	state.MapFromApi(apiResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *systemConfigHttpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan model.HttpConfigurationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Apply schema defaults for fields that are null
	if plan.Retries.IsNull() {
		plan.Retries = types.Int32Value(common.HTTP_SETTINGS_DEFAULT_RETRIES)
	}
	if plan.Timeout.IsNull() {
		plan.Timeout = types.Int32Value(common.HTTP_SETTINGS_DEFAULT_TIMEOUT)
	}

	r.updateHttpSettings(ctx, &plan, &resp.Diagnostics, &resp.State)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *systemConfigHttpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Call API to Create
	ctx = r.AuthContext(ctx)
	httpResponse, err := r.Client.ManageSonatypeHTTPSystemSettingsAPI.ResetHttpSettings(ctx).Execute()

	// Handle Error
	if err != nil {
		errors.HandleAPIError(
			"Error removing System HTTP configuration",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		errors.HandleAPIError(
			"Failed removing System HTTP configuration",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}

	// Remove resource from State
	resp.State.RemoveResource(ctx)
}

func (r *systemConfigHttpResource) updateHttpSettings(ctx context.Context, plan *model.HttpConfigurationModel, respDiags *diag.Diagnostics, respState *tfsdk.State) {
	// Call API to Create
	ctx = r.AuthContext(ctx)
	httpSettings := v3.NewHttpSettingsXoWithDefaults()
	plan.MapToApi(httpSettings)
	httpResponse, err := r.Client.ManageSonatypeHTTPSystemSettingsAPI.UpdateHttpSettings(ctx).Body(*httpSettings).Execute()

	// Handle Errors
	if err != nil {
		errors.HandleAPIError(
			"Error setting System HTTP Settings",
			&err,
			httpResponse,
			respDiags,
		)
		return
	}
	if httpResponse.StatusCode != http.StatusNoContent {
		errors.HandleAPIError(
			"Setting System HTTP Settings was not successful",
			&err,
			httpResponse,
			respDiags,
		)
		return
	}

	// Read Data back from API
	apiResponse, httpResponse, err := r.Client.ManageSonatypeHTTPSystemSettingsAPI.GetHttpSettings(ctx).Execute()

	// Handle any errors
	if err != nil {
		errors.HandleAPIError(
			"Failed to read System HTTP Settings",
			&err,
			httpResponse,
			respDiags,
		)
		return
	}

	// Update State from Response
	plan.MapFromApi(apiResponse)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := respState.Set(ctx, plan)
	respDiags.Append(diags...)
	if respDiags.HasError() {
		return
	}
}
