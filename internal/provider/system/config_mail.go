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

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
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
	resp.TypeName = req.ProviderTypeName + "_repository_maven_hosted"
}

// Schema defines the schema for the resource.
func (r *systemConfigMailResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configure the System EmaiL Server",
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
				Default:     int64default.StaticInt64(25),
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
				Optional:    true,
			},
			"start_tls_required": schema.BoolAttribute{
				Description: "Require STARTTLS support",
				Optional:    true,
			},
			"ssl_on_connect_enabled": schema.BoolAttribute{
				Description: "Enable SSL/TLS encryption upon connection",
				Optional:    true,
			},
			"ssl_server_identity_check_enabled": schema.BoolAttribute{
				Description: "Enable server identity check",
				Optional:    true,
			},
			"nexus_trust_store_enabled": schema.BoolAttribute{
				Description: "Use certificate connected to the Nexus Repository Truststore",
				Optional:    true,
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
	api_request := r.Client.EmailAPI.SetEmailConfiguration(ctx).Body(requestPayload)
	api_response, err := api_request.Execute()

	// Handle Error
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting Mail Server configuration",
			fmt.Sprintf("Error setting Mail Server configuration: %d: %s", api_response.StatusCode, api_response.Status),
		)
		return
	} else if api_response.StatusCode != http.StatusCreated && api_response.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error setting Mail Server configuration",
			fmt.Sprintf("Unexpected Response Code whilst setting Mail Server configuration: %d: %s", api_response.StatusCode, api_response.Status),
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
	api_response, http_response, err := r.Client.EmailAPI.GetEmailConfiguration(ctx).Execute()

	if err != nil {
		if http_response.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"System Email Configuration does not exist",
				fmt.Sprintf("Unable to read System Email Configuration: %d: %s", http_response.StatusCode, http_response.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading System Email Configuration",
				fmt.Sprintf("Unable to read System Email Configuration: %s: %s", http_response.Status, err),
			)
		}
		return
	} else {
		// Update State
		state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		state.Enabled = types.BoolPointerValue(api_response.Enabled)
		if api_response.Host != nil {
			state.Host = types.StringPointerValue(api_response.Host)
		} else {
			state.Host = types.StringNull()
		}

		port := api_response.Port
		state.Port = types.Int64Value(int64(port))

		if api_response.Username != nil {
			state.Username = types.StringPointerValue(api_response.Username)

		// state.Name = types.StringValue(*api_response.Name)
		// state.Format = types.StringValue(*api_response.Format)
		// state.Type = types.StringValue(*api_response.Type)
		// state.Url = types.StringValue(*api_response.Url)
		// state.Online = types.BoolValue(api_response.Online)
		// state.Storage.BlobStoreName = types.StringValue(api_response.Storage.BlobStoreName)
		// state.Storage.StrictContentTypeValidation = types.BoolValue(api_response.Storage.StrictContentTypeValidation)
		// state.Storage.WritePolicy = types.StringValue(api_response.Storage.WritePolicy)
		// if api_response.Cleanup != nil {
		// 	policies := make([]types.String, len(api_response.Cleanup.PolicyNames), 0)
		// 	for i, p := range api_response.Cleanup.PolicyNames {
		// 		policies[i] = types.StringValue(p)
		// 	}
		// 	state.Cleanup = &model.RepositoryCleanupModel{
		// 		PolicyNames: policies,
		// 	}
		// }
		// state.Maven.ContentDisposition = types.StringValue(*api_response.Maven.ContentDisposition)
		// state.Maven.LayoutPolicy = types.StringValue(*api_response.Maven.LayoutPolicy)
		// state.Maven.VersionPolicy = types.StringValue(*api_response.Maven.VersionPolicy)
		// if api_response.Component != nil && api_response.Component.ProprietaryComponents != nil {
		// 	state.Component = &model.RepositoryComponentModel{
		// 		ProprietaryComponents: types.BoolValue(*api_response.Component.ProprietaryComponents),
		// 	}
		// } else {
		// 	state.Component = nil
		// }

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

	// Update API Call
	// requestPayload := sonatyperepo.MavenHostedRepositoryApiRequest{
	// 	Name:   *plan.Name.ValueStringPointer(),
	// 	Maven:  sonatyperepo.MavenAttributes{},
	// 	Online: plan.Online.ValueBool(),
	// 	Storage: sonatyperepo.HostedStorageAttributes{
	// 		BlobStoreName:               plan.Storage.BlobStoreName.ValueString(),
	// 		StrictContentTypeValidation: plan.Storage.StrictContentTypeValidation.ValueBool(),
	// 		WritePolicy:                 plan.Storage.WritePolicy.ValueString(),
	// 	},
	// }
	// if !plan.Maven.ContentDisposition.IsNull() {
	// 	requestPayload.Maven.ContentDisposition = plan.Maven.ContentDisposition.ValueStringPointer()
	// }
	// if !plan.Maven.LayoutPolicy.IsNull() {
	// 	requestPayload.Maven.LayoutPolicy = plan.Maven.LayoutPolicy.ValueStringPointer()
	// }
	// if !plan.Maven.VersionPolicy.IsNull() {
	// 	requestPayload.Maven.VersionPolicy = plan.Maven.VersionPolicy.ValueStringPointer()
	// }
	// if len(plan.Cleanup.PolicyNames) > 0 {
	// 	policies := make([]string, len(plan.Cleanup.PolicyNames), 0)
	// 	for _, p := range plan.Cleanup.PolicyNames {
	// 		policies = append(policies, p.ValueString())
	// 	}
	// 	requestPayload.Cleanup = &sonatyperepo.CleanupPolicyAttributes{
	// 		PolicyNames: policies,
	// 	}
	// }
	// if !plan.Component.ProprietaryComponents.IsNull() {
	// 	requestPayload.Component = &sonatyperepo.ComponentAttributes{
	// 		ProprietaryComponents: plan.Component.ProprietaryComponents.ValueBoolPointer(),
	// 	}
	// }
	// apiUpdateRequest := r.Client.RepositoryManagementAPI.UpdateMavenHostedRepository(ctx, state.Name.ValueString()).Body(requestPayload)

	// // Call API
	// httpResponse, err := apiUpdateRequest.Execute()

	// // Handle Error(s)
	// if err != nil {
	// 	if httpResponse.StatusCode == 404 {
	// 		resp.State.RemoveResource(ctx)
	// 		resp.Diagnostics.AddWarning(
	// 			"Maven Hosted Repository to update did not exist",
	// 			fmt.Sprintf("Unable to update Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
	// 		)
	// 	} else {
	// 		resp.Diagnostics.AddError(
	// 			"Error Updating Maven Hosted Repository",
	// 			fmt.Sprintf("Unable to update Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
	// 		)
	// 	}
	// 	return
	// } else if httpResponse.StatusCode == http.StatusNoContent {
	// 	// Map response body to schema and populate Computed attribute values
	// 	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// 	// Set state to fully populated data
	// 	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	// 	if resp.Diagnostics.HasError() {
	// 		return
	// 	}
	// } else {
	// 	resp.Diagnostics.AddError(
	// 		"Unknown Error Updating Maven Hosted Repository",
	// 		fmt.Sprintf("Unable to update Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
	// 	)
	// }
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

	DeleteRepository(r.Client, &ctx, state.Name.ValueString(), resp)
}
