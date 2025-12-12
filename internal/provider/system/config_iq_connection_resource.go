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
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// systemConfigMailResource is the resource implementation.
type systemConfigIqConnectionResource struct {
	common.BaseResource
}

// NewSystemConfigIqConnectionResource is a helper function to simplify the provider implementation.
func NewSystemConfigIqConnectionResource() resource.Resource {
	return &systemConfigIqConnectionResource{}
}

// Metadata returns the resource type name.
func (r *systemConfigIqConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_iq_connection"
}

// Schema defines the schema for the resource.
func (r *systemConfigIqConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configure the Sonatype IQ Server Connection",
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: "Whether to use Sonatype Repository Firewall",
				Required:    true,
			},
			"url": schema.StringAttribute{
				Description: "The address of your Sonatype IQ Server",
				Required:    true,
				Optional:    false,
			},
			"nexus_trust_store_enabled": schema.BoolAttribute{
				Description: "Use certificates stored in the Nexus Repository Manager truststore to connect to Sonatype IQ Server",
				Required:    true,
				Optional:    false,
			},
			"authentication_method": schema.StringAttribute{
				Description: "Username to use for authentication with SMTP Server",
				Required:    true,
				Optional:    false,
				Validators: []validator.String{
					stringvalidator.OneOf(
						common.IQ_AUTHENTICATON_TYPE_USER,
						common.IQ_AUTHENTICATON_TYPE_PKI,
					),
				},
			},
			"username": schema.StringAttribute{
				Description: "User with access to Sonatype Repository Firewall",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "Credentials for the Sonatype Repository Firewall User",
				Required:    true,
				Sensitive:   true,
			},
			"connection_timeout": schema.Int32Attribute{
				Description: "Seconds to wait for activity before stopping and retrying the connection.",
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(defaultConnectionTimeoutSeconds),
				Validators: []validator.Int32{
					int32validator.Between(1, 3600),
				},
			},
			"properties": schema.StringAttribute{
				Description: "Additional properties to configure for Sonatype Repository Firewall",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"show_iq_server_link": schema.BoolAttribute{
				Description: "Show Sonatype Repository Firewall link in Browse menu when server is enabled",
				Required:    true,
				Optional:    false,
			},
			"fail_open_mode_enabled": schema.BoolAttribute{
				Description: "Allow by default when quarantine is enabled and the connection to Sonatype IQ Server fails",
				Required:    true,
				Optional:    false,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *systemConfigIqConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Call Update API
	plan := r.doUpdateRequest(ctx, &req.Plan, &resp.Diagnostics)
	if plan == nil {
		return
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
func (r *systemConfigIqConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.IqConnectionModel
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
	apiResponse, httpResponse, err := r.Client.ManageSonatypeRepositoryFirewallConfigurationAPI.GetConfiguration(ctx).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"Sonatype IQ Connection does not exist",
				fmt.Sprintf("Unable to read Sonatype IQ Connection Configuration: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading Sonatype IQ Connection Configuration",
				fmt.Sprintf("Unable to read Sonatype IQ Connection Configuration: %s: %s", httpResponse.Status, err),
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
func (r *systemConfigIqConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Call Update API
	plan := r.doUpdateRequest(ctx, &req.Plan, &resp.Diagnostics)
	if plan == nil {
		return
	}

	// Update State
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *systemConfigIqConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Disable API Call
	httpResponse, err := r.Client.ManageSonatypeRepositoryFirewallConfigurationAPI.DisableIq(ctx).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"Sonatype IQ Connection does not exist",
				fmt.Sprintf("Unable to disable Sonatype IQ Connection Configuration: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error disabling Sonatype IQ Connection Configuration",
				fmt.Sprintf("Unable to disable Sonatype IQ Connection Configuration: %s: %s", httpResponse.Status, err),
			)
		}
		return
	}

	// Remove resource from State
	resp.State.RemoveResource(ctx)
}

func (r *systemConfigIqConnectionResource) doUpdateRequest(ctx context.Context, reqPlan *tfsdk.Plan, respDiags *diag.Diagnostics) *model.IqConnectionModel {
	var plan model.IqConnectionModel
	respDiags.Append(reqPlan.Get(ctx, &plan)...)

	if respDiags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", respDiags.Errors()))
		return nil
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	apiModel := sonatyperepo.NewIqConnectionXoWithDefaults()
	plan.MapToApi(apiModel)
	httpResponse, err := r.Client.ManageSonatypeRepositoryFirewallConfigurationAPI.UpdateConfiguration(ctx).Body(*apiModel).Execute()

	// Handle Error
	if err != nil {
		respDiags.AddError(
			"Error setting Sonatype IQ Connection configuration",
			fmt.Sprintf("Error setting Sonatype IQ Connection configuration: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return nil
	} else if httpResponse.StatusCode != http.StatusNoContent {
		respDiags.AddError(
			"Error setting Sonatype IQ Connection configuration",
			fmt.Sprintf("Unexpected Response Code whilst setting Sonatype IQ Connection configuration: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
	}

	return &plan
}
