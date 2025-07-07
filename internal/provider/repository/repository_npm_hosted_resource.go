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

package repository

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"terraform-provider-sonatyperepo/internal/provider/repository/format"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// repositoryNpmHostedResource is the resource implementation.
// type repositoryNpmHostedResource struct {
// 	common.BaseResource
// 	RepositoryFormat format.RepositoryFormat
// }

// Generic to all Repository Resources
type repositoryResource struct {
	common.BaseResource
	RepositoryFormat format.RepositoryFormat
	RepositoryType   format.RepositoryType
}

// Metadata returns the resource type name.
func (r *repositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.RepositoryFormat.GetResourceName(req)
}

// Set Schema for this Resource
func (r *repositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := getHostedStandardSchema(common.REPO_FORMAT_NPM)
	maps.Copy(schema.Attributes, r.RepositoryFormat.GetFormatSchemaAttributes())
	resp.Schema = schema
}

func (r *repositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	plan, diags := r.RepositoryFormat.GetPlanAsModel(ctx, req.Plan)
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
	httpResponse, err := r.RepositoryFormat.DoCreateRequest(plan, r.Client, ctx)

	// Handle Errors
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating %s %s Repository", r.RepositoryFormat.GetKey(), r.RepositoryType.String()),
			fmt.Sprintf("Error response: %s", httpResponse.Status),
		)
		return
	}
	if !slices.Contains(r.RepositoryFormat.GetApiCreateSuccessResposneCodes(), httpResponse.StatusCode) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creation of %s %s Repository was not successful", r.RepositoryFormat.GetKey(), r.RepositoryType.String()),
			fmt.Sprintf("Error response: %s", httpResponse.Status),
		)
	}

	// Call Read API as that contains more complete information for mapping to State
	apiResponse, httpResponse, err := r.RepositoryFormat.DoReadRequest(plan, r.Client, ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("%s %s Repository did not exist to read", r.RepositoryType.String(), r.RepositoryFormat.GetKey()),
				fmt.Sprintf("Error Response: %s", httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading %s %s Repository", r.RepositoryType.String(), r.RepositoryFormat.GetKey()),
				fmt.Sprintf("Error response: %s - %s", httpResponse.Status, err),
			)
		}
		return
	}

	stateModel := r.RepositoryFormat.UpdateStateFromApi(plan, apiResponse)
	stateModel = (r.RepositoryFormat.UpdatePlanForState(stateModel)).(model.RepositoryNpmHostedModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *repositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	stateModel, diags := r.RepositoryFormat.GetStateAsModel(ctx, req.State)
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
	apiResponse, httpResponse, err := r.RepositoryFormat.DoReadRequest(stateModel, r.Client, ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("%s %s Repository did not exist to read", r.RepositoryType.String(), r.RepositoryFormat.GetKey()),
				fmt.Sprintf("Error Response: %s", httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading %s %s Repository", r.RepositoryType.String(), r.RepositoryFormat.GetKey()),
				fmt.Sprintf("Error response: %s - %s", httpResponse.Status, err),
			)
		}
		return
	}

	// Update State from Response
	r.RepositoryFormat.UpdateStateFromApi(stateModel, apiResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *repositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	planModel, diags := r.RepositoryFormat.GetPlanAsModel(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)

	// Retrieve values from state
	stateModel, diags := r.RepositoryFormat.GetStateAsModel(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	// Request Context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Make API requet
	httpResponse, err := r.RepositoryFormat.DoUpdateRequest(planModel, stateModel, r.Client, ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("%s %s Repository did not exist to read", r.RepositoryType.String(), r.RepositoryFormat.GetKey()),
				fmt.Sprintf("Error Response: %s", httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading %s %s Repository", r.RepositoryType.String(), r.RepositoryFormat.GetKey()),
				fmt.Sprintf("Error response: %s - %s", httpResponse.Status, err),
			)
		}
		return
	}

	// Call Read API as that contains more complete information for mapping to State
	apiResponse, httpResponse, err := r.RepositoryFormat.DoReadRequest(planModel, r.Client, ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("%s %s Repository did not exist to read", r.RepositoryType.String(), r.RepositoryFormat.GetKey()),
				fmt.Sprintf("Error Response: %s", httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading %s %s Repository", r.RepositoryType.String(), r.RepositoryFormat.GetKey()),
				fmt.Sprintf("Error response: %s - %s", httpResponse.Status, err),
			)
		}
		return
	}

	stateModel = r.RepositoryFormat.UpdateStateFromApi(planModel, apiResponse)
	stateModel = (r.RepositoryFormat.UpdatePlanForState(stateModel)).(model.RepositoryNpmHostedModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *repositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	state, diags := r.RepositoryFormat.GetStateAsModel(ctx, req.State)
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

	// Make API requet
	httpResponse, err := r.RepositoryFormat.DoDeleteRequest(state, r.Client, ctx)

	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("%s %s Repository did not exist to delete", r.RepositoryType.String(), r.RepositoryFormat.GetKey()),
				fmt.Sprintf("Error Response: %s", httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error deleting %s %s Repository", r.RepositoryFormat.GetKey(), r.RepositoryFormat),
				fmt.Sprintf("Error response: %s", httpResponse.Status),
			)
		}
		return
	}
	if httpResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unexpected response when deleting %s %s Repository", r.RepositoryFormat.GetKey(), r.RepositoryFormat),
			fmt.Sprintf("Error response: %s", httpResponse.Status),
		)
	}
}

// ----------------

// NewRepositoryNpmHostedResource is a helper function to simplify the provider implementation.
func NewRepositoryNpmHostedResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.NpmRepositoryFormat{},
		RepositoryType:   format.REPO_TYPE_HOSTED,
	}
}

// // Update updates the resource and sets the updated Terraform state on success.
// func (r *repositoryNpmHostedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// 	// Retrieve values from plan & state
// 	var plan model.RepositoryMavenHostedModel
// 	var state model.RepositoryMavenHostedModel

// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
// 	if resp.Diagnostics.HasError() {
// 		tflog.Error(ctx, fmt.Sprintf("Getting plan data has errors: %v", resp.Diagnostics.Errors()))
// 		return
// 	}
// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
// 		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
// 		return
// 	}

// 	ctx = context.WithValue(
// 		ctx,
// 		sonatyperepo.ContextBasicAuth,
// 		r.Auth,
// 	)

// 	// Update API Call
// 	requestPayload := sonatyperepo.MavenHostedRepositoryApiRequest{
// 		Name:   *plan.Name.ValueStringPointer(),
// 		Maven:  sonatyperepo.MavenAttributes{},
// 		Online: plan.Online.ValueBool(),
// 		Storage: sonatyperepo.HostedStorageAttributes{
// 			BlobStoreName:               plan.Storage.BlobStoreName.ValueString(),
// 			StrictContentTypeValidation: plan.Storage.StrictContentTypeValidation.ValueBool(),
// 			WritePolicy:                 plan.Storage.WritePolicy.ValueString(),
// 		},
// 	}
// 	if !plan.Maven.ContentDisposition.IsNull() {
// 		requestPayload.Maven.ContentDisposition = plan.Maven.ContentDisposition.ValueStringPointer()
// 	}
// 	if !plan.Maven.LayoutPolicy.IsNull() {
// 		requestPayload.Maven.LayoutPolicy = plan.Maven.LayoutPolicy.ValueStringPointer()
// 	}
// 	if !plan.Maven.VersionPolicy.IsNull() {
// 		requestPayload.Maven.VersionPolicy = plan.Maven.VersionPolicy.ValueStringPointer()
// 	}
// 	if plan.Cleanup != nil && len(plan.Cleanup.PolicyNames) > 0 {
// 		policies := make([]string, len(plan.Cleanup.PolicyNames), 0)
// 		for _, p := range plan.Cleanup.PolicyNames {
// 			policies = append(policies, p.ValueString())
// 		}
// 		requestPayload.Cleanup = &sonatyperepo.CleanupPolicyAttributes{
// 			PolicyNames: policies,
// 		}
// 	}
// 	if !plan.Component.ProprietaryComponents.IsNull() {
// 		requestPayload.Component = &sonatyperepo.ComponentAttributes{
// 			ProprietaryComponents: plan.Component.ProprietaryComponents.ValueBoolPointer(),
// 		}
// 	}
// 	apiUpdateRequest := r.Client.RepositoryManagementAPI.UpdateMavenHostedRepository(ctx, state.Name.ValueString()).Body(requestPayload)

// 	// Call API
// 	httpResponse, err := apiUpdateRequest.Execute()

// 	// Handle Error(s)
// 	if err != nil {
// 		if httpResponse.StatusCode == 404 {
// 			resp.State.RemoveResource(ctx)
// 			resp.Diagnostics.AddWarning(
// 				"Maven Hosted Repository to update did not exist",
// 				fmt.Sprintf("Unable to update Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
// 			)
// 		} else {
// 			resp.Diagnostics.AddError(
// 				"Error Updating Maven Hosted Repository",
// 				fmt.Sprintf("Unable to update Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
// 			)
// 		}
// 		return
// 	} else if httpResponse.StatusCode == http.StatusNoContent {
// 		// Map response body to schema and populate Computed attribute values
// 		// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

// 		// Set state to fully populated data
// 		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
// 		if resp.Diagnostics.HasError() {
// 			return
// 		}
// 	} else {
// 		resp.Diagnostics.AddError(
// 			"Unknown Error Updating Maven Hosted Repository",
// 			fmt.Sprintf("Unable to update Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
// 		)
// 	}
// }

// // Delete deletes the resource and removes the Terraform state on success.
// func (r *repositoryNpmHostedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// 	// Retrieve values from state
// 	var state model.RepositoryMavenHostedModel

// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
// 		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
// 		return
// 	}

// 	ctx = context.WithValue(
// 		ctx,
// 		sonatyperepo.ContextBasicAuth,
// 		r.Auth,
// 	)

// 	DeleteRepository(r.Client, &ctx, state.Name.ValueString(), resp)
// }
