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
	"io"
	"net/http"
	"time"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

// securitySsrfProtectionResource is the resource implementation.
type securitySsrfProtectionResource struct {
	common.BaseResource
}

// NewSecuritySsrfProtectionResource is a helper function to simplify the provider implementation.
func NewSecuritySsrfProtectionResource() resource.Resource {
	return &securitySsrfProtectionResource{}
}

// Metadata returns the resource type name.
func (r *securitySsrfProtectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_ssrf_protection"
}

// Schema defines the schema for the resource.
func (r *securitySsrfProtectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	emptySet := types.SetValueMust(types.StringType, []attr.Value{})

	resp.Schema = tfschema.Schema{
		Description: "Manage SSRF (Server-Side Request Forgery) Protection settings. As of Sonatype Nexus Repository 3.92, requests to private (RFC1918) addresses are blocked unless explicitly allowed here.",
		Attributes: map[string]tfschema.Attribute{
			"enabled": schema.ResourceRequiredBool("Whether or not SSRF Protection is enabled"),
			"allowed_domains": schema.ResourceOptionalStringSetWithDefault(
				"Domain names allowed to bypass SSRF Protection",
				setdefault.StaticValue(emptySet),
			),
			"allowed_ips": schema.ResourceOptionalStringSetWithDefault(
				"IP addresses allowed to bypass SSRF Protection",
				setdefault.StaticValue(emptySet),
			),
			"last_updated": schema.ResourceLastUpdated(),
		},
	}
}

// fetchSsrfConfig reads the current SSRF Protection configuration from the API.
func (r *securitySsrfProtectionResource) fetchSsrfConfig(ctx context.Context) (*sonatyperepo.SsrfProtectionConfigurationXO, *http.Response, error) {
	httpResponse, err := r.Client.SecurityManagementSSRFProtectionAPI.GetConfiguration(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, httpResponse, fmt.Errorf("could not read response body: %w", err)
	}

	var cfg sonatyperepo.SsrfProtectionConfigurationXO
	if err := json.Unmarshal(body, &cfg); err != nil {
		return nil, httpResponse, fmt.Errorf("could not parse response: %w", err)
	}

	return &cfg, httpResponse, nil
}

// ImportState imports the resource into Terraform state.
func (r *securitySsrfProtectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Since this is a singleton resource (there's only one SSRF Protection configuration),
	// we don't need to parse the import ID. We just read the current configuration.
	ctx = r.AuthContext(ctx)

	cfg, httpResponse, err := r.fetchSsrfConfig(ctx)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				"Your user is unauthorized to access this resource or feature during import.",
			)
		} else {
			errors.HandleAPIError(
				"Error importing SSRF Protection settings",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	var state model.SsrfProtectionModel
	state.MapFromApi(cfg)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully imported SSRF Protection system resource")
}

// Create creates the resource and sets the initial Terraform state.
func (r *securitySsrfProtectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.SsrfProtectionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	payload := sonatyperepo.SsrfProtectionConfigurationXO{}
	plan.MapToApi(&payload)

	httpResponse, err := r.Client.SecurityManagementSSRFProtectionAPI.UpdateConfiguration(ctx).Body(payload).Execute()

	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				common.ERROR_MESSAGE_UNAUTHORIZED,
			)
		} else {
			errors.HandleAPIError(
				"Error updating SSRF Protection settings",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	if httpResponse.StatusCode >= http.StatusOK && httpResponse.StatusCode < http.StatusMultipleChoices {
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		errors.HandleAPIError(
			"Update of SSRF Protection settings was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *securitySsrfProtectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.SsrfProtectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	cfg, httpResponse, err := r.fetchSsrfConfig(ctx)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				common.ERROR_MESSAGE_UNAUTHORIZED,
			)
		} else {
			errors.HandleAPIError(
				"Error reading SSRF Protection settings",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	state.MapFromApi(cfg)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *securitySsrfProtectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.SsrfProtectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	payload := sonatyperepo.SsrfProtectionConfigurationXO{}
	plan.MapToApi(&payload)

	httpResponse, err := r.Client.SecurityManagementSSRFProtectionAPI.UpdateConfiguration(ctx).Body(payload).Execute()

	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				common.ERROR_MESSAGE_UNAUTHORIZED,
			)
		} else {
			errors.HandleAPIError(
				"Error updating SSRF Protection settings",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	if httpResponse.StatusCode >= http.StatusOK && httpResponse.StatusCode < http.StatusMultipleChoices {
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Failed to update SSRF Protection settings",
			fmt.Sprintf("Unable to update SSRF Protection settings: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// Since this is a singleton configuration resource, this resets SSRF Protection back to
// its secure, out-of-the-box default: enabled, with no allowed domains or IPs.
func (r *securitySsrfProtectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.SsrfProtectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.AuthContext(ctx)

	payload := sonatyperepo.SsrfProtectionConfigurationXO{
		Enabled:        true,
		AllowedDomains: []string{},
		AllowedIPs:     []string{},
	}

	httpResponse, err := r.Client.SecurityManagementSSRFProtectionAPI.UpdateConfiguration(ctx).Body(payload).Execute()

	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				common.ERROR_MESSAGE_UNAUTHORIZED,
			)
		} else {
			errors.HandleAPIError(
				"Error removing SSRF Protection settings",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}
}
