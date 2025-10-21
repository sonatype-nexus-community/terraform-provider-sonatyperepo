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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Ensure resource satisfies various resource interfaces.
var (
	_ resource.Resource                = &systemConfigMailResource{}
	_ resource.ResourceWithImportState = &systemConfigMailResource{}
)

// systemConfigMailResource is the resource implementation.
type systemConfigMailResource struct {
	common.BaseResource
}

// NewSystemConfigMailResource is a helper function to simplify the provider implementation.
func NewSystemConfigMailResource() resource.Resource {
	return &systemConfigMailResource{}
}

// Metadata returns the resource type name.
func (r *systemConfigMailResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_config_mail"
}

// Schema defines the schema for the resource.
func (r *systemConfigMailResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configure the System Email Server",
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: "Whether Email Server is enabled",
				Required:    true,
			},
			"host": schema.StringAttribute{
				Description: "SMTP Server Hostname",
				Required:    true,
				Optional:    false,
			},
			"port": schema.Int64Attribute{
				Description: "SMTP Server Port",
				Required:    true,
				Optional:    false,
			},
			"username": schema.StringAttribute{
				Description: "Username to use for authentication with SMTP Server",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password to use for authentication with SMTP Server",
				Optional:    true,
				Sensitive:   true,
			},
			"from_address": schema.StringAttribute{
				Description: "From Address to use when sending emails",
				Required:    true,
				Optional:    false,
			},
			"subject_prefix": schema.StringAttribute{
				Description: "A prefix to use in Subject Lines for emails that are sent",
				Optional:    true,
			},
			"start_tls_enabled": schema.BoolAttribute{
				Description: "Enable STARTTLS support for insecure connections",
				Required:    true,
				Optional:    false,
			},
			"start_tls_required": schema.BoolAttribute{
				Description: "Require STARTTLS support",
				Required:    true,
				Optional:    false,
			},
			"ssl_on_connect_enabled": schema.BoolAttribute{
				Description: "Enable SSL/TLS encryption upon connection",
				Required:    true,
				Optional:    false,
			},
			"ssl_server_identity_check_enabled": schema.BoolAttribute{
				Description: "Enable server identity check",
				Required:    true,
				Optional:    false,
			},
			"nexus_trust_store_enabled": schema.BoolAttribute{
				Description: "Use certificate connected to the Nexus Repository Truststore",
				Required:    true,
				Optional:    false,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// ImportState imports the resource into Terraform state.
func (r *systemConfigMailResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Since this is a singleton resource (system email configuration), 
	// we don't need to validate the ID - any non-empty string is acceptable
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID cannot be empty. Use any non-empty string (e.g., 'system-email-config') to import the system email configuration.",
		)
		return
	}

	// Set the ID to a fixed value since this is a singleton resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("last_updated"), types.StringValue("system-email-config"))...)

	tflog.Info(ctx, fmt.Sprintf("Imported system email configuration with ID: %s", req.ID))
}

// Create creates the resource and sets the initial Terraform state.
func (r *systemConfigMailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.EmailConfigurationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

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

	requestPayload := sonatyperepo.ApiEmailConfiguration{
		Enabled:                       plan.Enabled.ValueBoolPointer(),
		Host:                          plan.Host.ValueStringPointer(),
		Port:                          int32(*plan.Port.ValueInt64Pointer()),
		Username:                      plan.Username.ValueStringPointer(),
		Password:                      plan.Password.ValueStringPointer(),
		FromAddress:                   plan.FromAddress.ValueStringPointer(),
		SubjectPrefix:                 plan.SubjectPrefix.ValueStringPointer(),
		StartTlsEnabled:               plan.StartTLSEnabled.ValueBoolPointer(),
		StartTlsRequired:              plan.StartTLSRequired.ValueBoolPointer(),
		SslOnConnectEnabled:           plan.SSLOnConnectEnabled.ValueBoolPointer(),
		SslServerIdentityCheckEnabled: plan.SSLServerIdentityCheckEnabled.ValueBoolPointer(),
		NexusTrustStoreEnabled:        plan.NexusTrustStoreEnabled.ValueBoolPointer(),
	}
	apiResponse, err := r.Client.EmailAPI.SetEmailConfiguration(ctx).Body(requestPayload).Execute()

	// Handle Error
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting Mail Server configuration",
			fmt.Sprintf("Error setting Mail Server configuration: %d: %s", apiResponse.StatusCode, apiResponse.Status),
		)
		return
	} else if apiResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error setting Mail Server configuration",
			fmt.Sprintf("Unexpected Response Code whilst setting Mail Server configuration: %d: %s", apiResponse.StatusCode, apiResponse.Status),
		)
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *systemConfigMailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.EmailConfigurationModel

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
	apiResponse, httpResponse, err := r.Client.EmailAPI.GetEmailConfiguration(ctx).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"System Email Configuration does not exist",
				fmt.Sprintf("Unable to read System Email Configuration: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading System Email Configuration",
				fmt.Sprintf("Unable to read System Email Configuration: %s: %s", httpResponse.Status, err),
			)
		}
		return
	} else {
		// Update State
		state.Enabled = types.BoolPointerValue(apiResponse.Enabled)
		if apiResponse.Host != nil {
			state.Host = types.StringPointerValue(apiResponse.Host)
		} else {
			state.Host = types.StringNull()
		}

		port := apiResponse.Port
		state.Port = types.Int64Value(int64(port))

		if apiResponse.Username != nil {
			state.Username = types.StringPointerValue(apiResponse.Username)
		}

		// Ignore Password during reads for security reasons

		if apiResponse.FromAddress != nil {
			state.FromAddress = types.StringPointerValue(apiResponse.FromAddress)
		} else {
			state.FromAddress = types.StringNull()
		}
		
		if apiResponse.SubjectPrefix != nil {
			state.SubjectPrefix = types.StringPointerValue(apiResponse.SubjectPrefix)
		} else {
			state.SubjectPrefix = types.StringNull()
		}
		
		if apiResponse.StartTlsEnabled != nil {
			state.StartTLSEnabled = types.BoolPointerValue(apiResponse.StartTlsEnabled)
		} else {
			state.StartTLSEnabled = types.BoolNull()
		}
		
		if apiResponse.StartTlsRequired != nil {
			state.StartTLSRequired = types.BoolPointerValue(apiResponse.StartTlsRequired)
		} else {
			state.StartTLSRequired = types.BoolNull()
		}
		
		if apiResponse.SslOnConnectEnabled != nil {
			state.SSLOnConnectEnabled = types.BoolPointerValue(apiResponse.SslOnConnectEnabled)
		} else {
			state.SSLOnConnectEnabled = types.BoolNull()
		}
		
		if apiResponse.SslServerIdentityCheckEnabled != nil {
			state.SSLServerIdentityCheckEnabled = types.BoolPointerValue(apiResponse.SslServerIdentityCheckEnabled)
		} else {
			state.SSLServerIdentityCheckEnabled = types.BoolNull()
		}
		
		if apiResponse.NexusTrustStoreEnabled != nil {
			state.NexusTrustStoreEnabled = types.BoolPointerValue(apiResponse.NexusTrustStoreEnabled)
		} else {
			state.NexusTrustStoreEnabled = types.BoolNull()
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *systemConfigMailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.EmailConfigurationModel
	var state model.EmailConfigurationModel

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

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Call API to Update
	requestPayload := sonatyperepo.ApiEmailConfiguration{
		Enabled:                       plan.Enabled.ValueBoolPointer(),
		Host:                          plan.Host.ValueStringPointer(),
		Port:                          int32(*plan.Port.ValueInt64Pointer()),
		Username:                      plan.Username.ValueStringPointer(),
		Password:                      plan.Password.ValueStringPointer(),
		FromAddress:                   plan.FromAddress.ValueStringPointer(),
		SubjectPrefix:                 plan.SubjectPrefix.ValueStringPointer(),
		StartTlsEnabled:               plan.StartTLSEnabled.ValueBoolPointer(),
		StartTlsRequired:              plan.StartTLSRequired.ValueBoolPointer(),
		SslOnConnectEnabled:           plan.SSLOnConnectEnabled.ValueBoolPointer(),
		SslServerIdentityCheckEnabled: plan.SSLServerIdentityCheckEnabled.ValueBoolPointer(),
		NexusTrustStoreEnabled:        plan.NexusTrustStoreEnabled.ValueBoolPointer(),
	}
	apiResponse, err := r.Client.EmailAPI.SetEmailConfiguration(ctx).Body(requestPayload).Execute()

	// Handle Error
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting Mail Server configuration",
			fmt.Sprintf("Error setting Mail Server configuration: %d: %s", apiResponse.StatusCode, apiResponse.Status),
		)
		return
	} else if apiResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error setting Mail Server configuration",
			fmt.Sprintf("Unexpected Response Code whilst setting Mail Server configuration: %d: %s", apiResponse.StatusCode, apiResponse.Status),
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
func (r *systemConfigMailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.EmailConfigurationModel

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

	apiResponse, err := r.Client.EmailAPI.DeleteEmailConfiguration(ctx).Execute()

	// Handle Error
	if err != nil {
		resp.Diagnostics.AddError(
			"Error removing SMTP Mail Server configuration",
			fmt.Sprintf("Error removing SMTP Mail Server configuration: %d: %s", apiResponse.StatusCode, apiResponse.Status),
		)
		return
	} else if apiResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error removing SMTP Mail Server configuration",
			fmt.Sprintf("Unexpected Response Code whilst removing SMTP Mail Server configuration: %d: %s", apiResponse.StatusCode, apiResponse.Status),
		)
	}
}