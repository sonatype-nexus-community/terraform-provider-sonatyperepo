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

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &securitySslTruststoreResource{}
	_ resource.ResourceWithImportState = &securitySslTruststoreResource{}
)

// securitySslTruststoreResource is the resource implementation.
type securitySslTruststoreResource struct {
	common.BaseResource
}

// NewSecuritySslTruststoreResource is a helper function to simplify the provider implementation.
func NewSecuritySslTruststoreResource() resource.Resource {
	return &securitySslTruststoreResource{}
}

// Metadata returns the resource type name.
func (r *securitySslTruststoreResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_ssl_truststore"
}

// Schema defines the schema for the resource.
func (r *securitySslTruststoreResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Description: "Manage SSL certificates in the Nexus truststore.",
		Attributes: map[string]tfschema.Attribute{
			"pem": schema.ResourceRequiredStringWithPlanModifier(
				"PEM-encoded certificate to add to the truststore",
				[]planmodifier.String{stringplanmodifier.RequiresReplace()},
			),
			"id":                          schema.ResourceComputedString("Certificate ID"),
			"fingerprint":                 schema.ResourceComputedString("SHA-1 fingerprint of the certificate"),
			"serial_number":               schema.ResourceComputedString("Serial number of the certificate"),
			"subject_common_name":         schema.ResourceComputedString("Subject common name"),
			"subject_organization":        schema.ResourceComputedString("Subject organization"),
			"subject_organizational_unit": schema.ResourceComputedString("Subject organizational unit"),
			"issuer_common_name":          schema.ResourceComputedString("Issuer common name"),
			"issuer_organization":         schema.ResourceComputedString("Issuer organization"),
			"issuer_organizational_unit":  schema.ResourceComputedString("Issuer organizational unit"),
			"issued_on":                   schema.ResourceComputedInt64("Certificate issued on (epoch milliseconds)"),
			"expires_on":                  schema.ResourceComputedInt64("Certificate expires on (epoch milliseconds)"),
			"last_updated":                schema.ResourceLastUpdated(),
		},
	}
}

// ImportState imports the resource into Terraform state.
func (r *securitySslTruststoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create creates the resource and sets the initial Terraform state.
func (r *securitySslTruststoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.SecuritySslTruststoreModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	apiResponse, httpResponse, err := r.Client.SecurityCertificatesAPI.AddCertificate(ctx).Body(plan.Pem.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				"Your user is unauthorized to access this resource or feature.",
			)
		} else {
			errors.HandleAPIError(
				"Error adding certificate to truststore",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	plan.MapFromApi(apiResponse)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *securitySslTruststoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.SecuritySslTruststoreModel

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

	apiResponse, httpResponse, err := r.Client.SecurityCertificatesAPI.GetTrustStoreCertificates(ctx).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				"Your user is unauthorized to access this resource or feature.",
			)
		} else {
			errors.HandleAPIError(
				"Error reading truststore certificates",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Find the certificate by ID
	var found bool
	for _, cert := range apiResponse {
		if cert.Id != nil && *cert.Id == state.Id.ValueString() {
			state.MapFromApi(&cert)
			found = true
			break
		}
	}

	if !found {
		// Certificate was deleted outside of Terraform
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *securitySslTruststoreResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// No-op: pem has RequiresReplace, so Terraform will never call Update.
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *securitySslTruststoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.SecuritySslTruststoreModel

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

	httpResponse, err := r.Client.SecurityCertificatesAPI.RemoveCertificate(ctx, state.Id.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			// Certificate already removed, nothing to do
			return
		}
		if httpResponse.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unauthorized",
				"Your user is unauthorized to access this resource or feature.",
			)
		} else {
			errors.HandleAPIError(
				"Error removing certificate from truststore",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}
}
